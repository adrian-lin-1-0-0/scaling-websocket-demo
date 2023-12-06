package transport

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
	"golang.org/x/net/websocket"
)

type chatService struct {
	userUsecase domain.UserUsecase
}

func NewChatService(userUsecase domain.UserUsecase) *chatService {
	return &chatService{
		userUsecase: userUsecase,
	}
}

func (s chatService) Handler(ws *websocket.Conn) {
	userId := ws.Request().URL.Query().Get("userId")
	log.Default().Println("New connection userId : ", userId)

	receiver, err := s.userUsecase.Register(domain.UserId(userId))
	if err != nil {
		fmt.Println("Can't register", err)
		panic(err)
	}

	go func() {
		for {
			msg := <-receiver
			fmt.Println("Reciver message", msg)

			bytes, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Can't marshal", err)
				break
			}

			if err = websocket.Message.Send(ws, bytes); err != nil {
				fmt.Println("Can't send :", err)
				break
			}
		}
	}()

	for {
		var request domain.Message

		data := make([]byte, 1024)
		n, err := ws.Read(data)
		if err != nil {
			log.Println("Can't receive", err)
			break
		}

		err = json.Unmarshal(data[:n], &request)
		if err != nil {
			log.Println("Can't unmarshal", err)
			break
		}

		request.UserId = domain.UserId(userId)
		s.userUsecase.Send(request)
	}
}
