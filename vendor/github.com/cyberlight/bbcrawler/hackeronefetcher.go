package bbcrawler

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type hackerOneFetcher struct{}

func (ho hackerOneFetcher) Fetch(fetchUrl string, params map[string]string) ([]byte, error) {
	client := &http.Client{}
	u, err := url.Parse(fetchUrl)
	if err != nil {
		return []byte{}, err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36")
	req.Header.Set("referer", "https://hackerone.com/")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	res, _ := client.Do(req)
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
