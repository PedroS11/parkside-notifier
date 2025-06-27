package main

import (
	"context"
	"errors"
	"fmt"

	rqlitehttp "github.com/rqlite/rqlite-go-http"
)

func CreateClient() *rqlitehttp.Client {
	client, err := rqlitehttp.NewClient("http://myrqlite-host-1:4001", nil)
	if err != nil {
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
		fmt.Println("ERROR", err.Error())
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
			PositionalParams: []any{url, 1},
		},
	}, nil)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	errored, _, errorMessage := response.HasError()

	if errored {
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	fmt.Printf("ExecuteResponse: %+v\n", response)

	return response, nil
}
