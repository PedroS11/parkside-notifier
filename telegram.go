package main

import (
	"context"
	"fmt"
	"os"
	"parksideNotifier/interfaces"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(ctx context.Context) (*bot.Bot, context.Context) {
	opts := []bot.Option{
		bot.WithCheckInitTimeout(time.Second * 5),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)

	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	return b, ctx
}

func SendMediaGroup(b *bot.Bot, ctx context.Context, c interfaces.Card) {
	media1 := &models.InputMediaPhoto{
		Media:     c.Img,
		Caption:   c.Date + " veja [" + c.Name + "](" + c.Url + ")",
		ParseMode: models.ParseModeMarkdownV1,
	}

	fmt.Println(media1)

	_, err := b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
		ChatID: os.Getenv("CHANNEL_ID"),
		Media: []models.InputMedia{
			media1,
		},
	})

	if err != nil {
		fmt.Println("ERROR", err.Error())
	}

}
