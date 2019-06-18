package main

import (
	"net/http"
	// "encoding/json"
	// "fmt"
	// "io/ioutil"
	"log"
	"sync"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
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
	clientQueue map[string] chan []byte

	topic 	string
	mtx		sync.RWMutex
}

func createTopic (topic string) {
	t := &Topic {
		clientQueue: make(map[string]chan []byte),
		topic: topic,
		mtx: sync.RWMutex{},
	}
	handler := topicHandler(t)
	http.HandleFunc("/topic/"+topic, handler)
}

func topicHandler (t *Topic) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("websocket connection failed")
			http.Error(w, "Couldn't open ws connection",http.StatusBadRequest)
			return
		}

		t.mtx.Lock()
		id := uuid.Must(uuid.NewV4()).String()
		ch, ok := t.clientQueue[id]
		if !ok {
			ch = make(chan []byte)
			t.clientQueue[id] = ch
		}
		t.mtx.Unlock()
			
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
func (t *Topic) Publish (payload []byte) error {
	go func() {
		for _, ch := range t.clientQueue{
			select {
				case ch <- payload:
				default:
			}
		}
	}()
	return nil
} 

func main () {
	http.ListenAndServe(":3000", nil)
}