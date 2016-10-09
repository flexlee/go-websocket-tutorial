package main

import (
	"net/http"

	"github.com/flexlee/go-websocket-tutorial/listener"
	"github.com/flexlee/go-websocket-tutorial/ws"
)

func main() {
	hub := ws.InitWS()
	ws.StartWS()
	listener.InitDBListener(hub)
	// ws.InitDBListener(&hub)

	http.HandleFunc("/ws", ws.WsPage)
	http.HandleFunc("/", ws.HomePage)
	http.ListenAndServe(":8080", nil)
}
