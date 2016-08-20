package main

import (
	"fmt"
	"github.com/cyberlight/bbcrawler"
	"time"
)

func main() {
	config := &bbcrawler.HackerOneCrawlerConfig{
		SearchUrl:           "URL",
		HacktivitySearchUrl: "URL",
		PathToLocalDb:       "DB",
		FireBaseUrl:         "URL",
		FireBaseToken:       "TOKEN",
	}
	crawler := bbcrawler.NewHackerOneCrowler(config)
	h1HackCrawler := bbcrawler.NewH1HacktivityCrowler(config)

	go crawler.Crawl()
	go h1HackCrawler.Crawl()

	doneCrawler := false
	doneH1Crawler := false
	for {
		select {
		case <-crawler.Done:
			news := crawler.GetNewRecords().([]bbcrawler.HackerOneRecord)
			fmt.Println()
			fmt.Println("<=====New Records: ", len(news))
			crawler.ClearNewRecords()
			doneCrawler = true
		case <-h1HackCrawler.Done:
			news := h1HackCrawler.GetNewRecords().([]bbcrawler.H1HactivityRecord)
			fmt.Println()
			fmt.Println("<=====New Records Hacktivity: ", len(news))
			h1HackCrawler.ClearNewRecords()
			doneH1Crawler = true
		case <-time.After(1 * time.Minute):
			if doneCrawler {
				doneCrawler = false
				go crawler.Crawl()
			}
			if doneH1Crawler {
				doneH1Crawler = false
				go h1HackCrawler.Crawl()
			}
		}
	}
}
