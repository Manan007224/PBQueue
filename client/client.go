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

type client struct {
	// a client could have more than one subscriptions
	subscribers	map[<-chan []byte]*subscriber
	mtx	sync.RWMutex
}


// this struct represents a subscriber that has a subscription
// Because a client can have more than one subscriptions

type subscriber struct {
	incomingMsg	chan<- []byte // incomingMsg channel can only take incoming messages
	exit 		chan bool
	topic       string
}

func publish(topic string, payload []byte) error {
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

func subscribe(s *subscriber) error {
	log.Printf("subscribing to the topic")

	uri := fmt.Sprintf("ws://%s/sub?topic=%s", server, topic)
	conn, _, err := websocket.DefaultDialer.Dial(uri, make(http.Header))

	if err != nil {
		fmt.Println("error in connecting to the default server localhost:3000")
		return nil, err
	}

	go func() {
		for {
			t, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("couldn't read message")
				conn.Close()
				return
			}
			switch t {
			case <-s.exit:
				conn.Close()
				close(s.exit)
				return
			default:
				s.incomingMsg <- p
			}
		}
	}()
	return nil
}

func (c *client) Unsubscribe(ch <-chan []byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	sub, ok := c.subscribers[ch]

	if ok {
		sub.exit <- true
	}
	return nil
}

func (c *client) Subscribe(topic string) (<-chan []byte, error) {
	ch := make(chan []byte)
	s := &subscriber{
		incomingMsg: ch,
		exit: make(chan bool)
	}
	c.subscribers[ch] = s

	err := subscribe(s)
	if err != nil {
		return ch, err
	}
	return ch, nil
}

func (c *client) Publish(topic string, payload []byte) error {
	return publish(topic, payload)
}