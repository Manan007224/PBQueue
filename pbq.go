package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"sync"
	"github.com/gorilla/websocket"
)

type Pbq struct {
	topics	map[string] []*Subscriber
	mtx		sync.RWMutex
}


func NewPBQ () *Pbq {
	return &PBQ {
		topics : map[string]subscribers{},
		subscribers : subscribers{},
		mtx : sync.RWMutex{},
	}
}

func (p *Pbq) Attach() (*Subscriber, error) {
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

func (p *Pbq) Pub(msg Message) error {
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

func (p *Pbq) Sub(topic string, s *Subscriber) (<-chan Message, error) {
	ch := make(chan Message, 100)
	p.mtx.Lock()
	p.topics[topic] = append(p.topics[topic], s.id)
	s.subscriptions[topic] = ch
	p.mtx.UnLock()
	return ch
}

// Method to remove a subscriber from a topic in the pbq
func (p *Pbq) Unsub(msg Message, subc <-chan Message) error {
	p.mtx.Lock()
	subscribers, ok := p.topics[msg.topics]
	messages, okc := s.subscriptions[msg.topic]
	p.mtx.UnLock()

	if !ok {
		return nil
	}

	if !okc {
		delete(s.subscriptions, msg.topic)
	}

	temp := make(chan *Subscriber)
	for _, sub := range subscribers {
		if sub != subc {
			temp = append(temp, sub)
		}
		continue
	}
	p.mtx.Lock()
	p.topics[msg.topic] = temp
	p.mtx.UnLock()
}