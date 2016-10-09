package main

import (
	"net/http"

	"github.com/flexlee/go-websocket-tutorial/listener"
	"github.com/flexlee/go-websocket-tutorial/ws"
)

func main() {
	hub := ws.InitWS()
	ws.StartWS()

	var connInfo string = "host=0.0.0.0 port=32770 dbname=todopy_pg user=keli password=something_super_secret_change_in_production sslmode=disable"
	var listenCh string = "todos_updates"
	listener.InitDBListener(hub, connInfo, listenCh)

	http.HandleFunc("/ws", ws.WsPage)
	http.HandleFunc("/", ws.HomePage)
	http.ListenAndServe(":8080", nil)
}
