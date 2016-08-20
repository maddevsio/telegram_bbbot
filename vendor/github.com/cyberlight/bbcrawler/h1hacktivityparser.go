package bbcrawler

import (
	"encoding/json"
)

type h1HacktivityParser struct{}

func (h h1HacktivityParser) Read(data []byte) (interface{}, error) {
	var jsonResponse H1HactivityResponse
	err := json.Unmarshal(data, &jsonResponse)
	return jsonResponse, err
}
