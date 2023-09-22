ckage main

import (
	"rrt-bot/bot"
	"log"
	"os"
)

func main() {
	botToken, ok := os.LookupEnv("BOOT_TOKEN")
