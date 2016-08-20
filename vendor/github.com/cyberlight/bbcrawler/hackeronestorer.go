package bbcrawler

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"sync"
	"time"
)

type HackerOneStore struct {
	PathToDb   string
	newRecords []HackerOneRecord
	sync.RWMutex
}

func (h *HackerOneStore) IsEmpty() (bool, error) {
	ErrorEmptyDb := fmt.Errorf("No bucket found %s", "All")
	fmt.Println("Store running")
	h.Lock()
	defer h.Unlock()
	db, err := bolt.Open(h.PathToDb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()

	if err != nil {
		return false, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("All"))
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

func (h *HackerOneStore) Store(data interface{}) error {
	fmt.Println("Store running")
	h.Lock()
	defer h.Unlock()
	db, err := bolt.Open(h.PathToDb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()

	if err != nil {
		return err
	}

	if response, ok := data.(HackerOneResponse); ok {
		for _, v := range response.Results {
			fmt.Print(".")
			jsonStr, _ := json.Marshal(v)
			db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte("All"))
				if err != nil {
					return fmt.Errorf("create All bucket: %s", err)
				}
				if b.Get([]byte(v.Handle)) == nil {
					fmt.Print("+")
					h.newRecords = append(h.newRecords, v)
					return b.Put([]byte(v.Handle), jsonStr)
				}
				return nil
			})
		}
	} else if rec, ok := data.(HackerOneRecord); ok {
		jsonStr, _ := json.Marshal(rec)
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("All"))
			if err != nil {
				return fmt.Errorf("create All bucket: %s", err)
			}
			if b.Get([]byte(rec.Handle)) == nil {
				return b.Put([]byte(rec.Handle), jsonStr)
			}
			return nil
		})
	} else if !ok {
		fmt.Println("Fail converting to HackerOneResponse")
	}
	return nil
}

func (h HackerOneStore) GetNewRecords() interface{} {
	h.RLock()
	defer h.RUnlock()
	return h.newRecords
}

func (h *HackerOneStore) Clear() {
	h.Lock()
	defer h.Unlock()
	h.newRecords = nil
}
