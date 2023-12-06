package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain/lru"
)

type lookupService struct {
	lookupUsecase domain.LookupUsecase
	cache         *lru.Cache
}

func NewLookupService(lookupUsecase domain.LookupUsecase, cache *lru.Cache) *lookupService {
	return &lookupService{
		lookupUsecase: lookupUsecase,
		cache:         cache,
	}
}

func (s lookupService) LookupHandler(w http.ResponseWriter, r *http.Request) {

	var payload domain.LookupUserPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Default().Println("LookupHandler : ", payload.UserId)

	endpoint, ok := s.cache.Get(string(payload.UserId))

	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	var response domain.LookupUserPayload

	response.UserId = payload.UserId
	response.Endpoint = string(endpoint.(domain.Endpoint))

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s lookupService) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var payload domain.LookupUserPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Default().Println("RegisterHandler : ", payload.UserId, payload.Endpoint)

	err = s.lookupUsecase.Register(domain.UserId(payload.UserId), domain.Endpoint(payload.Endpoint))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
