package main

import (
	"fmt"
	"log"
	"../client"
)

func main() {
	ch, err := client.Subscribe("example")
	if err != nil {
		log.Println(err)
		log.Println("wtf is happening")
		return
	}

	for e := range ch {
		log.Println(string(e))
	}

	log.Println("Channel closed")
}