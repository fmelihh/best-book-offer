package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Book struct {
	Url         string  `json:"url"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
}

const BaseUrl = "https://books.toscrape.com"

func main() {
	start := time.Now()
	c := colly.NewCollector(colly.AllowedDomains("books.toscrape.com"), colly.Async(true))
	url := BaseUrl + "/catalogue/page-1.html"
	var items = make([]Book, 0, 20)

	c.OnHTML("ol.row li article.product_pod", func(element *colly.HTMLElement) {
		bookUrl := element.ChildAttr("a", "href")
		bookUrl = element.Request.AbsoluteURL(bookUrl)

		if err := c.Visit(bookUrl); err != nil {
			log.Fatal(err)
		}
	})

	c.OnHTML(".product_page", func(element *colly.HTMLElement) {

		bookTitle := element.ChildText(".row .product_main h1")
		bookTitle = strings.TrimSpace(bookTitle)

		priceSection := element.ChildText("#content_inner > article > table > tbody > tr:nth-child(4) > td")
		priceSection = strings.TrimSpace(priceSection)

		priceUnicodeArr := []rune(priceSection)

		currency := string(priceUnicodeArr[0])

		priceString := string(priceUnicodeArr[1:])
		price, err := strconv.ParseFloat(priceString, 64)
		if err != nil {
			log.Fatal(err)
		}

		description := element.ChildText("#content_inner > article > p")
		description = strings.TrimSpace(description)

		book := Book{
			Name:        bookTitle,
			Description: description,
			Price:       price,
			Url:         element.Request.URL.String(),
			Currency:    currency,
		}
		items = append(items, book)

	})

	c.OnHTML("ul.pager li.next", func(element *colly.HTMLElement) {

		nextUrl := element.ChildAttr("a", "href")
		nextUrl = element.Request.AbsoluteURL(nextUrl)

		fmt.Printf("%s page has ended. Switching to the next page %s", element.Request.URL.String(), nextUrl)
		err := c.Visit(nextUrl)
		if err != nil {
			log.Fatal(err)
		}

	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("visiting", request.URL)

	})

	err := c.Visit(url)
	c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Crawling process has ended. Total crawled data is %d", len(items))

	file, err := os.Create("books.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Url", "Name", "Price", "Currency", "Description"}
	if err := writer.Write(header); err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		stringPrice := strconv.FormatFloat(item.Price, 'f', -1, 64)
		row := []string{item.Url, item.Name, stringPrice, item.Currency, item.Description}
		if err := writer.Write(row); err != nil {
			log.Fatal(err)
		}
	}

	t := time.Since(start)
	fmt.Printf("\n\nSpending Time: %v", t)
	fmt.Printf("\nTotal Crawled data saved to the books.csv file.")
}
