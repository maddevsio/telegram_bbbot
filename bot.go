package main

import (
	"github.com/caarlos0/env"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	//"gopkg.in/telegram-bot-api.v4"
	"log"
	"fmt"
)

type config struct {
	Token string `env:"TELEGRAM_BBBOT_TOKEN"`
}

var (
	ConfigurationError =
		fmt.Errorf("Configuration error. Please check setting of ENV variables: %s", "TELEGRAM_BBBOT_TOKEN")
	TelegramBotApiError = func (err error) error { return fmt.Errorf("Telegram Bot API Error: %s", err.Error()) }
)

func main() {
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(ConfigurationError)
	}
	log.Println("TOKEN: ", cfg.Token)
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(TelegramBotApiError(err))
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
