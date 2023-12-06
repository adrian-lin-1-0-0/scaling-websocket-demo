package usecases

import (
	"log"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
)

var _ domain.UserUsecase = (*userUsecase)(nil)

type userUsecase struct {
	self        domain.Endpoint
	userRepo    domain.UserRepo
	messagePeer domain.MessagePeer
	pubsub      domain.PubSub
	lookupPeer  domain.LookupPeer
}

func NewUserUsecase(
	self domain.Endpoint,
	userRepo domain.UserRepo,
	messagePeer domain.MessagePeer,
	pubsub domain.PubSub,
	lookupPeer domain.LookupPeer,
) *userUsecase {

	userUsecase := &userUsecase{
		self:        self,
		userRepo:    userRepo,
		messagePeer: messagePeer,
		pubsub:      pubsub,
		lookupPeer:  lookupPeer,
	}

	userUsecase.loop()

	return userUsecase
}

func (s *userUsecase) loop() {
	go func() {
		for {
			msg := <-s.pubsub.Subscribe("message")
			s.Send(msg)
		}
	}()
}

func (s *userUsecase) Send(msg domain.Message) error {

	realremote, err := s.lookupPeer.Lookup(msg.To)
	if err != nil {
		log.Default().Println("lookup error", err)
		return err
	}

	log.Default().Println("lookupPeer : ", msg.To, realremote)

	if realremote == s.self {
		s.pubsub.Publish(string(msg.To), msg)
		return nil
	}

	log.Default().Println("send data to remote : ", msg.To, realremote)

	err = s.messagePeer.SendTo(realremote, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *userUsecase) Register(userId domain.UserId) (<-chan domain.Message, error) {

	err := s.lookupPeer.Register(userId, s.self)
	if err != nil {
		return nil, err
	}

	return s.pubsub.Subscribe(string(userId)), nil
}
