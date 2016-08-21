package main

import (
	"fmt"
	"github.com/cyberlight/bbcrawler"
	"time"
)


func main() {
	config := &bbcrawler.HackerOneCrawlerConfig{
		SearchUrl:     "URL",
		HacktivitySearchUrl: "URL",
		PathToLocalDb: "DB",
		FireBaseUrl:   "URL",
		FireBaseToken: "TOKEN",
		BugCrowdProgramsUrl: "URL",
	}
	crawler := bbcrawler.NewHackerOneCrowler(config)
	h1HackCrawler := bbcrawler.NewH1HacktivityCrowler(config)
	bugCrowdNPCrawler := bbcrawler.NewBugCrowdCrawler(config)

	go crawler.Crawl()
	go h1HackCrawler.Crawl()
	go bugCrowdNPCrawler.Crawl()

	doneCrawler := false
	doneH1Crawler := false
	doneBgNewProgCrawler := false

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
		case <-bugCrowdNPCrawler.Done:
			news := bugCrowdNPCrawler.GetNewRecords().([]bbcrawler.BugCrowdNewProgramsRecord)
			fmt.Println()
			fmt.Println("<=====New Records BugCrowd: ", len(news))
			bugCrowdNPCrawler.ClearNewRecords()
			doneBgNewProgCrawler = true
		case <-time.After(1 * time.Minute):
			if doneCrawler {
				doneCrawler = false
				go crawler.Crawl()
			}
			if doneH1Crawler {
				doneH1Crawler = false
				go h1HackCrawler.Crawl()
			}
			if doneBgNewProgCrawler {
				doneBgNewProgCrawler = false
				go bugCrowdNPCrawler.Crawl()
			}
		}
	}
}
