package main

import (
	"github.com/cyberlight/bbcrawler"
	"time"
	"fmt"
)

func main() {
	config := &bbcrawler.HackerOneCrawlerConfig{
		SearchUrl:     "URL",
		PathToLocalDb: "DB",
		FireBaseUrl:   "https://project_name.firebaseio.com",
		FireBaseToken: "TOKEN",
	}
	crawler := bbcrawler.NewHackerOneCrowler(config)

	go crawler.Crawl()
	done := false
	for {
		select {
		case <-crawler.Done:
			news := crawler.GetNewRecords().([]bbcrawler.HackerOneRecord)
			fmt.Println()
			fmt.Println("<=====New Records: ", len(news))
			crawler.ClearNewRecords()
			done = true
		case <-time.After(1 * time.Minute):
			if done {
				done = false
				go crawler.Crawl()
			}
		}
	}
}
