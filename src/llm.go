package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"parksideNotifier/src/interfaces"
)

func GetProductsWithOpenAI(imageUrls []string) []interfaces.Product {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: OPENAI_API_KEY environment variable not set")
		os.Exit(1)
	}

	// Create the request payload
	textContent := interfaces.OpenAITextContent{
		Type: "input_text",
		Text: "Analyze this flyer of products, if there're any products from brand Parkside. If no products are found, return an empty array.",
	}

	content := []interface{}{
		textContent,
	}

	for _, imageUrl := range imageUrls {
		content = append(content, interfaces.OpenAIImageContent{
			Type:     "input_image",
			ImageURL: imageUrl,
		})
	}

	input := interfaces.OpenAIInput{
		Role:    "user",
		Content: content,
	}

	request := interfaces.OpenAIRequest{
		Model:        "gpt-4o-mini",
		Instructions: "Return only a JSON array of all products. Each product should have a name and price. Expected response format: [{\"name\":\"Product 1\",\"price\":100}]. Make sure the output is a single string of the JSON array, not raw JSON or with the json word",
		Input:        []interfaces.OpenAIInput{input},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Error marshaling OpenAI request: %v\n", err)
		os.Exit(1)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating OpenAI request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending OpenAI request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading OpenAI response: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("OpenAI API error (status %d): %s\n", resp.StatusCode, string(bodyBytes))
		os.Exit(1)
	}

	// Parse response
	var openaiResponse interfaces.OpenAIResponse
	err = json.Unmarshal(bodyBytes, &openaiResponse)
	if err != nil {
		fmt.Printf("Error parsing OpenAI response: %v\n", err)
		os.Exit(1)

	}

	if openaiResponse.Error.Message != "" {
		fmt.Printf("OpenAI API error: %s\n", openaiResponse.Error.Message)
		os.Exit(1)
	}

	if len(openaiResponse.Output) == 0 {
		fmt.Println("No response from OpenAI")
		return []interfaces.Product{}
	}

	// Response sometimes comes as ```json'[{"name":"Parksidе Aspirador/ Soprador de Folhas Elétrico 2600 W","price":29.99}]'``` so we need to clean it
	jsonProducts := CleanJsonString(openaiResponse.Output[0].Content[0].Text)
	fmt.Printf("OpenAI Response: %s\n", jsonProducts)

	// Parse the JSON response
	var products []interfaces.Product
	err = json.Unmarshal([]byte(jsonProducts), &products)
	if err != nil {
		fmt.Printf("Error parsing product JSON from OpenAI: %v\n", err)
		os.Exit(1)
	}

	return products
}
