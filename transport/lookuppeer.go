package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
)

type lookupPeer struct {
	userRepo domain.UserRepo
}

var _ domain.LookupPeer = (*lookupPeer)(nil)

func NewLookupPeer(userRepo domain.UserRepo) *lookupPeer {

	return &lookupPeer{
		userRepo: userRepo,
	}
}

func (s *lookupPeer) Lookup(userId domain.UserId) (domain.Endpoint, error) {

	remote, err := s.userRepo.Lookup(userId)
	if err != nil {
		return "", err
	}

	if remote == "" {
		return "", fmt.Errorf("Not found")
	}

	var payload domain.LookupUserPayload

	payload.UserId = string(userId)

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	response, err := http.Post("http://"+string(remote)+"/lookup", "application/json", bytes.NewReader(data))
	if err != nil {
		log.Default().Println("lookup remote error", err)
		return "", err
	}

	defer response.Body.Close()

	var lookupUserPayload domain.LookupUserPayload

	err = json.NewDecoder(response.Body).Decode(&lookupUserPayload)
	if err != nil {
		log.Default().Println("lookup remote error", err)
		return "", err
	}

	return domain.Endpoint(lookupUserPayload.Endpoint), nil
}

func (s *lookupPeer) Register(userId domain.UserId, endpoint domain.Endpoint) error {

	remote, err := s.userRepo.Lookup(userId)
	if err != nil {
		return err
	}

	if remote == "" {
		return fmt.Errorf("Not found")
	}

	var payload domain.LookupUserPayload

	payload.UserId = string(userId)
	payload.Endpoint = string(endpoint)

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = http.Post("http://"+string(remote)+"/register", "application/json", bytes.NewReader(data))
	if err != nil {
		log.Default().Println("Register error", err)
		return err
	}

	return nil

}
