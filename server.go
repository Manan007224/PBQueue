package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

var (
	WSUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}
	}
	PBQ = NewPBQ() // default broker
)




func main () {
	http.HandleFunc("/pub", pub)
	http.HandleFunc("/sub", sub)
	http.ListenAndServe(":3000", nil)
	fmt.Println("Server is running at port 3000")
}