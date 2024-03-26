package main

import (
	event_consumer "TelegramGOBot/consumer/event-consumer"
	"TelegramGOBot/events/telegram"
	"flag"
	"log"

	tgClient "TelegramGOBot/clients/telegram"
	"TelegramGOBot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Printf("Service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("Service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"Token for access to telegram bot",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("Token is not avaleble")
	}
	return *token
}
