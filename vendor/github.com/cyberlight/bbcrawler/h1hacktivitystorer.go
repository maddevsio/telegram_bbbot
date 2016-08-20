package bbcrawler

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"strconv"
	"sync"
	"time"
)

const (
	HACKTIVITY_BACKET = "Hacktivity"
)

type H1HacktivityStore struct {
	PathToDb   string
	newRecords []H1HactivityRecord
	sync.RWMutex
}

func (h *H1HacktivityStore) IsEmpty() (bool, error) {
	ErrorEmptyDb := fmt.Errorf("No bucket found %s", HACKTIVITY_BACKET)
	fmt.Println("Store Hacktivity running")
	h.Lock()
	defer h.Unlock()
	db, err := bolt.Open(h.PathToDb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()

	if err != nil {
		return false, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(HACKTIVITY_BACKET))
		if b == nil {
			return ErrorEmptyDb
		}
		return nil
	})

	if err != nil && err == ErrorEmptyDb {
		return true, nil
	}
	return false, err
}

func (h *H1HacktivityStore) Store(data interface{}) error {
	fmt.Println("Store Hacktivity running")
	h.Lock()
	defer h.Unlock()
	db, err := bolt.Open(h.PathToDb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()

	if err != nil {
		return err
	}

	if response, ok := data.(H1HactivityResponse); ok {
		for _, v := range response.Reports {
			fmt.Print("h.")
			jsonStr, _ := json.Marshal(v)
			db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte(HACKTIVITY_BACKET))
				if err != nil {
					return fmt.Errorf("create Hacktivity bucket: %s", err)
				}
				if b.Get([]byte(strconv.Itoa(v.Id))) == nil {
					fmt.Print("h+")
					h.newRecords = append(h.newRecords, v)
					return b.Put([]byte(strconv.Itoa(v.Id)), jsonStr)
				}
				return nil
			})
		}
	} else if rec, ok := data.(H1HactivityRecord); ok {
		jsonStr, _ := json.Marshal(rec)
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(HACKTIVITY_BACKET))
			if err != nil {
				return fmt.Errorf("create Hacktivity bucket: %s", err)
			}
			if b.Get([]byte(strconv.Itoa(rec.Id))) == nil {
				return b.Put([]byte(strconv.Itoa(rec.Id)), jsonStr)
			}
			return nil
		})
	} else if !ok {
		fmt.Println("Fail converting to H1HactivityResponse or H1ActivityRecord")
	}
	return nil
}

func (h H1HacktivityStore) GetNewRecords() interface{} {
	h.RLock()
	defer h.RUnlock()
	return h.newRecords
}

func (h *H1HacktivityStore) Clear() {
	h.Lock()
	defer h.Unlock()
	h.newRecords = nil
}
