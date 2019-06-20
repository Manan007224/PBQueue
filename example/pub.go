package main

import (
	"log"
	"time"
	"../client"
	"fmt"
)


var (
	c = client.New()
)

func main() {

	tick := time.NewTicker(time.Second)

	go publish()

	for _ = range tick.C {
		if err := c.Publish("foo", []byte(`bar`)); err != nil {
			log.Println(err)
			break
		}
	}
}

func publish() {
	ch, err := c.Subscribe("foo")
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Unsubscribe(ch)

	for i := 0; i < 10; i++ {
		select {
		case e := <-ch:
			log.Println(string(e))
			fmt.Println(string(e))
		}
	}
}