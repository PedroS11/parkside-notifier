package main

import (
	"fmt"
	"parksideNotifier/interfaces"
	"strings"

	"github.com/gocolly/colly"
)

func filterParkside(cards []interfaces.Card) []interfaces.Card {
	result := make([]interfaces.Card, 0)
	for _, card := range cards {
		if strings.Contains(strings.ToLower(card.Name), "parkside") || strings.Contains(strings.ToLower(card.Name), "ferramenta") {
			result = append(result, card)
		}
	}

	return result
}

func GetParksideCards() []interfaces.Card {
	c := colly.NewCollector(
		colly.AllowedDomains("www.lidl.pt"),
	)

	var cards []interfaces.Card

	// triggered when the scraper encounters an error
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// triggered when a CSS selector matches an element
	c.OnHTML(".AHeroStageItems__List", func(e *colly.HTMLElement) {
		// printing all URLs associated with the <a> tag on the page
		e.ForEach("li", func(i int, h *colly.HTMLElement) {
			card := interfaces.Card{
				Url:  h.ChildAttr("a", "href"),
				Name: h.ChildText(".AHeroStageItems__Item--Headline"),
				Date: h.ChildText(".AHeroStageItems__Item--SubHeadline"),
				Img:  "https://lidl.pt" + h.ChildAttr("img", "src"),
			}
			cards = append(cards, card)
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	c.Visit("https://www.lidl.pt")

	return filterParkside(cards)
}
