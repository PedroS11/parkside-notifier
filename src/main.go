package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func getFlyersAndNotify() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := CreateClient()

	flyers := GetFlyers()

	slog.Info(fmt.Sprintf("Found %d flyers", len(flyers)))

	bot, ctx := Start(ctx)

	for i, flyer := range flyers {
		isNotified, err := WasUrlNotified(client, ctx, flyer.Url)

		if err != nil {
			errorMessage := fmt.Sprintf("Error checking URL: %s", err.Error())

			slog.Error(errorMessage)
			SendErrorMessage(bot, ctx, errorMessage)

			continue // Skip this card and move on to the next one
		}

		// Only call SendMediaGroup if the URL was not notified
		if !isNotified {
			products, err := GetProductsFromUrls(flyer.Images)

			if err != nil {
				errorMessage := fmt.Sprintf("Error getting products from URLs: %s", err.Error())

				slog.Error(errorMessage)
				SendErrorMessage(bot, ctx, errorMessage)

				continue
			}

			// Sleep 1min to avoid open ai returning 429 - Too Many Requests error
			if i < len(flyers)-1 {
				time.Sleep(time.Minute)
			}

			if len(products) == 0 {
				UpdateMessage(client, ctx, flyer.Url, 1)

				slog.Warn(fmt.Sprintf("Flyer %s %s has no Parkside products", flyer.Name, flyer.Url))
				continue
			}

			flyer.Products = append(flyer.Products, products...)

			slog.Info(fmt.Sprintf("Flyer %s has %v\n", flyer.Url, flyer.Products))

			SendMediaGroup(bot, ctx, flyer)
			UpdateMessage(client, ctx, flyer.Url, 1)
		} else {
			slog.Info(fmt.Sprintf("Flyer %s on %s was already processed", flyer.Name, flyer.Url))
		}

		slog.Info(fmt.Sprintf("All %d flyers were analysed successfully", len(flyers)))
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	getFlyersAndNotify()
	// // create a scheduler
	// s, err := CreateCronJob(getFlyersAndNotify)
	// if err != nil {
	// 	LogError("main", err)
	// 	os.Exit(1)
	// }

	// // start the scheduler
	// s.Start()

	// // block until you are ready to shut down
	// select {}
}

// Package main ...

// FUNCIINOU
// package main

// import (
// 	"fmt"
// 	"log/slog"
// 	"os"
// )

// func main() {
// 	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

// 	fmt.Println(GetFlyers())
// }
