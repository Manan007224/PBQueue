package main

import (
	"log"
	"fmt"
	"../client"
)


func main() {

	c := client.New()

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