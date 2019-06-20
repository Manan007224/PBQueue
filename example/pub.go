package main

import (
	"log"
	"time"
	"../client"
)

func main() {

	c := client.New()
	tick := time.NewTicker(time.Second)

	for _ = range tick.C {
		if err := c.Publish("foo", []byte(`bar`)); err != nil {
			log.Println(err)
			break
		}
	}
}