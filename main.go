package main

import (
	"net/http"

	"github.com/flexlee/go-websocket-tutorial/ws"
)

func main() {
	ws.InitWS()
	// go hub.start()
	http.HandleFunc("/ws", ws.WsPage)
	http.HandleFunc("/", ws.HomePage)
	http.ListenAndServe(":8080", nil)
}
