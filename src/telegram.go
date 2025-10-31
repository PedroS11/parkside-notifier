package main

import (
	"context"
	"fmt"
	"os"
	"parksideNotifier/src/interfaces"
	"strings"
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

func SendMediaGroup(b *bot.Bot, ctx context.Context, flyer interfaces.Flyer) {
	template := "Dia %s no folheto [%s](%s), existirão os seguintes produtos Parkside: %s"

	productsName := []string{}

	for _, product := range flyer.Products {
		productsName = append(productsName, product.Name)
	}

	mediaMessage := &models.InputMediaPhoto{
		Media:     flyer.PreviewImage,
		Caption:   fmt.Sprintf(template, flyer.Date, flyer.Name, flyer.Url, strings.Join(productsName, ",")),
		ParseMode: models.ParseModeMarkdownV1,
	}

	_, err := b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
		ChatID: os.Getenv("CHANNEL_ID"),
		Media: []models.InputMedia{
			mediaMessage,
		},
	})

	if err != nil {
		fmt.Println("ERROR", err.Error())
	}
}
