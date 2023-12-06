package transport

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
)

type messagePeer struct {
	prefix string
}

var _ domain.MessagePeer = (*messagePeer)(nil)

func NewMessagePeer(prefix string) *messagePeer {
	return &messagePeer{
		prefix: prefix,
	}
}

func (s *messagePeer) SendTo(endpoint domain.Endpoint, msg domain.Message) error {

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = http.Post("http://"+string(endpoint)+s.prefix, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Default().Println("SendTo error", err)
		return err
	}

	return nil
}
