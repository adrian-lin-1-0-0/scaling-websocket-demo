package bus

import (
	"sync"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
)

var _ domain.PubSub = (*MessageBus)(nil)

type MessageBus struct {
	mu       sync.RWMutex
	subs     map[string]map[chan domain.Message]struct{}
	chanSize int
}

func NewMessageBus(chanSize int) *MessageBus {
	return &MessageBus{
		subs:     make(map[string]map[chan domain.Message]struct{}),
		chanSize: chanSize,
	}
}

func (mb *MessageBus) Subscribe(topic string) <-chan domain.Message {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	ch := make(chan domain.Message, mb.chanSize)
	mb.subs[topic] = map[chan domain.Message]struct{}{ch: {}}

	return ch
}

func (mb *MessageBus) Publish(topic string, msg domain.Message) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	if topic == "message" {

		if len(mb.subs[string(msg.To)]) != 0 {
			for ch := range mb.subs[string(msg.To)] {
				ch <- msg
			}
			return
		}
	}

	subs := mb.subs[topic]
	for ch := range subs {
		ch <- msg
	}
}

func (mb *MessageBus) Unsubscribe(topic string, ch chan domain.Message) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	delete(mb.subs[topic], ch)
}
