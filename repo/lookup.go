package repo

import (
	"fmt"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain/consistenthash"
)

const (
	defaultReplicas = 50
)

var _ domain.Lookup = (*lookupRepo)(nil)

type lookupRepo struct {
	basePath  string
	self      domain.Endpoint
	hashs     *consistenthash.Map
	discovery domain.Discovery
}

func NewLookupRepo(basePath string, self domain.Endpoint, discovery domain.Discovery) *lookupRepo {

	lookupRepo := &lookupRepo{
		self:      self,
		basePath:  basePath,
		hashs:     consistenthash.New(defaultReplicas, nil),
		discovery: discovery,
	}

	lookupRepo.discovery.Register(domain.Endpoint(
		basePath + string(self),
	))

	if err := lookupRepo.setNodes(); err != nil {
		panic(err)
	}

	lookupRepo.autoDiscovery()

	return lookupRepo
}

func (r *lookupRepo) AddService(endpoint domain.Endpoint) error {
	r.hashs.Add(string(endpoint))
	return nil
}

func (r *lookupRepo) RemoveService(endpoint domain.Endpoint) error {
	r.hashs.Remove(string(endpoint))
	return nil
}

func (r *lookupRepo) GetService(key string) (domain.Endpoint, error) {
	return domain.Endpoint(r.hashs.Get(key)), nil
}

func (r *lookupRepo) setNodes() error {
	if r.discovery == nil {
		panic("discovery is nil")
	}

	endpoints, err := r.discovery.GetServices(r.basePath)
	if err != nil {
		return fmt.Errorf("get services failed: %v", err)
	}

	services := make([]string, len(endpoints))
	for i, endpoint := range endpoints {
		services[i] = string(endpoint[len(r.basePath):])
	}

	r.hashs.Add(services...)

	return nil

}

func (r *lookupRepo) autoDiscovery() {
	if r.discovery == nil {
		panic("discovery is nil")
	}

	eventChan, err := r.discovery.WatchPrefix(r.basePath)
	if err != nil {
		panic(err)
	}

	go func() {
		for event := range eventChan {
			if event.Type == domain.EventPut {
				r.AddService(domain.Endpoint(event.Key[len(r.basePath):]))
			}
			if event.Type == domain.EventDelete {
				r.RemoveService(domain.Endpoint(event.Key[len(r.basePath):]))
			}
		}
	}()

}
