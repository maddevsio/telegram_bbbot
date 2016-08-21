package bbcrawler

import (
	"sync"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
	"encoding/json"
)

const (
	BUGCROWD_NEW_BUCKET = "BugcrowdNew"
)

type BugCrowdStore struct {
	PathToDb   string
	newRecords []BugCrowdNewProgramsRecord
	sync.RWMutex
}

func (h *BugCrowdStore) IsEmpty() (bool, error) {
	ErrorEmptyDb := fmt.Errorf("No bucket found %s", BUGCROWD_NEW_BUCKET)
	h.Lock()
	defer h.Unlock()
	db, err := bolt.Open(h.PathToDb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()

	if err != nil {
		return false, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUGCROWD_NEW_BUCKET))
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

func (h *BugCrowdStore) Store(data interface{}) error {
	fmt.Println("Store ", BUGCROWD_NEW_BUCKET, " running")
	h.Lock()
	defer h.Unlock()
	db, err := bolt.Open(h.PathToDb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()

	if err != nil {
		return err
	}

	if recs, ok := data.([]BugCrowdNewProgramsRecord); ok {
		for _, rec := range recs {
			jsonStr, _ := json.Marshal(rec)
			fmt.Print(".")
			err = db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte(BUGCROWD_NEW_BUCKET))
				if err != nil {
					return fmt.Errorf("create BugcrowdNew bucket: %s", err)
					return err
				}
				if b.Get([]byte(rec.Name)) == nil {
					fmt.Print(".+")
					h.newRecords = append(h.newRecords, rec)
					return b.Put([]byte(rec.Name), jsonStr)
				}
				return nil
			})
			if err != nil {
				fmt.Println("Error check BugCrowd new program ", recs)
				return err
			}
		}
		return nil
	} else if rec, ok := data.(BugCrowdNewProgramsRecord); ok {
		jsonStr, _ := json.Marshal(rec)
		return db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(BUGCROWD_NEW_BUCKET))
			if err != nil {
				return fmt.Errorf("create Hacktivity bucket: %s", err)
			}
			if b.Get([]byte(rec.Name)) == nil {
				fmt.Print(".+")
				h.newRecords = append(h.newRecords, rec)
				return b.Put([]byte(rec.Name), jsonStr)
			}
			return nil
		})
	} else {
		return fmt.Errorf("Can't convert to type %s", "BugCrowdNewProgramsRecord")
	}
}

func (h BugCrowdStore) GetNewRecords() interface{} {
	h.RLock()
	defer h.RUnlock()
	return h.newRecords
}

func (h *BugCrowdStore) Clear() {
	h.Lock()
	defer h.Unlock()
	h.newRecords = nil
}
