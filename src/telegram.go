package main

import (
	"context"
	"fmt"
	"log/slog"
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
		slog.Error("Telegram Start", "error", err.Error())
		os.Exit(1)
	}

	return b, ctx
}

func SendMediaGroup(b *bot.Bot, ctx context.Context, flyer interfaces.Flyer) (*models.Message, error) {
	template := "*%s no folheto [%s](%s), vão estar disponíveis %d produtos Parkside*\n"
	caption := fmt.Sprintf(template,
		EscapeMarkdownV2(flyer.Date),
		EscapeMarkdownV2(flyer.Name),
		EscapeMarkdownV2(flyer.Url),
		len(flyer.Products),
	)

	var productsMessage string
	for _, product := range flyer.Products {
		productsMessage += fmt.Sprintf(" • %s por %s€ \n", EscapeMarkdownV2(product.Name), EscapeMarkdownV2(fmt.Sprintf("%.2f", product.Price)))
	}

	mediaMessage := &models.InputMediaPhoto{
		Media:     flyer.PreviewImage,
		Caption:   caption,
		ParseMode: models.ParseModeMarkdown,
	}

	messages, err := b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
		ChatID: os.Getenv("CHANNEL_ID"),
		Media: []models.InputMedia{
			mediaMessage,
		},
	})

	if err != nil {
		slog.Error("SendMediaGroup", "error", err.Error())
		return messages[0], err
	}

	message, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    os.Getenv("CHANNEL_ID"),
		Text:      productsMessage,
		ParseMode: models.ParseModeMarkdown,
	})

	if err != nil {
		slog.Error("SendMessage", "error", err.Error())
		return message, err
	}

	return message, nil
}

func SendErrorMessage(b *bot.Bot, ctx context.Context, message string) (*models.Message, error) {
	response, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: os.Getenv("DEBUG_CHANNEL_ID"),
		Text:   fmt.Sprintln("Error", message),
	})

	if err != nil {
		slog.Error("SendErrorMessage", "error", err.Error())
		return response, err
	}

	return response, nil
}
