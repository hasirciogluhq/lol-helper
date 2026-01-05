package gui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type ItemManager struct {
	nameToID map[string]int
	mutex    sync.RWMutex
}

func NewItemManager() *ItemManager {
	im := &ItemManager{
		nameToID: make(map[string]int),
	}
	go im.fetchItemData()
	return im
}

func (im *ItemManager) fetchItemData() {
	// Fetch item.json from DDragon
	// Using a recent version, ideally this should be dynamic but hardcoded for stability now
	resp, err := http.Get("https://ddragon.leagueoflegends.com/cdn/14.1.1/data/en_US/item.json")
	if err != nil {
		fmt.Printf("Failed to fetch item data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data map[string]struct {
			Name string `json:"name"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Failed to decode item data: %v\n", err)
		return
	}

	im.mutex.Lock()
	defer im.mutex.Unlock()

	for idStr, item := range result.Data {
		var id int
		fmt.Sscanf(idStr, "%d", &id)
		im.nameToID[strings.ToLower(item.Name)] = id
	}
}

func (im *ItemManager) GetItemID(name string) int {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	// Try exact match first
	if id, ok := im.nameToID[strings.ToLower(name)]; ok {
		return id
	}

	// Try partial match if exact fails (simple fuzzy)
	lowerName := strings.ToLower(name)
	for itemName, id := range im.nameToID {
		if strings.Contains(itemName, lowerName) || strings.Contains(lowerName, itemName) {
			return id
		}
	}

	return 0
}
