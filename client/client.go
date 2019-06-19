package main 

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"fmt"
	"bytes"
)

var (
	server = "localhost:3000"	
)

func Publish (topic string, payload []byte) error {
	resp, err := http.Post(fmt.Sprintf("http://%s/pub?topic=%s", server, topic), "application/json", 
		bytes.NewBuffer(payload))

	if err != nil {
		return err
	}

	resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Non 200 response %d", resp.StatusCode)
	}
	return nil
}

func Subscribe (topic string) (<-chan []byte, error) {
	log.Printf("subscribing to the topic")

	uri := fmt.Sprintf("ws://%s/sub?topic=%s", server, topic)
	conn, _, err := websocket.DefaultDialer.Dial(uri, make(http.Header))

	if err != nil {
		fmt.Println("error in connecting to the default server localhost:3000")
		return nil, err
	}

	chanMessage := make(chan []byte, 100)

	go func() {
		for {
			t, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("couldn't read message")
				conn.Close()
				return
			}
			switch t {
			case websocket.CloseMessage:
				log.Println("Close message, clossing channel")
				conn.Close()
				close(chanMessage)
				return
			default:
				chanMessage <- p
			}
		}
	}()
	return chanMessage, nil
}