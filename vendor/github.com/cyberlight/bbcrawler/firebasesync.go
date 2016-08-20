package bbcrawler

import (
	"fmt"
	"github.com/melvinmt/firebase"
	"net/url"
)

type FireBaseSync struct {
	Token   string
	BaseUrl string
}

func (f FireBaseSync) Write(data interface{}, path string) error {
	queryUrl, err := f.getFireBaseDataUrl(path)
	if err != nil {
		return err
	}
	fmt.Println("QuertUrl: ", queryUrl)
	ref := firebase.NewReference(queryUrl).Auth(f.Token)
	if err = ref.Write(data); err != nil {
		return err
	}
	return nil
}

func (f FireBaseSync) Read(data interface{}, path string) error {
	queryUrl, err := f.getFireBaseDataUrl(path)
	if err != nil {
		return err
	}
	dataRef := firebase.NewReference(queryUrl).Export(false)
	if err = dataRef.Value(data); err != nil {
		return err
	}
	return nil

}

func (f FireBaseSync) getFireBaseDataUrl(path string) (string, error) {
	u, err := url.Parse(f.BaseUrl)
	if err != nil {
		return "", err
	}
	u.Path = path
	return u.String(), nil
}
