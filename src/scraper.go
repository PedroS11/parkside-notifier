package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func GetParksideCards() []string {
	urls := getFlyerUrls()

	return urls
}

func getFlyerUrls() []string {
	c := colly.NewCollector(
		colly.AllowedDomains("www.lidl.pt"),
	)

	var flyerUrls []string

	// triggered when the scraper encounters an error
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	isSemanaisFound := false

	// triggered when a CSS selector matches an element
	c.OnHTML(".subcategory", func(e *colly.HTMLElement) {
		if !isSemanaisFound {
			// printing all URLs associated with the <a> tag on the page
			e.ForEach("a", func(i int, h *colly.HTMLElement) {
				fmt.Println(i, h.Attr("href"))
				flyerUrls = append(flyerUrls, strings.Replace(h.Attr("href"), "/ar/0", "/view/flyer/page/1", 1))
			})
			isSemanaisFound = true
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	c.Visit("https://www.lidl.pt/c/folhetos/s10020672")

	return flyerUrls
}

func ParseFlyer22(url string) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.lidl.pt"),
	)

	// triggered when the scraper encounters an error
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// triggered when a CSS selector matches an element
	c.OnHTML(".page__wrapper", func(e *colly.HTMLElement) {
		// printing all URLs associated with the <a> tag on the page
		e.ForEach("img", func(i int, h *colly.HTMLElement) {
			fmt.Println(i, h.Attr("src"))
		})
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", string(r.Body))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	c.Visit(url)
}
