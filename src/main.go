package main

import (
	"fmt"
	"net/http"
	"os"

	// "github.com/flexlee/websocket-dashboard/src/listener"
	"github.com/flexlee/websocket-dashboard/src/pubsubredis"
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
	redisAddr := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	redisPort := os.Getenv("REDIS_PORT_6379_TCP_PORT")
	redisURL := fmt.Sprintf("%s:%s", redisAddr, redisPort)
	fmt.Println(redisURL)

	connInfo := fmt.Sprintf("host=%s port=%s dbname=portfolio user=dbUser password=something_super_secret_change_in_production sslmode=disable", dbAddr, dbPort)

	// listenCh := "portfolio_update"
	// listener.InitDBListener(hub, connInfo, listenCh)

	go pubsubredis.SubscribeRedis(hub, redisURL)
	updatedb.UpdateDB(connInfo, redisURL)

	http.HandleFunc("/ws", ws.WsPage)
	http.HandleFunc("/", ws.HomePage)
	http.ListenAndServe(":8080", nil)
}
