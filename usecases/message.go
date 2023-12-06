package usecases

import (
	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
)

type messageUsecase struct {
	pubsub domain.PubSub
}

func NewMessageUsecase(pubsub domain.PubSub) *messageUsecase {
	return &messageUsecase{
		pubsub: pubsub,
	}
}

func (s *messageUsecase) Send(msg domain.Message) error {
	s.pubsub.Publish(string(msg.To), msg)
	return nil
}
