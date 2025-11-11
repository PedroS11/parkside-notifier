// package main

// import (
// 	"context"
// 	"fmt"
// 	"log/slog"
// 	"os"
// 	"os/signal"
// 	"time"

// 	_ "github.com/joho/godotenv/autoload"
// )

// func getFlyersAndNotify() {
// 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
// 	defer cancel()

// 	client := CreateClient()

// 	flyers := GetFlyers()

// 	slog.Info(fmt.Sprintf("Found %d flyers", len(flyers)))

// 	bot, ctx := Start(ctx)

// 	for i, flyer := range flyers {
// 		isNotified, err := WasUrlNotified(client, ctx, flyer.Url)

// 		if err != nil {
// 			errorMessage := fmt.Sprintf("Error checking URL: %s", err.Error())

// 			slog.Error(errorMessage)
// 			SendErrorMessage(bot, ctx, errorMessage)

// 			continue // Skip this card and move on to the next one
// 		}

// 		// Only call SendMediaGroup if the URL was not notified
// 		if !isNotified {
// 			products, err := GetProductsFromUrls(flyer.Images)

// 			if err != nil {
// 				errorMessage := fmt.Sprintf("Error getting products from URLs: %s", err.Error())

// 				slog.Error(errorMessage)
// 				SendErrorMessage(bot, ctx, errorMessage)

// 				continue
// 			}

// 			// Sleep 1min to avoid open ai returning 429 - Too Many Requests error
// 			if i < len(flyers)-1 {
// 				time.Sleep(time.Minute)
// 			}

// 			if len(products) == 0 {
// 				UpdateMessage(client, ctx, flyer.Url, 1)

// 				slog.Warn(fmt.Sprintf("Flyer %s %s has no Parkside products", flyer.Name, flyer.Url))
// 				continue
// 			}

// 			flyer.Products = append(flyer.Products, products...)

// 			slog.Info(fmt.Sprintf("Flyer %s has %v\n", flyer.Url, flyer.Products))

// 			SendMediaGroup(bot, ctx, flyer)
// 			UpdateMessage(client, ctx, flyer.Url, 1)
// 		} else {
// 			slog.Info(fmt.Sprintf("Flyer %s on %s was already processed", flyer.Name, flyer.Url))
// 		}

// 		slog.Info(fmt.Sprintf("All %d flyers were analysed successfully", len(flyers)))
// 	}
// }

// func main() {
// 	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
// 	getFlyersAndNotify()
// 	// // create a scheduler
// 	// s, err := CreateCronJob(getFlyersAndNotify)
// 	// if err != nil {
// 	// 	LogError("main", err)
// 	// 	os.Exit(1)
// 	// }

// 	// // start the scheduler
// 	// s.Start()

// 	// // block until you are ready to shut down
// 	// select {}
// }

// Package main ...

// FUNCIINOU
package main

import (
	"fmt"
	"os"
	"parksideNotifier/src/interfaces"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

func main() {
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
	l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	browser := rod.New().Client(l.MustClient()).MustConnect()

	// You may want to start a server to watch the screenshots of the remote browser.
	launcher.Open(browser.ServeMonitor(""))

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

	fmt.Println(cards)

	utils.Pause()
}
