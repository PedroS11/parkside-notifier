package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"
)

func getCardsAndNotify() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := CreateClient()

	parksideCards := GetParksideCards()

	fmt.Println("parksideCards", len(parksideCards))
	bot, ctx := Start(ctx)

	for _, card := range parksideCards {
		isNotified, err := WasUrlNotified(client, ctx, card.Url)

		if err != nil {
			fmt.Println("Error checking URL:", err)
			continue // Skip this card and move on to the next one
		}

		// Only call SendMediaGroup if the URL was not notified
		if !isNotified {
			SendMediaGroup(bot, ctx, card)
		}
	}
}
func main() {
	// create a scheduler
	s, err := CreateCronJob(getCardsAndNotify)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {}
}
