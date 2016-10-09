package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader
var hub Hub

func InitWS() *Hub {

	upgrader = websocket.Upgrader{}
	hub = Hub{
		Broadcast:    make(chan []byte),
		addClient:    make(chan *Client),
		removeClient: make(chan *Client),
		clients:      make(map[*Client]bool),
	}
	return &hub
}

func StartWS() {
	go hub.start()
}

type Hub struct {
	clients      map[*Client]bool
	Broadcast    chan []byte
	addClient    chan *Client
	removeClient chan *Client
}

func (hub *Hub) start() {
	for {
		select {
		case client := <-hub.addClient:
			hub.clients[client] = true
		case client := <-hub.removeClient:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
		case message := <-hub.Broadcast:
			for client := range hub.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}

type Client struct {
	ws *websocket.Conn
	// Hub passes Broadcast messages to this channel
	send chan []byte
}

// Hub Broadcasts a new message and this fires
func (c *Client) write() {
	// make to the close the connection in case the loop exits
	defer func() {
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte("Error writing message, closing WebSocket..."))
				return
			}

			c.ws.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// New message received so pass it to the Hub
func (c *Client) read() {
	defer func() {
		hub.removeClient <- c
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			hub.removeClient <- c
			c.ws.Close()
			break
		}

		hub.Broadcast <- message
	}
}

func WsPage(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		http.NotFound(w, req)
		return
	}

	client := &Client{
		ws:   conn,
		send: make(chan []byte),
	}

	hub.addClient <- client

	go client.write()
	go client.read()
}

// homePage serves the chat box
func HomePage(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "index.html")
}
