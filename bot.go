package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
	"strings"
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


	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	updatesChan := make(chan tgbotapi.Update, 100)

	handleWebHook := func(c *gin.Context) {
		var update tgbotapi.Update
		err := c.BindJSON(&update)
		if err != nil {
			log.Println(err)
		} else {
			updatesChan <- update
		}
	}

	router.GET("/" + bot.Token, handleWebHook)
	router.POST("/" + bot.Token, handleWebHook)

	go router.Run(":" + cfg.Port)

	for update := range updatesChan {
		log.Printf("%+v\n", update)

		if update.Message == nil {
			continue
		}

		if update.Message.Chat.IsGroup()  {
			u, err := bot.GetMe()
			if err != nil {
				log.Printf("Error: %s", err.Error())
			}else {
				if update.Message.NewChatMember != nil {
					if u.UserName != update.Message.NewChatMember.UserName {
						_, err := bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping))
						if err != nil {
							log.Printf("Error: %s", err.Error())
						}

						hiText := fmt.Sprintf("Привет! %s Я баг баунти бот!",
							update.Message.NewChatMember.FirstName)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, hiText)

						bot.Send(msg)

						count, err := bot.GetChatMembersCount(update.Message.Chat.ChatConfig())
						if err == nil {
							hiText = fmt.Sprintf("Нас уже %d !", count)
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, hiText)
							bot.Send(msg)
						}
					} else {
						hiText := "Всем привет! Я баг баунти бот!"
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, hiText)
						bot.Send(msg)
					}
				}
			}
		}

		if bot.IsMessageToMe(*update.Message) {
			_, err := bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping))
			if err != nil {
				log.Printf("Error: %s", err.Error())
			}
			if strings.Contains(update.Message.Text, "что нового") {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пока ничего нового! :)")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}
}
