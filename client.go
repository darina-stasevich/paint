package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (client *Client) readPump() {
	defer func() {
		client.hub.unregister <- client
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("error while reading: %v", err)
			break
		}

		client.hub.broadcast <- message
	}
}

func (client *Client) writePump() {
	defer func() {
		client.conn.Close()
	}()

	for {
		message, ok := <-client.send
		if !ok {
			client.conn.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}

		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

func serveWs(hm *HubManager, w http.ResponseWriter, r *http.Request) {

	log.Println("serve ws")
	roomID := r.PathValue("roomID")
	if roomID == "" {
		http.Error(w, "Не указан ID комнаты", http.StatusBadRequest)
		return
	}

	hub := hm.getOrCreateHub(roomID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.readPump()
	go client.writePump()

}
