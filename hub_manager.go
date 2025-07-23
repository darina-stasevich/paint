package main

import (
	"log"
	"sync"
)

type HubManager struct {
	hubs     map[string]*Hub
	hubMutex sync.RWMutex
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func (hm *HubManager) getOrCreateHub(roomID string) *Hub {
	hm.hubMutex.Lock()
	defer hm.hubMutex.Unlock()

	if _, ok := hm.hubs[roomID]; ok == false {
		hub := NewHub()
		hm.hubs[roomID] = hub
		go hub.run()
		log.Printf("Создан новый хаб для комнаты %s", roomID)
	}

	return hm.hubs[roomID]
}
