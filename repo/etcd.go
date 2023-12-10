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
	client  *etcd.Client
	self    string
	ctx     context.Context
	ttl     time.Duration
	leaseID etcd.LeaseID
	once    sync.Once
	leaser  etcd.Lease
}

const (
	defaultTTL = 10 * time.Second
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
		client: cli,
		self:   self,
		ctx:    ctx,
		ttl:    defaultTTL,
		once:   sync.Once{},
	}, nil
}

func (d *Discovery) GetServices(basePath string) ([]domain.Endpoint, error) {
	resp, err := d.client.Get(d.ctx, basePath, etcd.WithPrefix())
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

func (d *Discovery) WatchPrefix(basePath string) (domain.EventChan, error) {
	watchChan := d.client.Watch(d.ctx, basePath, etcd.WithPrefix())
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

func (d *Discovery) InitLease(ttl int64) error {
	lease := etcd.NewLease(d.client)
	d.leaser = lease
	leaseResp, err := lease.Grant(d.ctx, ttl)
	if err != nil {
		return fmt.Errorf("grant lease failed: %v", err)
	}
	d.leaseID = leaseResp.ID
	return nil
}

func (d *Discovery) Register(endpoint domain.Endpoint) error {
	d.once.Do(func() {
		err := d.InitLease(int64(d.ttl.Seconds()))
		if err != nil {
			return
		}
		go d.keepAlive()
	})

	if d.leaseID == 0 {
		panic(errors.New("leaseID is nil"))
	}

	_, err := d.client.Put(d.ctx, string(endpoint), "", etcd.WithLease(d.leaseID))

	if err != nil {
		panic("register service failed :" + err.Error())
	}

	log.Default().Println("register service success")

	return err
}

func (d *Discovery) Unregister(endpoint domain.Endpoint) error {
	_, err := d.client.Delete(d.ctx, string(endpoint))
	return err
}

func (c *Discovery) keepAlive() {
	leaseKeepAliveQueue, err := c.leaser.KeepAlive(c.ctx, c.leaseID)
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
