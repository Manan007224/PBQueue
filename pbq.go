package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"sync"
	"github.com/gorilla/websocket"
)

type Pbq struct {
	topics			map[string] []Subscriber
	mtx					sync.RWMutex
}


func NewPBQ () *Pbq {
	return &PBQ {
		topics : map[string]subscribers{},
		subscribers : subscribers{},
		mtx : sync.RWMutex{},
	}
}

func (p *Pbq) Attach() (*Subscriber, error) {
	// create a random id and attach it to the subscriber later
	rid := make([]byte, 50)
	err := rand.Read(rid)
	if err != nil {
		return nil, err
	}
	rid = hex.EncodeToString(rid)
	s := &Subscriber {
		messages : make(chan Message),
		mtx : sync.RWMutex{},
		alive : true,
		id : rid,
		createdAt : time.Now().UnixNano(),
	}
	return s, nil
} 

func (p *Pbq) Pub(topic string, msg Message) error {
	p.mtx.Lock()
	subscribers, ok := p.topics[msg.topic]
	p.mtx.UnLock()

	if !ok || len(subscribers) < 1 {
		return nil
	}

	go func () {
		for _, sub := range subscribers {
			subchannel, ok := sub[msg.topic]
			if !ok {
				continue
			} 
			select {
				case subchannel <- msg:
				default:
			}
		}
	}()
	return nil
}

func (p *Pbq) Sub(msg Message, s *Subscriber) (<-chan Message, error) {
	ch := make(chan Message, 100)
	p.mtx.Lock()
	p.topics[msg.topic] = append(p.topics[msg.topic], s.id)
	s.subscriptions[msg.topic] = ch
	p.mtx.UnLock()
	return ch
}

func (p *Pbq) Unsub(msg Message) error {
	p.mtx.Lock()
	subscribers, ok := p.topics[msg.topics]
	p.mtx.UnLock()
}