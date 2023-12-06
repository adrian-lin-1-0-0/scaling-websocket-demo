package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
	etcd "go.etcd.io/etcd/client/v3"
)

var (
	ErrServiceNameEmpty     = errors.New("service name is empty")
	ErrServiceEndpointEmpty = errors.New("service endpoint is empty")
	ErrGetServicesFailed    = errors.New("get services failed")
)

type Discovery struct {
	client    *etcd.Client
	self      string
	ctx       context.Context
	ttl       time.Duration
	heartbeat time.Duration
	leaseID   etcd.LeaseID
	once      sync.Once
}

const (
	defaultTTL       = 10 * time.Second
	defaultHeartbeat = 3 * time.Second
)

var _ domain.Discovery = (*Discovery)(nil)

func NewDiscovery(ctx context.Context, endpoints []string, self string) (*Discovery, error) {

	cli, err := etcd.New(etcd.Config{
		Endpoints:   endpoints,
		Context:     ctx,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Discovery{
		client:    cli,
		self:      self,
		ctx:       ctx,
		ttl:       defaultTTL,
		heartbeat: defaultHeartbeat,
		once:      sync.Once{},
	}, nil
}

func (c *Discovery) GetServices(basePath string) ([]domain.Endpoint, error) {
	resp, err := c.client.Get(c.ctx, basePath, etcd.WithPrefix())
	if err != nil {
		return nil, ErrGetServicesFailed
	}

	var endpoints []domain.Endpoint
	for _, kv := range resp.Kvs {
		// if string(kv.Key) == c.self {
		// 	continue
		// }
		endpoints = append(endpoints, domain.Endpoint(string(kv.Key)))
	}
	return endpoints, nil
}

func (c *Discovery) WatchPrefix(basePath string) (domain.EventChan, error) {
	watchChan := c.client.Watch(c.ctx, basePath, etcd.WithPrefix())
	eventChan := make(chan domain.Event)
	go func() {
		for resp := range watchChan {
			for _, ev := range resp.Events {
				var event domain.Event
				switch ev.Type {
				case etcd.EventTypePut:
					event.Type = domain.EventPut
				case etcd.EventTypeDelete:
					event.Type = domain.EventDelete
				}
				event.Key = string(ev.Kv.Key)
				event.Value = string(ev.Kv.Value)
				eventChan <- event
			}
		}
	}()
	return eventChan, nil
}

func (c *Discovery) InitLease(ttl int64) error {
	lease := etcd.NewLease(c.client)
	leaseResp, err := lease.Grant(c.ctx, ttl)
	if err != nil {
		return fmt.Errorf("grant lease failed: %v", err)
	}
	c.leaseID = leaseResp.ID
	return nil
}

func (c *Discovery) Register(endpoint domain.Endpoint) error {
	c.once.Do(func() {
		err := c.InitLease(int64(c.ttl.Seconds()))
		if err != nil {
			return
		}
		go c.loop()
	})

	if c.leaseID == 0 {
		panic(errors.New("leaseID is nil"))
	}

	_, err := c.client.Put(c.ctx, string(endpoint), "", etcd.WithLease(c.leaseID))

	if err != nil {
		panic("register service failed :" + err.Error())
	}

	log.Default().Println("register service success")

	return err
}

func (c *Discovery) Unregister(endpoint domain.Endpoint) error {
	_, err := c.client.Delete(c.ctx, string(endpoint))
	return err
}

func (c *Discovery) loop() {
	tick := time.NewTicker(c.heartbeat)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			c.keepAlive()
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Discovery) keepAlive() {
	lease := etcd.NewLease(c.client)
	leaseKeepAliveQueue, err := lease.KeepAlive(c.ctx, c.leaseID)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case _, ok := <-leaseKeepAliveQueue:
			if !ok {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}
