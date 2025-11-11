package main

import (
	"fmt"
	"log/slog"
	"os"
	"parksideNotifier/src/interfaces"
	"strings"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func crawlFlyers(browser *rod.Browser) []interfaces.Flyer {

	page := browser.MustPage("https://www.lidl.pt/c/folhetos/s10020672")

	page.MustElement("#onetrust-reject-all-handler").MustClick()

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

func parseFlyer(browser *rod.Browser, flyerUrl string) []string {
	page := browser.MustPage(flyerUrl)

	// Reduce browser size to only show one page instead of displaying two pages from the PDF
	page.MustSetViewport(912, 1368, 1, false)

	fmt.Println("Crawling", slog.String("url", flyerUrl))

	// Reject cookies
	// page.MustElement("#onetrust-reject-all-handler").MustClick()

	fmt.Println("COOKIES")

	var flyerPageUrls []string

	foundFirstPage := false

	for {
		flyerPage := page.MustElement(".page--current")
		fmt.Println(page.MustInfo().URL)
		foundFinalPage := false
		var nextPage *rod.Element

		url := flyerPage.MustElement("img").MustProperty("src")

		flyerPageUrls = append(flyerPageUrls, url.String())

		navigationArrows, _ := page.Elements(".button--navigation-lidl")

		fmt.Println("DEPOIS", len(navigationArrows))

		if len(navigationArrows) == 1 {
			if foundFirstPage {
				foundFinalPage = true
			} else {
				foundFirstPage = true
				nextPage = navigationArrows[0]
			}
		} else {
			nextPage = navigationArrows[1]

		}

		if foundFinalPage {
			break
		}

		nextPage.MustClick()
	}

	return flyerPageUrls
}

func GetFlyers() []interfaces.Flyer {
	// This example is to launch a browser remotely, not connect to a running browser remotely,
	// to connect to a running browser check the "../connect-browser" example.
	// Rod provides a docker image for beginners, run the below to start a launcher.Manager:
	//
	//     docker run --rm -p 7317:7317 ghcr.io/go-rod/rod
	//
	// For available CLI flags run: docker run --rm ghcr.io/go-rod/rod rod-manager -h
	// For more information, check the doc of launcher.Manager
	l := launcher.MustNewManaged(os.Getenv("ROD_URL"))

	// You can also set any flag remotely before you launch the remote browser.
	// Available flags: https://peter.sh/experiments/chromium-command-line-switches
	l.Set("disable-gpu").Delete("disable-gpu")

	// Launch with headful mode
	l.Headless(true).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	browser := rod.New().Client(l.MustClient()).MustConnect()

	// You may want to start a server to watch the screenshots of the remote browser.
	// launcher.Open(browser.ServeMonitor(""))

	flyers := crawlFlyers(browser)

	// flyers := []interfaces.Flyer{{
	// 	Url: "https://www.lidl.pt/l/pt/folhetos/novidades-a-partir-de-10-11/view/flyer/page/1?lf=HHZ",
	// }}

	slog.Info(fmt.Sprintf("There are %d flyers", len(flyers)))

	var wg sync.WaitGroup

	for i := range flyers {
		wg.Add(1)
		go func(flyer *interfaces.Flyer) {
			defer wg.Done()

			images := parseFlyer(browser, flyer.Url)
			flyer.Images = images
		}(&flyers[i])
	}

	wg.Wait()

	return flyers
}
