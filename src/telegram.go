package main

import (
	"context"
	"fmt"
	"os"
	"parksideNotifier/src/interfaces"
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
	template := "*%s no folheto [%s](%s), vão estar disponíveis os seguintes produtos Parkside:*\n"
	intro := fmt.Sprintf(template,
		EscapeMarkdownV2(flyer.Date),
		EscapeMarkdownV2(flyer.Name),
		EscapeMarkdownV2(flyer.Url),
	)

	for _, product := range flyer.Products {
		intro += fmt.Sprintf(" • %s por %s€ \n", EscapeMarkdownV2(product.Name), EscapeMarkdownV2(fmt.Sprintf("%.2f", product.Price)))
	}

	mediaMessage := &models.InputMediaPhoto{
		Media:     flyer.PreviewImage,
		Caption:   intro,
		ParseMode: models.ParseModeMarkdown,
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
