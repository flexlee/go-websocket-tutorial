package pubsubredis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flexlee/websocket-dashboard/src/ws"

	"github.com/garyburd/redigo/redis"
)

func SubscribeRedis(hub *ws.Hub) {

	redisClient, err := redis.Dial("tcp", ":32772")
	if err != nil {
		fmt.Println(err)
	}
	defer redisClient.Close()
	psc := redis.PubSubConn{Conn: redisClient}
	psc.Subscribe("chan1")

	for {
		switch n := psc.Receive().(type) {
		case redis.Message:
			// fmt.Printf("Message: %s %s\n", n.Channel, n.Data)

			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(n.Data), "", "\t")
			if err != nil {
				fmt.Println("Error processing JSON: ", err)
			}

			// Need to slowdown sending msg to channel, would overwhelm websocket
			time.Sleep(1 * time.Millisecond)
			fmt.Println(string(prettyJSON.Bytes()))
			hub.Broadcast <- prettyJSON.Bytes()
		case redis.PMessage:
			fmt.Printf("PMessage: %s %s\n", n.Pattern, n.Channel, n.Data)
		case redis.Subscription:
			fmt.Printf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
			if n.Count == 0 {
			}
		case error:
			fmt.Printf("error: %v\n", n)
		}
	}
}
