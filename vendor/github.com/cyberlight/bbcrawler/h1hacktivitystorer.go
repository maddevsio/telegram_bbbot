package bbcrawler

import (
	"fmt"
	"strconv"
	"sync"
)

type H1HacktivityStore struct {
	PathToDb   string
	newRecords []H1HactivityRecord
	all map[string]H1HactivityRecord
	sync.RWMutex
}

func (h *H1HacktivityStore) IsEmpty() (bool, error) {
	fmt.Println("Store Hacktivity running")
	h.RLock()
	defer h.RUnlock()
	return len(h.all) == 0, nil
}

func (h *H1HacktivityStore) Store(data interface{}) error {
	fmt.Println("Store Hacktivity running")
	h.Lock()
	defer h.Unlock()
	if response, ok := data.(H1HactivityResponse); ok {
		for _, v := range response.Reports {
			fmt.Print("h.")
			if _, ok := h.all[strconv.Itoa(v.Id)]; !ok {
				fmt.Print("h+")
				h.newRecords = append(h.newRecords, v)
				h.all[strconv.Itoa(v.Id)] = v
			}
		}
	} else if rec, ok := data.(H1HactivityRecord); ok {
		if _, ok := h.all[strconv.Itoa(rec.Id)]; !ok {
			fmt.Print("h+")
			h.all[strconv.Itoa(rec.Id)] = rec
		}
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
