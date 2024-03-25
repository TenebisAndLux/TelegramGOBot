package TelegramGOBot

import (
	"TelegramGOBot/clients/telegram"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient = telegram.New(mustToken())

	//fetcher = fetcher.New(tgClient)

	//processor = processor.New(tgClient)

	//consumer.Start(fetcher, processor)
}

func mustToken(string, error) {
	token := flag.String(
		"token-bot-token",
		"",
		"Token for access to telegram bot",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("Token is not avaleble")
	}
}
