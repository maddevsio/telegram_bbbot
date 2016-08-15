package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
)

type config struct {
	Token string `env:"TELEGRAM_BBBOT_TOKEN"`
	Url   string `env:"TELEGRAM_BBBOT_URL"`
	Port  string `env:"PORT"`
}

var (
	ConfigurationError  = fmt.Errorf("Configuration error. Please check setting of ENV variables: %s", "TELEGRAM_BBBOT_TOKEN")
	TelegramBotApiError = func(err error) error { return fmt.Errorf("Telegram Bot API Error: %s", err.Error()) }
)

func main() {
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(ConfigurationError)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(TelegramBotApiError(err))
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(cfg.Url))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	go router.Run(":" + cfg.Port)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
