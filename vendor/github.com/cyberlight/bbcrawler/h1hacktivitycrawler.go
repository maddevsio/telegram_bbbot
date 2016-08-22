package bbcrawler

import (
	"fmt"
	"strconv"
	"sync"
)

const (
	FIREBASE_H1_HACKTIVITY_DATABASE = "h1hacktivity"
)

type H1HacktivityReporter struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Username string `json:"username"`
}

type H1HacktivityTeamProfile struct {
	Name string `json:"name"`
}

type H1HacktivityTeamProfilePict struct {
	Medium string `json:"medium"`
	Small  string `json:"small"`
}

type H1HacktivityTeam struct {
	Handle             string                      `json:"handle"`
	Url                string                      `json:"url"`
	Profile            H1HacktivityTeamProfile     `json:"profile"`
	ProfilePictureUrls H1HacktivityTeamProfilePict `json:"profile_picture_urls"`
}

type H1HactivityRecord struct {
	BountyDisclosed             bool                 `json:"bounty_disclosed"`
	FormattedBounty             string               `json:"formatted_bounty"`
	Id                          int                  `json:"id"`
	LatestDisclosableAction     string               `json:"latest_disclosable_action"`
	LatestDisclosableActivityAt string               `json:"latest_disclosable_activity_at"`
	ReadableSubstate            string               `json:"readable_substate"`
	Reporter                    H1HacktivityReporter `json:"reporter"`
	Substate                    string               `json:"substate"`
	Swag                        bool                 `json:"swag"`
	Team                        H1HacktivityTeam     `json:"team"`
	Title                       string               `json:"title"`
	Url                         string               `json:"url"`
	VoteCount                   int                  `json:"vote_count"`
	Voters                      []string             `json:"voters"`
}

type H1HactivityResponse struct {
	Count   int                 `json:"count"`
	Reports []H1HactivityRecord `json:"reports"`
}

type H1HacktivityCrawler struct {
	sync.RWMutex
	fetcher            Fetcher
	reader             Reader
	store              Storer
	fbSync             FireBaseSyncer
	Done               chan bool
	hackerOneSearchUrl string
	config             *HackerOneCrawlerConfig
}

func (h *H1HacktivityCrawler) Crawl() {
	fmt.Println("Check database consistancy Hacktivity")
	if empty, err := h.store.IsEmpty(); empty && err == nil {
		h.syncDb()
	} else if err != nil {
		fmt.Println("<===== Sync error: ", err)
	}

	data, err := h.fetcher.Fetch(h.config.HacktivitySearchUrl, h.makeQuery(1))
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonResponse, err := h.reader.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := jsonResponse.(H1HactivityResponse)
	err = h.store.Store(response)
	if err != nil {
		fmt.Println("Error storing Hacktivity: ", err)
	}
	newRecords := h.store.GetNewRecords().([]H1HactivityRecord)
	if len(newRecords) > 0 {
		fmt.Println("New records Hacktivity:", len(newRecords))
		for _, v := range newRecords {
			err := h.fbSync.Write(v, fmt.Sprintf("%s/%s", FIREBASE_H1_HACKTIVITY_DATABASE, strconv.Itoa(v.Id)))
			if err != nil {
				fmt.Println("Error firebase Hacktivity write: ", err)
			}
		}
	} else {
		fmt.Println()
		fmt.Println("No new records Hacktivity found!")
	}
	h.Done <- true
}

func (h H1HacktivityCrawler) ClearNewRecords() {
	h.store.Clear()
}

func (h H1HacktivityCrawler) GetNewRecords() interface{} {
	return h.store.GetNewRecords()
}

func (h *H1HacktivityCrawler) syncDb() {
	fmt.Println("Hactivity Syncing ...")
	var d map[string]H1HactivityRecord
	err := h.fbSync.Read(&d, FIREBASE_H1_HACKTIVITY_DATABASE)
	if err != nil {
		fmt.Println("Hacktivity Syncing error!", err)
		return
	}
	for k, v := range d {
		fmt.Println("Write key: ", k)
		err := h.store.Store(v)
		if err != nil {
			fmt.Println("Error write to local database HAcktivity: ", err)
		}
	}
}

func (h *H1HacktivityCrawler) makeQuery(pageNum int) map[string]string {
	hackerOneQuery := make(map[string]string)
	hackerOneQuery["sort_type"] = "latest_disclosable_activity_at"
	hackerOneQuery["filter"] = "type:public"
	hackerOneQuery["page"] = fmt.Sprintf("%d", pageNum)
	return hackerOneQuery
}

func NewH1HacktivityCrowler(config *HackerOneCrawlerConfig) *H1HacktivityCrawler {
	return &H1HacktivityCrawler{
		fetcher: hackerOneFetcher{},
		reader:  h1HacktivityParser{},
		store: &H1HacktivityStore{
			PathToDb:   config.PathToLocalDb,
			newRecords: make([]H1HactivityRecord, 0),
			all: make(map[string]H1HactivityRecord),
		},
		fbSync: FireBaseSync{
			Token:   config.FireBaseToken,
			BaseUrl: config.FireBaseUrl,
		},
		Done:               make(chan bool),
		hackerOneSearchUrl: config.SearchUrl,
		config:             config,
	}
}
