package domain

type UserId string

type UserUsecase interface {
	Send(Message) error
	Register(UserId) (<-chan Message, error)
}

type UserRepo interface {
	Lookup(UserId) (Endpoint, error)
	Register(UserId, Endpoint) error
}
