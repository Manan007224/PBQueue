package main

import (
	"log"
	"fmt"
	"flag"
)

var fTopic = flag.Bool("server", false, "Run server")
var fClient1 = flag.Bool("client1", false, "Run client #1")
var fClient2 = flag.Bool("client2", false, "Run client #2")

func main() {
	flag.Parse()

	if *fTopic {
		topic()
	}

	if *fClient1 {
		client("1")
	}

	if *fClient2 {
		client("2")
	}
}

func topic() {
	t := createTopic("example")

	err := t.Publish([]byte("welcome to the Pbq"))

	if err != nil {
		fmt.Println("error in publishing")
	}
}

func client(id string) {
	go func() {
		c := &client {
			route: "/",
		}
		msgChannel, err := c.Subscribe("example")

		for {
			select {
			case m := <-msgChannel:
				log.Println("------ client" + id + string(m) + "------")
			}
		}
	}()
}