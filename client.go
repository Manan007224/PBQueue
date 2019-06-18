package main 

import (
	"sync"
)

type client struct {
	route		  string
	conn		  *websokcet.Conn
	dialer		  *websokcet.Dialer
	connId		  string
}

// connect to the webscoket server
func (c* client) connect (route string) error {
	uri := fmt.Sprintf("ws://localhost:3000/%s", route)
	c.dialer = webscoket.Default.Dialer
	c.dialer.HandshakeTimeout = 1 * time.Second
	log.Println("Dialing %s", uri)
	conn, err := c.Dialer.Dial(uri, http.Header{})
	c.conn = conn
	return c, err
}

func (c *client) Subscribe (topic string) (<-chan []byte, error) {
	log.Printf("subscribing to the topic")
	chanMessage := make(chan []byte, 100)

	if er := c.connect(topic); err != nil {
		log.Println("connection failed")
		// TODO
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