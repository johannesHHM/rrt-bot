package bot

import (
	"log"
	"strings"

	translator "github.com/Conight/go-googletrans"
)

var (
	BotTranslatePrefix byte
)

func parseContent(content string) (string, from string, dest string) {
	content = strings.ReplaceAll(content, BotPrefix+"gt", "")
	dest = "en"
	from = "auto"
	words := strings.Fields(content)
	for _, word := range words {
		if word[0] == BotTranslatePrefix {
			dest = word[1:]
			content = strings.ReplaceAll(content, word, "")
			break
		}
	}
	return content, from, dest
}

func translateMessage(content string) (result *translator.Translated) {
	message, from, dest := parseContent(content)

	t := translator.New()
	result, err := t.Translate(message, from, dest)
	if err != nil {
		log.Panicln("Failed to translate message")
	}
	return result
}
