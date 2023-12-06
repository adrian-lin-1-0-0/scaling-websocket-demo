package repo

import (
	"log"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain/lru"
)

var _ domain.UserRepo = (*userRepo)(nil)

type userRepo struct {
	lookup domain.Lookup
	cache  *lru.Cache
}

func NewUserRepo(lookup domain.Lookup, cache *lru.Cache) *userRepo {
	return &userRepo{
		lookup: lookup,
		cache:  cache,
	}
}

func (r *userRepo) Lookup(userId domain.UserId) (domain.Endpoint, error) {
	endpoint, err := r.lookup.GetService(string(userId))
	if err != nil {
		log.Default().Println("Lookup error", err)
		return "", err
	}

	log.Default().Println("Lookup : ", userId, endpoint)

	return endpoint, nil
}

func (r *userRepo) Register(userId domain.UserId, endpoint domain.Endpoint) error {
	log.Default().Println("Register : ", userId, endpoint)
	r.cache.Add(string(userId), endpoint)
	return nil
}
