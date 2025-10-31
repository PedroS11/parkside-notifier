package main

import (
	"strings"
)

func CleanJsonString(json string) string {
	cleanedString := strings.TrimPrefix(json, "```json")

	cleanedString = strings.TrimSuffix(cleanedString, "```")

	return strings.Trim(cleanedString, "'")
}
