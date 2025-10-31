package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func getFlyersAndNotify() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// client := CreateClient()

	flyers := GetFlyers()

	fmt.Println("Found", len(flyers), "flyers")

	bot, ctx := Start(ctx)

	for i, flyer := range flyers {

		products := GetProductsWithOpenAI(flyer.Images)

		// Sleep 1min to avoid open ai returning 429 - Too Many Requests error
		if i < len(flyers)-1 {
			time.Sleep(time.Minute)
		}

		if len(products) == 0 {
			fmt.Printf("Flyer %s has no Parkside products", flyer.Url)
			continue
		}

		flyer.Products = append(flyer.Products, products...)
		fmt.Printf("Flyer %s has %v\n", flyer.Url, flyer.Products)

		// isNotified, err := WasUrlNotified(client, ctx, flyer.Url)

		// if err != nil {
		// 	fmt.Println("Error checking URL:", err)
		// 	continue // Skip this card and move on to the next one
		// }

		// // Only call SendMediaGroup if the URL was not notified
		// if !isNotified {
		SendMediaGroup(bot, ctx, flyer)
		// }
	}

	for _, flyer := range flyers {
		fmt.Printf("Flyer %s on %s has %v\n", flyer.Name, flyer.Url, flyer.Products)
	}
}

func main() {
	getFlyersAndNotify()

	// // create a scheduler
	// s, err := CreateCronJob(getFlyersAndNotify)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// // start the scheduler
	// s.Start()

	// // block until you are ready to shut down
	// select {}

	// GetProductsWithOpenAI("https://imgproxy.leaflets.schwarz/-0HK5TwsNHt8hvdlLp_-10eu1gc2oFY6wPPf_rpDgLM/rs:fit:1200:1200:1/g:no/czM6Ly9sZWFmbGV0cy9pbWFnZXMvMDE5YTEyMTctZDEwYS03ZjYxLTk2ZTEtYTMwYjE0ZjE5MmMyL3BhZ2UtMDdfNjg1ZjFmMTlkYzNmMDAwMDBjMTZjZjhmMDE2ODVkNzMucG5n.jpg")

}
