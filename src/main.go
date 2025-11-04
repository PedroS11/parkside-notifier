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

	slog.Info(fmt.Sprintln("Found", len(flyers), "flyers"))

	bot, ctx := Start(ctx)

	for i, flyer := range flyers {
		products, err := GetProductsFromUrls(flyer.Images)
		if err != nil {
			errorMessage := fmt.Sprintln("Error getting products from URLs:", err.Error())

			slog.Error(errorMessage)
			SendErrorMessage(bot, ctx, errorMessage)

			continue
		}

		// Sleep 1min to avoid open ai returning 429 - Too Many Requests error
		if i < len(flyers)-1 {
			time.Sleep(time.Minute)
		}

		if len(products) == 0 {
			slog.Warn(fmt.Sprintf("Flyer %s %s has no Parkside products", flyer.Name, flyer.Url))
			continue
		}

		flyer.Products = append(flyer.Products, products...)
		slog.Info(fmt.Sprintf("Flyer %s has %v\n", flyer.Url, flyer.Products))

		isNotified, err := WasUrlNotified(client, ctx, flyer.Url)

		if err != nil {
			errorMessage := fmt.Sprintln("Error checking URL:", err.Error())

			slog.Error(errorMessage)
			SendErrorMessage(bot, ctx, errorMessage)

			continue // Skip this card and move on to the next one
		}

		// Only call SendMediaGroup if the URL was not notified
		if !isNotified {
			SendMediaGroup(bot, ctx, flyer)
		}
	}

	for _, flyer := range flyers {
		slog.Info(fmt.Sprintf("Flyer %s on %s has %v\n", flyer.Name, flyer.Url, flyer.Products))
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
