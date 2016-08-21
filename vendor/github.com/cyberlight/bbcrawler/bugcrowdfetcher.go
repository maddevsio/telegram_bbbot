package bbcrawler

import (
	"github.com/PuerkitoBio/goquery"
)

type bugCrowdFetcher struct{}

type BugCrowdNewProgramsRecord struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func (ho bugCrowdFetcher) Fetch(fetchUrl string) (interface{}, error) {
	result := []BugCrowdNewProgramsRecord{}

	doc, err := goquery.NewDocument(fetchUrl)
	if err != nil {
		return nil, err
	}

	doc.Find("table.table tr.new").Each(func(i int, s *goquery.Selection) {
		ahref := s.Find("td a.tracked")
		link, exists := ahref.Attr("href")
		if exists {
			result = append(result, BugCrowdNewProgramsRecord{
				Name: ahref.Text(),
				Link: link,
			})
		}
	})

	return result, nil
}


