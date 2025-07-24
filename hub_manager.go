package main

import (
	"log"
	"sync"
)

type HubManager struct {
	hubs          map[string]*Hub
	hubMutex      sync.RWMutex
	unregisterHum chan string
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs:          make(map[string]*Hub),
		unregisterHum: make(chan string),
	}
}

func (hm *HubManager) deleteHub() {
	id := <-hm.unregisterHum

	hm.hubMutex.Lock()
	defer hm.hubMutex.Unlock()

	if _, ok := hm.hubs[id]; ok == true {
		delete(hm.hubs, id)
		log.Printf("room %v deleted", id)
	}
}

func (hm *HubManager) getOrCreateHub(roomID string) *Hub {
	hm.hubMutex.Lock()
	defer hm.hubMutex.Unlock()

	if _, ok := hm.hubs[roomID]; ok == false {
		hub := NewHub(hm, roomID)
		hm.hubs[roomID] = hub
		go hub.run()
		log.Printf("Создан новый хаб для комнаты %s", roomID)
	}

	return hm.hubs[roomID]
}
