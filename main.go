package main

import (
	"net/http"
	"log"
	"sync"
	"github.com/gorilla/websocket"
	"fmt"
	"io/ioutil"
	"flag"
	"os"
	"bufio"
	pbqclient "./client"
)

var (
	// server flags
	address = flag.String("address", ":3000", "server address")

	// client flags
	client = flag.Bool("client", false, "Run Pbq client")
	publish = flag.Bool("publish", false, "Publish via Pbq client")
	subscribe = flag.Bool("subscribe", false, "Subscribe via Pbq client")
	topic = flag.String("topic", "", "Topic for which client subscribes or publishes")
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

func (p *Pbq) pub(topic string, payload []byte) error {
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

func (p *Pbq) unsub(topic string, sub <-chan []byte) error {
	p.mtx.RLock()
	subscribers, ok := p.subQueue[topic]
	p.mtx.RUnlock()

	if !ok {
		return nil
	}

	var subs []chan []byte
	for _, subscriber := range subscribers {
		if subscriber == sub {
			continue
		}
		subs = append(subs, subscriber)
	}

	p.mtx.Lock()
	p.subQueue[topic] = subs
	p.mtx.Unlock()
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
	defer pbq.unsub(topic, ch)
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

func checkFlags() {
	flag.Parse()	

	if *client && !*publish && !*subscribe {
		log.Fatal("Specify whether to publish or subscribe")
	}

	if *client && len(*topic) == 0 {
		log.Fatal("topic not specified")
	}
}	

func handleClient() {
	pbqClient:= pbqclient.New()
	if *subscribe {
		channel, err := pbqClient.Subscribe(*topic)
		if err != nil {
			fmt.Println(err)
			return
		}
		for msg := range channel {
			fmt.Println(string(msg))
		}
	} else { // process publish
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			pbqClient.Publish(*topic, scanner.Bytes())
		}
	}	
}

func main () {

	checkFlags()

	if *client {
		handleClient()
		return
	} 
	http.HandleFunc("/pub", pub)
	http.HandleFunc("/sub", sub)
	fmt.Println("pbq listening on port 3000")
	http.ListenAndServe(":3000", nil)
}