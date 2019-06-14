package main 

import (
	"sync"
)

type Subscriber struct {
	subscribtions	map[string]chan Message 
	mtx 		 	sync.RWMutex
	alive    		bool
	creatdAt 		time.Time
	id       		string
}

func (s *Subscriber) GetMessages () <-chan []Message {
	return s.messages
}

func (s *Subscriber) IsAlive() bool {
	return this.alive
}