package bus

import (
	"testing"
	"time"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
)

func TestMessage(t *testing.T) {

	mb := NewMessageBus(1)

	userID := "userA"
	msgToUserA := "hello UserA"
	msgToUserB := "hello UserB"
	ch := mb.Subscribe(userID)
	messgeCh := mb.Subscribe("message")

	go func() {
		mb.Publish("message", domain.Message{
			UserId: "userB",
			Msg:    msgToUserA,
			To:     domain.UserId(userID),
		})

		mb.Publish("message", domain.Message{
			UserId: "userB",
			Msg:    msgToUserA,
			To:     domain.UserId(userID),
		})

		mb.Publish("message", domain.Message{
			UserId: "userA",
			Msg:    msgToUserB,
			To:     "userB",
		})
	}()

	go func() {

		//A in self
		for msg := range ch {
			t.Log(msg)
			if msg.Msg != msgToUserA {
				t.Error("msg is not hello")
			}
		}

		//B in remote
		for msg := range messgeCh {
			t.Log(msg)
			if msg.Msg != msgToUserB {
				t.Error("msg is not hello")
			}
		}
	}()

	time.Sleep(1 * time.Second)
}
