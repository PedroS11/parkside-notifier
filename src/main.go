package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"
)

func getFlyersAndNotify() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := CreateClient()

	flyers := GetFlyers()

	bot, ctx := Start(ctx)

	for _, flyer := range flyers {
		isNotified, err := WasUrlNotified(client, ctx, flyer.Url)

		if err != nil {
			fmt.Println("Error checking URL:", err)
			continue // Skip this card and move on to the next one
		}

		// Only call SendMediaGroup if the URL was not notified
		if !isNotified {
			SendMediaGroup(bot, ctx, flyer)
		}
	}
}

func main() {
	// create a scheduler
	s, err := CreateCronJob(getFlyersAndNotify)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {}
}
