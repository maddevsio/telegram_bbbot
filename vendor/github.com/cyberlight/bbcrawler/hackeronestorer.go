package bbcrawler

import (
	"fmt"
	"sync"
)

type HackerOneStore struct {
	PathToDb   string
	newRecords []HackerOneRecord
	all map[string]HackerOneRecord
	sync.RWMutex
}

func (h *HackerOneStore) IsEmpty() (bool, error) {
	fmt.Println("Store running")
	h.RLock()
	defer h.RUnlock()
	return len(h.all) == 0, nil
}

func (h *HackerOneStore) Store(data interface{}) error {
	fmt.Println("Store running")
	h.Lock()
	defer h.Unlock()

	if response, ok := data.(HackerOneResponse); ok {
		for _, v := range response.Results {
			fmt.Print(".")
			if _, ok := h.all[v.Handle]; !ok {
				fmt.Print("+")
				h.all[v.Handle] = v
				h.newRecords = append(h.newRecords, v)
			}
		}
	} else if rec, ok := data.(HackerOneRecord); ok {
		if _, ok := h.all[rec.Handle]; !ok {
			fmt.Print("+")
			h.all[rec.Handle] = rec
		}
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
