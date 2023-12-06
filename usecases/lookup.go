package usecases

import "github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"

var _ domain.LookupUsecase = (*lookupUsecase)(nil)

type lookupUsecase struct {
	userRepo domain.UserRepo
}

func NewLookupUsecase(userRepo domain.UserRepo) *lookupUsecase {
	return &lookupUsecase{
		userRepo: userRepo,
	}
}

func (s *lookupUsecase) Lookup(userId domain.UserId) (domain.Endpoint, error) {
	return s.userRepo.Lookup(userId)
}

func (s *lookupUsecase) Register(userId domain.UserId, endpoint domain.Endpoint) error {
	return s.userRepo.Register(userId, endpoint)
}
