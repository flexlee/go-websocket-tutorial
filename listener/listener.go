package listener

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flexlee/go-websocket-tutorial/ws"
	"github.com/lib/pq"
)

func waitForNotification(hub *ws.Hub, l *pq.Listener) {
	for {
		select {
		case notice := <-l.Notify:
			fmt.Println("Received data from channel [", notice.Channel, "] :")
			// Prepare notification payload for pretty print
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(notice.Extra), "", "\t")
			if err != nil {
				fmt.Println("Error processing JSON: ", err)
				return
			}
			fmt.Println(string(prettyJSON.Bytes()))
			hub.Broadcast <- prettyJSON.Bytes()
			return
		case <-time.After(600 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}

func InitDBListener(hub *ws.Hub) {
	var conninfo string = "host=0.0.0.0 port=32770 dbname=todopy_pg user=keli password=something_super_secret_change_in_production sslmode=disable"

	_, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("todos_updates")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")
	go func() {
		for {
			waitForNotification(hub, listener)
		}
	}()
}
