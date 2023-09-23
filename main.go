package main

import (
	"rrt-bot/bot"
	"log"
	"os"
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
	
	bot.BotToken = botToken
	bot.BotMapDir = botMapDir
	bot.Run()
}
