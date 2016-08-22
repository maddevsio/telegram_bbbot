package bbcrawler

import (
	"sync"
	"fmt"
)

type BugCrowdStore struct {
	PathToDb   string
	newRecords []BugCrowdNewProgramsRecord
	all map[string]BugCrowdNewProgramsRecord
	sync.RWMutex
}

func (h *BugCrowdStore) IsEmpty() (bool, error) {
	h.RLock()
	defer h.RUnlock()
	return len(h.all) == 0, nil
}

func (h *BugCrowdStore) Store(data interface{}) error {
	fmt.Println("Store BugCrowd running")
	h.Lock()
	defer h.Unlock()

	if recs, ok := data.([]BugCrowdNewProgramsRecord); ok {
		for _, rec := range recs {
			fmt.Print(".")
			if _, ok := h.all[rec.Name]; !ok {
				fmt.Print(".+")
				h.all[rec.Name] = rec
				h.newRecords = append(h.newRecords, rec)
			}
		}
		return nil
	} else if rec, ok := data.(BugCrowdNewProgramsRecord); ok {
		if _, ok := h.all[rec.Name]; !ok {
			fmt.Print(".+")
			h.all[rec.Name] = rec
			h.newRecords = append(h.newRecords, rec)
		}
	} else {
		return fmt.Errorf("Can't convert to type %s", "BugCrowdNewProgramsRecord")
	}
	return nil
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
