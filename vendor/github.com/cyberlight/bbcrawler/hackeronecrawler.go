package bbcrawler

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"sync"
)

var (
	wg sync.WaitGroup
)

type HackerOneCrawler struct {
	sync.RWMutex
	fetcher            Fetcher
	reader             Reader
	store              Storer
	fbSync             FireBaseSyncer
	pages              map[int]*HackerOneResponse
	Done               chan bool
	hackerOneSearchUrl string
	config             *HackerOneCrawlerConfig
}

func (h *HackerOneCrawler) hackerOneCrawl(url string, queryParams map[string]string) {
	response, err := h.hackerOneCrawlPage(url, queryParams, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	pages := response.Total / response.Limit
	if response.Total%response.Limit > 0 {
		pages += 1
	}

	var urlparams []map[string]string
	for i := 2; i <= pages; i++ {
		hackerOneQuery := h.makeQuery(i)
		urlparams = append(urlparams, hackerOneQuery)

		wg.Add(1)
		go h.hackerOneCrawlPage(
			h.hackerOneSearchUrl,
			hackerOneQuery, i)
	}
	return
}

func (h *HackerOneCrawler) hackerOneCrawlPage(url string, queryParams map[string]string, page int) (*HackerOneResponse, error) {
	defer wg.Done()
	fmt.Println("Crawl for page ", page)
	data, err := h.fetcher.Fetch(url, queryParams)
	if err != nil {
		return nil, err
	}
	jsonResponse, err := h.reader.Read(data)
	if err != nil {
		return nil, err
	}

	response := jsonResponse.(HackerOneResponse)

	fmt.Println("Page ", page, " count: ", len(response.Results))

	h.Lock()
	h.pages[page] = &response
	h.Unlock()

	return &response, nil
}

func (h *HackerOneCrawler) makeQuery(pageNum int) map[string]string {
	hackerOneQuery := make(map[string]string)
	hackerOneQuery["query"] = "bounties:yes" //"bounties:yes ibb:false"
	hackerOneQuery["sort"] = "published_at:descending"
	hackerOneQuery["page"] = fmt.Sprintf("%d", pageNum)
	return hackerOneQuery
}

func (h *HackerOneCrawler) syncDb() {
	fmt.Println("File desn't exist! Syncing ...")
	var d map[string]HackerOneRecord
	err := h.fbSync.Read(&d, "hackerone")
	if err != nil {
		fmt.Println("Syncing error!", err)
		return
	}
	for k, v := range d {
		fmt.Println("Write key: ", k)
		err := h.store.Store(v)
		if err != nil {
			log.Println("Error write to local database: ", err)
		}
	}
}

func (h *HackerOneCrawler) Crawl() {
	fmt.Println("Check database consistancy")
	if _, err := os.Stat(h.config.PathToLocalDb); os.IsNotExist(err) {
		h.syncDb()
	} else if empty, err := h.store.IsEmpty(); empty && err == nil {
		h.syncDb()
	}

	wg.Add(1)
	fmt.Println()
	fmt.Println("Begin Hacker one crawl")
	h.hackerOneCrawl(
		h.hackerOneSearchUrl,
		h.makeQuery(1))
	wg.Wait()

	fmt.Println("Pages retrieved: ", len(h.pages))
	for i := 1; i <= len(h.pages); i++ {
		response := h.pages[i]
		fmt.Println()
		fmt.Println("> Storing result for page: ", i, len(response.Results), response.Limit, response.Total)
		h.store.Store(*response)
	}

	newRecords := h.store.GetNewRecords().([]HackerOneRecord)
	if len(newRecords) > 0 {
		fmt.Println("New records:", len(newRecords))
		for _, v := range newRecords {
			err := h.fbSync.Write(v, fmt.Sprintf("hackerone/%s", v.Handle))
			if err != nil {
				fmt.Println("Error firebase write: ", err)
			}
		}
	} else {
		fmt.Println()
		fmt.Println("No new records found!")
	}
	h.Done <- true
}

func (h HackerOneCrawler) ClearNewRecords() {
	h.store.Clear()
}

func (h HackerOneCrawler) GetNewRecords() interface{} {
	return h.store.GetNewRecords()
}

func NewHackerOneCrowler(config *HackerOneCrawlerConfig) *HackerOneCrawler {
	return &HackerOneCrawler{
		fetcher: hackerOneFetcher{},
		reader:  hackerOneParser{},
		store: &HackerOneStore{
			PathToDb:   config.PathToLocalDb,
			newRecords: make([]HackerOneRecord, 0),
		},
		fbSync: FireBaseSync{
			Token:   config.FireBaseToken,
			BaseUrl: config.FireBaseUrl,
		},
		pages:              make(map[int]*HackerOneResponse),
		Done:               make(chan bool),
		hackerOneSearchUrl: config.SearchUrl,
		config:             config,
	}
}

type HackerOneCrawlerConfig struct {
	SearchUrl           string
	PathToLocalDb       string
	FireBaseUrl         string
	FireBaseToken       string
	HacktivitySearchUrl string
	BugCrowdProgramsUrl string
}
