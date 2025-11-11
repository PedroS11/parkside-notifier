package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"parksideNotifier/src/interfaces"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func GetProductsFromUrls(imageUrls []string) ([]interfaces.Product, error) {
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	ctx := context.Background()

	question := "Analyze these images from flyers of products, if there're any products from brand Parkside return them. If no products are found, return an empty array"

	var inputParams responses.ResponseInputParam

	// Create content slice starting with text input
	content := responses.ResponseInputMessageContentListParam{
		{
			OfInputText: &responses.ResponseInputTextParam{
				Text: question,
			},
		},
	}

	var imagesInputContent responses.ResponseInputMessageContentListParam

	for _, imageUrl := range imageUrls {
		imagesInputContent = append(imagesInputContent, responses.ResponseInputContentUnionParam{
			OfInputImage: &responses.ResponseInputImageParam{
				ImageURL: openai.String(imageUrl),
			},
		})
	}

	// Append image content
	content = append(content, imagesInputContent...)

	inputParams = append(inputParams, responses.ResponseInputItemUnionParam{
		OfInputMessage: &responses.ResponseInputItemMessageParam{
			Role:    "user",
			Content: content,
		},
	})

	slog.Info(fmt.Sprintf("Calling OpenAI for %d images with input: %v\n", len(imageUrls), inputParams))

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input:        responses.ResponseNewParamsInputUnion{OfInputItemList: inputParams},
		Model:        openai.ChatModelGPT4oMini,
		Instructions: openai.String("Return only a JSON array of all products. Each product should have a name and price. Expected response format: [{\"name\":\"Product 1\",\"price\":100}]. Make sure the output is a single string of the JSON array, not raw JSON or with the json word"),
	})

	if err != nil {
		slog.Error(fmt.Sprintf("Error sending request to OpenAI: %v\n", err.Error()))
		return []interfaces.Product{}, err
	}

	// Response sometimes comes as ```json'[{"name":"Parksidе Aspirador/ Soprador de Folhas Elétrico 2600 W","price":29.99}]'``` so we need to clean it
	jsonProducts := CleanJsonString(resp.OutputText())
	fmt.Printf("OpenAI Response: %s\n", jsonProducts)

	// Parse the JSON response
	var products []interfaces.Product
	err = json.Unmarshal([]byte(jsonProducts), &products)
	if err != nil {
		slog.Error(fmt.Sprintf("Error parsing product JSON %v from OpenAI: %v\n", jsonProducts, err))
		return []interfaces.Product{}, err
	}

	return products, nil
}
