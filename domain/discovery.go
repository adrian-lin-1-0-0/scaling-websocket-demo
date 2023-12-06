package domain

type Discovery interface {
	GetServices(string) ([]Endpoint, error)
	WatchPrefix(string) (EventChan, error)
	Register(Endpoint) error
	Unregister(Endpoint) error
}

type EventChan <-chan Event

type Endpoint string

type Service struct {
	Name     string
	Endpoint string
}

type Event struct {
	Type  EventType
	Key   string
	Value string
}

type EventType int

const (
	EventPut EventType = iota
	EventDelete
)
