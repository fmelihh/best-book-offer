package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

type Book struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Url         string  `json:"url"`
	Currency    string  `json:"currency"`
}

const BaseUrl = "https://books.toscrape.com"

func main() {
	c := colly.NewCollector()
	url := BaseUrl + "/catalogue/page-1.html"

	var items []Book

	c.OnHTML(".product_pod", func(element *colly.HTMLElement) {
		bookUrl := element.ChildAttr(".image_container a", "href")
		bookUrl = BaseUrl + "/catalogue/" + bookUrl
		err := c.Visit(bookUrl)
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
	if err != nil {
		log.Fatal(err)
	}
}
