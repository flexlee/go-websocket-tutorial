package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/flexlee/websocket-dashboard/src/listener"
	"github.com/flexlee/websocket-dashboard/src/updatedb"
	"github.com/flexlee/websocket-dashboard/src/ws"
)

func main() {
	hub := ws.InitWS()
	ws.StartWS()

	dbAddr := os.Getenv("POSTGRES_PORT_5432_TCP_ADDR")
	fmt.Println(dbAddr)
	dbPort := os.Getenv("POSTGRES_PORT_5432_TCP_PORT")
	fmt.Println(dbPort)

	connInfo := fmt.Sprintf("host=%s port=%s dbname=portfolio user=dbUser password=something_super_secret_change_in_production sslmode=disable", dbAddr, dbPort)

	listenCh := "portfolio_update"
	listener.InitDBListener(hub, connInfo, listenCh)
	go updatedb.UpdateDB(connInfo)

	http.HandleFunc("/ws", ws.WsPage)
	http.HandleFunc("/", ws.HomePage)
	http.ListenAndServe(":8080", nil)
}
