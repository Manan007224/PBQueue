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
	WSUpgrader = websocket.Upgrader {
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	PBQ = NewPBQ() // default broker
)

func sub(w http.ResponseWriter, r *http.Request) {
	conn, err := WSUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("websocket connection failed")
		http.Error(w, "Could not open websocket	connection", http.StatusBadRequest)
		return
	}
	// create a new subscriber
	subscriber, err := PBQ.Attach()

	if err != nil {
		log.Println("error in attaching a new broker")
		http.Error(w, "Could not retrieve events", http.StatusInternalServerError)	
		return
	}

	topic := r.URL.Query().Get("topic")
	ch, err := PBQ.sub(topic)

	if err != nil {
		log.Println("couldn't retrieve the topic")
		http.Error(w, "Could not retrieve events", http.StatusInternalServerError)
		return	
	}

	defer PBQ.unsub(topic, ch)

	for {
		select {
		case msg := <-ch:
			err = conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				log.Println("error in writing the message")
				return
			}
		}
	}
}

func pub(w http.ResponseWriter, r *http.Request) {
	var msg Message
	b, err := ioUtil.ReadAll(r.Body)
	err = json.Unmarshal(b, &msg)
	if err != nil {
		panic(err)
	}
	if err != nil {
		http.Error(w, "Pub Error", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	err = PBQ.Pub(msg)

	if err != nil {
		http.Error(w, "Pub Error", http.StatusInternalServerError)
		return
	}
}


func main () {
	http.HandleFunc("/pub", pub)
	http.HandleFunc("/sub", sub)
	http.ListenAndServe(":3000", nil)
	fmt.Println("Server is running at port 3000")
}