package main

import (
	"SeparatorMusicBot/pkg/telegram"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Load .env file error", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	mainBot, err := telegram.NewBot(token)
	if err != nil {
		log.Fatal("Invalid bot token", err)
	}

	mainBot.Start()
}
