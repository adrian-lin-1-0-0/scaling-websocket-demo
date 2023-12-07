package domain

type PubSub interface {
	Subscribe(topic string) <-chan Message
	Publish(topic string, msg Message)
	Unsubscribe(topic string, ch chan Message)
}
