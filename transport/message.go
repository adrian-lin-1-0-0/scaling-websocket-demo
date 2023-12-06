package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
)

type MessageService struct {
	messageUsecase domain.MessageUsecase
}

func NewMessageService(messageUsecase domain.MessageUsecase) *MessageService {
	return &MessageService{
		messageUsecase: messageUsecase,
	}
}

func (s MessageService) Handler(w http.ResponseWriter, r *http.Request) {

	var msg domain.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Default().Println("MessageService: ", msg)

	err = s.messageUsecase.Send(msg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
