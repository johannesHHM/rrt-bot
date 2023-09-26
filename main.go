package main

import (
	"log"
	"os"
	"rrt-bot/bot"
)

func main() {
	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("Env variable BOT_TOKEN not found")
	}
	botMapDir, ok := os.LookupEnv("BOT_MAP_DIR")
	if !ok {
		log.Fatal("Env variable BOT_MAP_DIR not found")
	}
	botPrefix, ok := os.LookupEnv("BOT_PREFIX")
	if !ok {
		log.Fatal("Env variable BOT_PREFIX not found")
	}
	botURLHTTPMaster, ok := os.LookupEnv("BOT_URL_HTTP_MASTER")
	if !ok {
		log.Fatal("Env variable BOT_URL_HTTP_MASTER not found")
	}
	bot.BotToken = botToken
	bot.BotMapDir = botMapDir
	bot.BotPrefix = botPrefix
	bot.URLHttpMaster = botURLHTTPMaster
	bot.BotTranslatePrefix = ':'
	bot.Run()
}
