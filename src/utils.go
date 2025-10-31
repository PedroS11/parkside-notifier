package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"parksideNotifier/src/interfaces"
	"strings"
)

func ImageToBase64(path string) string {
	resp, err := http.Get(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}

func FilterParksideProducts(products []interfaces.Product) []interfaces.Product {
	var parksideProducts []interfaces.Product

	for _, product := range products {
		if strings.Contains(strings.ToLower(product.Name), "parkside") {
			parksideProducts = append(parksideProducts, product)
		}
	}

	return parksideProducts
}
