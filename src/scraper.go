package main

import (
	"log/slog"
	"parksideNotifier/src/interfaces"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func crawlFlyers() []interfaces.Flyer {
	l := launcher.MustNewManaged("")

	// You can also set any flag remotely before you launch the remote browser.
	// Available flags: https://peter.sh/experiments/chromium-command-line-switches
	l.Set("disable-gpu").Delete("disable-gpu")

	// Launch with headful mode
	l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	browser := rod.New().Client(l.MustClient()).MustConnect()

	page := browser.MustPage("https://www.lidl.pt/c/folhetos/s10020672")

	page.MustElement("#onetrust-reject-all-handler").MustClick()
	// Get just the first one as it's Semanais
	subCategory := page.MustElement(".subcategory")
	promotionCards := subCategory.MustElements("a")

	var cards []interfaces.Flyer

	for _, promotionCard := range promotionCards {
		url := promotionCard.MustProperty("href").String()
		viewUrl := strings.Replace(url, "/ar/0", "/view/flyer/page/1", 1)
		card := interfaces.Flyer{
			Url:          viewUrl,
			Name:         promotionCard.MustElement(".flyer__name").MustText(),
			PreviewImage: promotionCard.MustElement(".flyer__image").MustProperty("src").String(),
			Date:         promotionCard.MustElement(".flyer__title").MustText(),
			Images:       []string{},
			Products:     []interfaces.Product{},
		}

		if card.Name == "Novidades" {
			cards = append(cards, card)
		}
	}

	return cards
}

func parseFlyer(flyerUrl string) []string {
	page := rod.New().NoDefaultDevice().MustConnect().MustPage(flyerUrl)
	page.MustWindowFullscreen()

	slog.Info("Crawling", slog.String("url", flyerUrl))

	// Reject cookies
	page.MustElement("#onetrust-reject-all-handler").MustClick()

	var flyerPageUrls []string

	for {
		flyerPages := page.MustElements(".page--current")
		foundFinalPage := false
		var nextPage *rod.Element

		for _, flyer := range flyerPages {
			url := flyer.MustElement("img").MustProperty("src")
			flyerPageUrls = append(flyerPageUrls, url.String())

			navigationArrows, _ := page.Timeout(1 * time.Second).Elements(".button--navigation-lidl")

			if len(navigationArrows) == 1 {
				previousPageButtonText := *navigationArrows[0].MustAttribute("aria-label")
				if previousPageButtonText == "Página anterior" {

					foundFinalPage = true

					break
				}
			}
			// Get last arrow button, it will be the move forward one
			// As i checked that when there's just one, it isn't the move backwards one
			nextPage = navigationArrows[len(navigationArrows)-1]
		}

		if foundFinalPage {
			break
		}

		nextPage.MustClick()
	}

	return flyerPageUrls
}

func GetFlyers() []interfaces.Flyer {
	flyers := crawlFlyers()

	var wg sync.WaitGroup

	for i := range flyers {
		wg.Add(1)
		go func(flyer *interfaces.Flyer) {
			defer wg.Done()

			images := parseFlyer(flyer.Url)
			flyer.Images = images
		}(&flyers[i])
	}

	wg.Wait()

	return flyers
}
