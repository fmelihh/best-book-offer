package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
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
	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"))
	url := BaseUrl + "/catalogue/page-1.html"

	var items []Book

	c.OnHTML("article.product_pod", func(element *colly.HTMLElement) {
		bookUrl := element.ChildAttr(".image_container a", "href")
		err := c.Visit(element.Request.AbsoluteURL(bookUrl))
		if err != nil {
			log.Fatal(err)
		}
	})

	c.OnHTML(".product_page", func(element *colly.HTMLElement) {

		bookTitle := element.ChildText(".row .product_main h1")
		bookTitle = strings.TrimSpace(bookTitle)

		priceSection := element.ChildText(".row .price_color")
		priceSection = strings.TrimSpace(priceSection)

		priceUnicodeArr := []rune(priceSection)

		currency := string(priceUnicodeArr[0])

		priceString := string(priceUnicodeArr[1:])
		price, _ := strconv.ParseFloat(priceString, 64)

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

	c.OnHTML(".pager", func(element *colly.HTMLElement) {
		nextUrl := element.ChildAttr(".next a", "href")

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
}
