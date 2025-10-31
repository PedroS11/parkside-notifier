package main

import (
	"strings"
)

func CleanJsonString(json string) string {
	cleanedString := strings.TrimPrefix(json, "```json")

	cleanedString = strings.TrimSuffix(cleanedString, "```")

	return strings.Trim(cleanedString, "'")
}

func EscapeMarkdownV2(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(s)
}
