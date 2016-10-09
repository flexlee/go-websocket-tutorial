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

			// Need to slowdown sending msg to channel, would overwhelm websocket
			time.Sleep(1 * time.Millisecond)
			// fmt.Println(string(prettyJSON.Bytes()))
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

// InitDBListener initializes DB connection and pass *ws.Hub for sedning notifications
func InitDBListener(hub *ws.Hub, connInfo string, listenCh string) {

	_, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(connInfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen(listenCh)
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
