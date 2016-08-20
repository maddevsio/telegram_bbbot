package bbcrawler

import "encoding/json"

type HackerOneBounty struct {
	BugCount      int `json:"bug_count"`
	MinimumBounty int `json:"minimum_bounty"`
}

type HackerOneRecord struct {
	Id                int             `json:"id"`
	Url               string          `json:"url"`
	Name              string          `json:"name"`
	Meta              HackerOneBounty `json:"meta"`
	About             string          `json:"about"`
	StrippedPolicy    string          `json:"stripped_policy"`
	Handle            string          `json:"handle"`
	ProfilePicture    string          `json:"profile_picture"`
	InternetBugBounty bool            `json:"internet_bug_bounty"`
}

type HackerOneResponse struct {
	Limit   int `json:"limit"`
	Total   int `json:"total"`
	Results []HackerOneRecord
}

type hackerOneParser struct{}

func (h hackerOneParser) Read(data []byte) (interface{}, error) {
	var jsonResponse HackerOneResponse
	err := json.Unmarshal(data, &jsonResponse)
	return jsonResponse, err
}
