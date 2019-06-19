package main

import (
	"net/http"
	"log"
	"sync"
	"github.com/gorilla/websocket"
	"fmt"
	"io/ioutil"
)

type Pbq struct {
	// the clients that have subscribed to this topic
	subQueue	map[string] []chan []byte
	mtx			sync.RWMutex
}

var (
	upgrader = websocket.Upgrader {
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	pbq = &Pbq{
		subQueue: make(map[string][]chan []byte),
	}
)

func (p *Pbq) sub(topic string) (<-chan []byte, error) {
	ch := make(chan []byte, 100)
	p.mtx.RLock()
	p.subQueue[topic] = append(p.subQueue[topic], ch)
	p.mtx.RUnlock()
	return ch, nil
}

func (p *Pbq) pub (topic string, payload []byte) error {
	p.mtx.RLock()
	subscribers, ok := p.subQueue[topic]
	p.mtx.RUnlock()

	if !ok {
		return nil
	}

	go func() {
		for _, subscriber := range subscribers {
			select {
			case subscriber <- payload:
			default:
			}
		}
	}()
	return nil
}


func sub(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket connection failed")
		http.Error(w, "Couldn't open ws connection",http.StatusBadRequest)
		return
	}

	topic := r.URL.Query().Get("topic")
	ch, err := pbq.sub(topic)
	if err != nil {
		log.Println("topic doesn't exist")
		http.Error(w, "Could not retrieve events",http.StatusInternalServerError)
	}
	for {
		select {
		case msg := <-ch:
			err = conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				log.Printf("error in sending the event")
				return
			}
		}
	}	
}

func pub(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Pub Error", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	err = pbq.pub(topic, b)
	if err != nil {
		http.Error(w, "Pub Error", http.StatusInternalServerError)
		return
	}
}


func main () {
	http.HandleFunc("/pub", pub)
	http.HandleFunc("/sub", sub)
	fmt.Println("pbq listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
