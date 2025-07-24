package main

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handlePaintRequest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	hm := NewHubManager()
	go hm.deleteHub()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Генерируем случайный номер комнаты
		roomID := strconv.Itoa(rand.Intn(10000))
		http.Redirect(w, r, "/paint/"+roomID, http.StatusFound)
	})

	mux.HandleFunc("/paint/{roomID}", handlePaintRequest)

	mux.HandleFunc("/ws/{roomID}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hm, w, r)
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
