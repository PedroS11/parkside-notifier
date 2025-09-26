package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"
)

// func getCardsAndNotify() {
// 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
// 	defer cancel()

// 	client := CreateClient()

// 	parksideCards := GetParksideCards()

// 	fmt.Println("parksideCards", len(parksideCards))
// 	bot, ctx := Start(ctx)

// 	for _, card := range parksideCards {
// 		isNotified, err := WasUrlNotified(client, ctx, card.Url)

// 		if err != nil {
// 			fmt.Println("Error checking URL:", err)
// 			continue // Skip this card and move on to the next one
// 		}

//			// Only call SendMediaGroup if the URL was not notified
//			if !isNotified {
//				SendMediaGroup(bot, ctx, card)
//			}
//		}
//	}

func main() {
	// // create a scheduler
	// s, err := CreateCronJob(getCardsAndNotify)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// // start the scheduler
	// s.Start()

	// // block until you are ready to shut down
	// select {}

	// url := "https://imgproxy.leaflets.schwarz/IU45FGYfnUNLePy3DkflXKqQme3e-ukpAwPnwn0ot_U/rs:fit:1200:1200:1/g:no/czM6Ly9sZWFmbGV0cy9pbWFnZXMvMDE5OTU3OWQtM2RlMy03Yjk0LTg0MjktMjQwNzllYWQ2YjI4L3BhZ2UtMDNfODVmNzYxNjY4MzM2Yjc1NjRlM2FmODQ5NWJiMTA1NzQucG5n.jpg"
	// base := ImageToBase64(url)
	// os.WriteFile("./base64.txt", []byte(base), 0666)

	// cards := GetParksideCards()
	// fmt.Println(cards)
	// ParseFlyer("https://www.lidl.pt/l/pt/folhetos/novidades-a-partir-de-22-09/view/flyer/page/1?lf=HHZ")
	// fmt.Println(GetUrls())
	urls := GetUrls()
	for _, url := range urls {
		fmt.Println(ParseFlyer(url))
	}

}
