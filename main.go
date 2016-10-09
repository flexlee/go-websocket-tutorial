package main

import "net/http"
import "github.com/flexlee/go_websocket_tutorial/websocket"

var hub = websocket.Hub{
	broadcast:    make(chan []byte),
	addClient:    make(chan *websocket.Client),
	removeClient: make(chan *websocket.Client),
	clients:      make(map[*websocket.Client]bool),
}

func main() {
	go hub.start()
	http.HandleFunc("/ws", websocket.WsPage)
	http.HandleFunc("/", websocket.HomePage)
	http.ListenAndServe(":8080", nil)
}
