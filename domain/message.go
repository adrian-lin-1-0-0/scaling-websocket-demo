package domain

type Message struct {
	UserId UserId `json:"userId"`
	Msg    string `json:"msg"`
	To     UserId `json:"to"`
}

type MessageUsecase interface {
	Send(Message) error
}

type MessagePeer interface {
	SendTo(Endpoint, Message) error
}
