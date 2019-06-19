package main 

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"time"
	"fmt"
)

type client struct {
	route		  string
	conn		  *websocket.Conn
	dialer		  *websocket.Dialer
	connId		  string
}

// connect to the webscoket server
func (c* client) connect (route string) error {
	uri := fmt.Sprintf("ws://localhost:3000/%s", route)
	c.dialer = websocket.DefaultDialer
	c.dialer.HandshakeTimeout = 1 * time.Second
	log.Println("Dialing %s", uri)
	conn, _, err := c.dialer.Dial(uri, http.Header{})
	c.conn = conn
	return err
}

func (c *client) Subscribe (topic string) (<-chan []byte, error) {
	log.Printf("subscribing to the topic")
	chanMessage := make(chan []byte, 100)

	if er := c.connect(topic); er != nil {
		log.Println("connection failed")
		// TODO
		return nil, er
	}

	go func() {
		for {
			t, p, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("couldn't read message")
				c.conn.Close()
				return
			}
			switch t {
			case websocket.CloseMessage:
				log.Println("Close message, clossing channel")
				c.conn.Close()
				close(chanMessage)
				return
			default:
				chanMessage <- p
			}
		}
	}()
	return chanMessage, nil
}