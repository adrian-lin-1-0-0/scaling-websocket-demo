package domain

import (
	"sync"
)

type PubSub interface {
	Subscribe(topic string) <-chan Message
	Publish(topic string, msg Message)
	Unsubscribe(topic string, ch chan Message)
}

type pubSub struct {
	mu   sync.RWMutex
	subs map[string][]chan Message
}

func NewpubSub() *pubSub {
	return &pubSub{
		subs: make(map[string][]chan Message),
	}
}

func (ps *pubSub) Subscribe(topic string) <-chan Message {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan Message, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)

	return ch
}

func (ps *pubSub) Publish(topic string, msg Message) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subs := ps.subs[topic]
	for _, ch := range subs {
		ch <- msg
	}
}
