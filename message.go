package main

type Message struct {
	topic 	string		`json:"topic"`
	payload interface{}	`json:"payload"`
}