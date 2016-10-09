package main

import "net/http"

func main() {
	go hub.start()
	http.HandleFunc("/ws", wsPage)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}
