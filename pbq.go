package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader {
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Topic struct {
	// the clients that have subscribed to this topic
	clientQueue map[connId] chan []byte

	topic 	string
	mtx		sync.RWMutex
}

func createTopic (topic string) error {
	t := &Topic {
		clientQueue: make(map[string]chan []byte),
		topic: topic,
		mtx: sync.RWMutex{},
	}
	http.HandleFunc("/topic/"+topic, topicHandler)
}

func topicHandler (t *Topic) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.upgrade(w, r, nil)
		if err != nil {
			log.Println("websocket connection failed")
			http.Error(w, "Couldn't open ws connection",http.StatusBadRequest)
			return
		}

		t.mtx.Lock()
		ch, ok := t.clientQueue[connID(uuid.NewV4.string())]
		if !ok {
			ch = make(chan []byte)
			clientQueue[c.id] = ch
		}
		t.mtx.UnLock()
			
		for {
			select {
			case msg := <-ch:
				err = conn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil {
					log.Printf("error in sending the event")
				}
			}
		}	
	}	
}

// send the message to all the subscribers
func (t *Topic) Publish (msg []byte) error {
	go func() {
		for _, q := range t.clientQueue{
			select {
				case: q.ch <- payload
			}
		}()
	}
} 

func main () {
}