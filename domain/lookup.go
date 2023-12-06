package domain

type Lookup interface {
	GetService(key string) (Endpoint, error)
	AddService(Endpoint) error
	RemoveService(Endpoint) error
}

type LookupUserPayload struct {
	UserId   string `json:"userId"`
	Endpoint string `json:"endpoint"`
}

type LookupUsecase interface {
	Lookup(UserId) (Endpoint, error)
	Register(UserId, Endpoint) error
}

type LookupPeer interface {
	Register(UserId, Endpoint) error
	Lookup(UserId) (Endpoint, error)
}
