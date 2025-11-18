package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	rqlitehttp "github.com/rqlite/rqlite-go-http"
)

func CreateClient() *rqlitehttp.Client {
	client, err := rqlitehttp.NewClient("http://"+os.Getenv("HTTP_ADDR"), nil)
	if err != nil {
		slog.Error("DB CreateClient", "error", err.Error())

		panic(err)
	}

	return client
}

func WasUrlNotified(client *rqlitehttp.Client, ctx context.Context, url string) (bool, error) {
	qResp, err := client.Query(ctx, rqlitehttp.SQLStatements{
		{
			SQL:              "SELECT * FROM message WHERE url = ?",
			PositionalParams: []any{url},
		},
	}, nil)

	if err != nil {
		slog.Error("WasUrlNotified", "error", err.Error())
		return false, err
	}

	rowFound := len(qResp.GetQueryResults()[0].Values) != 0

	if !rowFound {
		_, err := insertUrl(client, ctx, url)
		return false, err
	}

	var wasNotified int

	row := qResp.GetQueryResults()[0].Values[0]

	if notifiedFloat, ok := row[1].(float64); ok {
		wasNotified = int(notifiedFloat)
	} else {
		wasNotified = row[1].(int)
	}

	return wasNotified == 1, nil
}

func insertUrl(client *rqlitehttp.Client, ctx context.Context, url string) (*rqlitehttp.ExecuteResponse, error) {
	response, err := client.Execute(ctx, rqlitehttp.SQLStatements{
		{
			SQL:              "INSERT INTO message(url, notified) VALUES(?, ?)",
			PositionalParams: []any{url, 0},
		},
	}, nil)

	if err != nil {
		slog.Error("insertUrl error", err.Error(), url)
		return nil, err
	}

	errored, _, errorMessage := response.HasError()

	if errored {
		slog.Error("inserted url error", "error", errorMessage, "url", url)
		return nil, errors.New(errorMessage)
	}

	slog.Info("Calling insertUrl", "ExecuteResponse", response)

	return response, nil
}

func UpdateMessage(client *rqlitehttp.Client, ctx context.Context, url string, notify int) (*rqlitehttp.ExecuteResponse, error) {
	response, err := client.Execute(ctx, rqlitehttp.SQLStatements{
		{
			SQL:              "UPDATE message SET notified = ? WHERE url = ?",
			PositionalParams: []any{notify, url},
		},
	}, nil)

	if err != nil {
		slog.Error("UpdateMessage", "error", err.Error(), "notify", notify, "url", url)
		return nil, err
	}

	errored, _, errorMessage := response.HasError()

	if errored {
		slog.Error("Updated message", "error", errorMessage)

		return nil, errors.New(errorMessage)
	}

	slog.Info("Calling UpdateMessage", "ExecuteResponse", response)

	return response, nil
}
