package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"sync"
	"github.com/gorilla/websocket"
)

type PBQueue struct {
	Topics map[string][]chan []byte
	m 		 sync.RWMutex
}


// Publish a message to every subscriber
func (p *PBQueue) pub (topic string, msg []byte) (error) {

}

// Subsribe to a specific topic in the pubsub queue
func (p *PBQueue) sub (topic string) (<-chan []byte, error) {

}

// Unsubscribe from a topic in the pubsub queue
func (p *PBQueue) unsub (topic string, c <-chan []byte) (error) {

}



func main () {
	http.HandleFunc("/", setup)
	http.HandleFunc("/ws", webSocketHandler)
	http.ListenAndServe(":3000", nil)
	fmt.Println("Server is running at port 3000")
}