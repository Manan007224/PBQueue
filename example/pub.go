package main

import (
	"fmt"
)

func main() {
	t := createTopic("example")
	err := t.Publish([]byte("welcome to the Pbq"))

	if err != nil {
		fmt.Println("error in publishing")
	}
}
