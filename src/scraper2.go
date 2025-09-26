package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func GetUrls() []string {
	page := rod.New().NoDefaultDevice().MustConnect().MustPage("https://www.lidl.pt/c/folhetos/s10020672")
	page.MustWindowFullscreen()

	page.MustElement("#onetrust-reject-all-handler").MustClick()
	// Get just the first one as it's Semanais
	subCategory := page.MustElement(".subcategory")
	cards := subCategory.MustElements("a")

	var urls []string

	for _, card := range cards {
		url := card.MustProperty("href").String()

		urls = append(urls, strings.Replace(url, "/ar/0", "/view/flyer/page/1", 1))
	}

	return urls
}

func ParseFlyer(fylerUrl string) []string {
	page := rod.New().NoDefaultDevice().MustConnect().MustPage(fylerUrl)
	page.MustWindowFullscreen()

	page.MustElement("#onetrust-reject-all-handler").MustClick()
	var urls []string

	for {
		flyerPages := page.MustElements(".page--current")
		foundFinalPage := false
		var nextPage *rod.Element

		fmt.Println("\n\nCrawling", page.MustInfo().URL)

		for _, flyer := range flyerPages {
			url := flyer.MustElement("img").MustProperty("src")
			fmt.Println(url)
			urls = append(urls, url.String())

			navigationArrows, _ := page.Timeout(1 * time.Second).Elements(".button--navigation-lidl")
			if len(navigationArrows) == 1 && *navigationArrows[0].MustAttribute("aria-label") == "Página anterior" {
				foundFinalPage = true

				break
			}
			nextPage = navigationArrows[len(navigationArrows)-1]
		}

		if foundFinalPage {
			fmt.Println("ALL URLS", urls)
			break
		}

		nextPage.MustClick()
	}

	return urls

	// fmt.Println("FIM")

	// time.Sleep(time.Hour)
}
