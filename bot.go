package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/cyberlight/bbcrawler"
	"github.com/gin-gonic/gin"
	"gopkg.in/telegram-bot-api.v4"
)

type config struct {
	Token                 string `env:"TELEGRAM_BBBOT_TOKEN"`
	Url                   string `env:"TELEGRAM_BBBOT_URL"`
	Port                  string `env:"PORT"`
	FireBaseToken         string `env:"TELEGRAM_BBBOT_FIREBASE_TOKEN"`
	FireBaseUrl           string `env:"TELEGRAM_BBBOT_FIREBASE_URL"`
	PathToLocalDb         string `env:"TELEGRAM_BBBOT_PATH_TO_LOCAL_DB"`
	HOSearchUrl           string `env:"TELEGRAM_BBBOT_HO_SEARCH_URL"`
	BotChannel            string `env:"TELEGRAM_BBBOT_CHANNEL"`
	PingHost              string `env:"TELEGRAM_BBBOT_HOST"`
	H1HackSearchUrl       string `env:"TELEGRAM_BBBOT_H1_HACK_SEARCH_URL"`
	BugCrowdNewProgramUrl string `env:"TELEGRAM_BBBOT_BUGCROWD_NEW_PROG_URL"`
	HackerOneBaseUrl      string `env:"TELEGRAM_BBBOT_H1_BASE_URL"`
}

var (
	ConfigurationError  = fmt.Errorf("Configuration error. Please check setting of ENV variables: %s", "TELEGRAM_BBBOT_TOKEN")
	TelegramBotApiError = func(err error) error { return fmt.Errorf("Telegram Bot API Error: %s", err.Error()) }
)

func botReceiveUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("%+v\n", update)
}

func initRouting(bot *tgbotapi.BotAPI, cfg config, updateChan chan tgbotapi.Update) {
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	handleWebHook := func(c *gin.Context) {
		var update tgbotapi.Update
		err := c.BindJSON(&update)
		if err != nil {
			log.Println(err)
		} else {
			updateChan <- update
		}
	}

	router.GET("/"+bot.Token, handleWebHook)
	router.POST("/"+bot.Token, handleWebHook)

	go router.Run(":" + cfg.Port)
}

func main() {
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(ConfigurationError)
	}

	hoConfig := &bbcrawler.HackerOneCrawlerConfig{
		FireBaseToken:       cfg.FireBaseToken,
		FireBaseUrl:         cfg.FireBaseUrl,
		PathToLocalDb:       cfg.PathToLocalDb,
		SearchUrl:           cfg.HOSearchUrl,
		HacktivitySearchUrl: cfg.H1HackSearchUrl,
		BugCrowdProgramsUrl: cfg.BugCrowdNewProgramUrl,
	}

	hoCrawler := bbcrawler.NewHackerOneCrowler(hoConfig)
	h1HackCrawler := bbcrawler.NewH1HacktivityCrowler(hoConfig)
	bugCrowdNewProgCrawler := bbcrawler.NewBugCrowdCrawler(hoConfig)

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

	updatesChan := make(chan tgbotapi.Update, 100)
	initRouting(bot, cfg, updatesChan)
	go hoCrawler.Crawl()
	go h1HackCrawler.Crawl()
	go bugCrowdNewProgCrawler.Crawl()

	doneH1Crawler := false
	doneH1Hacktivity := false
	doneBCNewProg := false

	for {
		select {
		case update := <-updatesChan:
			fmt.Println("<=== Received Updates from telegram")
			botReceiveUpdate(bot, update)
		case <-hoCrawler.Done:
			fmt.Println("<=== Received signal Done: Hackerone")
			records := hoCrawler.GetNewRecords().([]bbcrawler.HackerOneRecord)
			if len(records) > 0 {
				for _, v := range records {
					msg := tgbotapi.NewMessageToChannel(cfg.BotChannel,
						fmt.Sprintf(
							"*%s*\n\n"+
								"```text \n"+
								"%s"+
								"```\n"+
								"%s\n",
							v.Handle,
							v.StrippedPolicy,
							fmt.Sprintf(cfg.HackerOneBaseUrl+"%s", v.Url),
						))
					msg.ParseMode = "Markdown"
					bot.Send(msg)
				}
			}
			hoCrawler.ClearNewRecords()
			fmt.Println("<=== Clear records: Hackerone")
			doneH1Crawler = true
		case <-h1HackCrawler.Done:
			fmt.Println("<=== Received signal Done: Hackerone Hactivity")
			records := h1HackCrawler.GetNewRecords().([]bbcrawler.H1HactivityRecord)
			if len(records) > 0 {
				for _, v := range records {
					msg := tgbotapi.NewMessageToChannel(cfg.BotChannel,
						fmt.Sprintf(
							"_Hacktivity_ from *%s*\n"+
								"```text \n"+
								"%s\n"+
								"```\n"+
								"%s\n",
							v.Reporter.Username,
							v.Title,
							fmt.Sprintf(cfg.HackerOneBaseUrl+"%s", v.Url),
						))
					msg.ParseMode = "Markdown"
					bot.Send(msg)
				}
			}
			h1HackCrawler.ClearNewRecords()
			fmt.Println("<=== Clear records: Hackerone Hactivity")
			doneH1Hacktivity = true
		case <-bugCrowdNewProgCrawler.Done:
			fmt.Println("<=== Received signal Done: CrowdCom")
			records := bugCrowdNewProgCrawler.GetNewRecords().([]bbcrawler.BugCrowdNewProgramsRecord)
			if len(records) > 0 {
				for _, v := range records {
					msg := tgbotapi.NewMessageToChannel(cfg.BotChannel,
						fmt.Sprintf(
							"\n_Bugcrowd.com_ new program *%s*\n"+
								"%s\n",
							v.Name,
							v.Link,
						))
					msg.ParseMode = "Markdown"
					bot.Send(msg)
				}
			}
			bugCrowdNewProgCrawler.ClearNewRecords()
			fmt.Println("<=== Clear records: CrowdCom")
			doneBCNewProg = true
		case <-time.After(2 * time.Minute):
			fmt.Println("== Ping request ==")
			c := http.Client{
				Timeout: 10 * time.Second,
			}
			c.Get(cfg.PingHost)

			fmt.Println("<=== doneH1Crawler: ", doneH1Crawler)
			if doneH1Crawler {
				doneH1Crawler = false
				go hoCrawler.Crawl()
			}
			fmt.Println("<=== doneH1Hacktivity: ", doneH1Hacktivity)
			if doneH1Hacktivity {
				doneH1Hacktivity = false
				go h1HackCrawler.Crawl()
			}
			fmt.Println("<=== doneBCNewProg: ", doneBCNewProg)
			if doneBCNewProg {
				doneBCNewProg = false
				go bugCrowdNewProgCrawler.Crawl()
			}
		}
	}
}
