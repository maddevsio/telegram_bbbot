package bbcrawler

import (
	"fmt"
	"os"
	"sync"
)

const (
	FIREBASE_BUGCROWD_NEW_DATABASE = "bugcrowd_new"
)

type BugCrowdCrawler struct {
	sync.RWMutex
	fetcher UniversalFetcher
	store   Storer
	fbSync  FireBaseSyncer
	Done    chan bool
	config  *HackerOneCrawlerConfig
}

func (b *BugCrowdCrawler) Crawl() {
	fmt.Println("Check database consistancy Bugcrowd")
	if _, err := os.Stat(b.config.PathToLocalDb); os.IsNotExist(err) {
		b.syncDb()
	} else if empty, err := b.store.IsEmpty(); empty && err == nil {
		b.syncDb()
	} else if err != nil {
		fmt.Println("<===== Sync Bug crowd error: ", err)
	}

	data, err := b.fetcher.Fetch(b.config.BugCrowdProgramsUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	records := data.([]BugCrowdNewProgramsRecord)
	err = b.store.Store(records)
	if err != nil {
		fmt.Println("Error storing Bugcrowd: ", err)
	}

	newRecords := b.store.GetNewRecords().([]BugCrowdNewProgramsRecord)
	if len(newRecords) > 0 {
		fmt.Println()
		fmt.Println("New records Hacktivity:", len(newRecords))
		for _, v := range newRecords {
			err := b.fbSync.Write(v, fmt.Sprintf("%s/%s", FIREBASE_BUGCROWD_NEW_DATABASE, v.Name))
			if err != nil {
				fmt.Println("Error firebase Bugcrowd write: ", err)
			}
		}
	} else {
		fmt.Println()
		fmt.Println("No new records Bugcrowd found!")
	}
	b.Done <- true
}

func (h *BugCrowdCrawler) syncDb() {
	fmt.Println("BugCrowd new Syncing ...")
	var d map[string]BugCrowdNewProgramsRecord
	err := h.fbSync.Read(&d, FIREBASE_BUGCROWD_NEW_DATABASE)
	if err != nil {
		fmt.Println("BugCrowd new error!", err)
		return
	}
	for k, v := range d {
		fmt.Println("BugCrowd -> Write key: ", k)
		err := h.store.Store(v)
		if err != nil {
			fmt.Println("Error write to local database Bugcrowd new: ", err)
		}
	}
	h.store.Clear()
}

func (h BugCrowdCrawler) ClearNewRecords() {
	h.store.Clear()
}

func (h BugCrowdCrawler) GetNewRecords() interface{} {
	return h.store.GetNewRecords()
}

func NewBugCrowdCrawler(config *HackerOneCrawlerConfig) *BugCrowdCrawler {
	return &BugCrowdCrawler{
		fetcher: bugCrowdFetcher{},
		store: &BugCrowdStore{
			PathToDb:   config.PathToLocalDb,
			newRecords: make([]BugCrowdNewProgramsRecord, 0),
			all: make(map[string]BugCrowdNewProgramsRecord),
		},
		fbSync: FireBaseSync{
			Token:   config.FireBaseToken,
			BaseUrl: config.FireBaseUrl,
		},
		Done:   make(chan bool),
		config: config,
	}
}
