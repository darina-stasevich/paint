package main

import (
	"encoding/json"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	history    [][]byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		history:    make([][]byte, 0),
	}
}

func (hub *Hub) broadcastUserCount() {

	message := map[string]interface{}{
		"type":  "user_count",
		"count": len(hub.clients),
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	for c, _ := range hub.clients {
		select {
		case c.send <- jsonMessage:
		default:
			close(c.send)
			delete(hub.clients, c)
		}
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
			hub.broadcastUserCount()
			go func() {
				for _, message := range hub.history {
					client.send <- message
				}
			}()
		case client := <-hub.unregister:
			delete(hub.clients, client)
			hub.broadcastUserCount()
		case message := <-hub.broadcast:
			dirtyClients := false
			for client, _ := range hub.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hub.clients, client)
					dirtyClients = true
				}
			}
			hub.history = append(hub.history, message)
			if dirtyClients {
				hub.broadcastUserCount()
			}
		}
	}
}
