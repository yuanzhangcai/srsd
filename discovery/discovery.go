package discovery

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/yuanzhangcai/srsd/service"
)

// Event 监听事件
type Event = clientv3.Event

// Discovery 服务发现组件
type Discovery struct {
	opts    *Options
	cli     *clientv3.Client
	m       sync.RWMutex
	cancel  context.CancelFunc
	srvList map[string]*service.Services
}

// NewDiscovery 创建服务发现组件
func NewDiscovery(opts ...Option) *Discovery {
	opt := newOptions(opts...)

	return &Discovery{
		opts:    opt,
		srvList: make(map[string]*service.Services),
	}
}

// Start 开启服务发现
func (c *Discovery) Start(name string) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.cli == nil {
		cli, err := c.createEtcdClient()
		if err != nil {
			return err
		}

		c.cli = cli
	}

	err := c.loadAll(name)
	if err != nil {
		return err
	}

	return nil
}

func (c *Discovery) createEtcdClient() (*clientv3.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opts.Timeout)
	defer cancel()

	return clientv3.New(clientv3.Config{
		Context:     ctx,
		Endpoints:   c.opts.Addresses,
		DialTimeout: c.opts.Timeout,
		Username:    c.opts.Username,
		Password:    c.opts.Password,
	})
}

func (c *Discovery) getServiceName(key string) string {
	key = strings.Replace(key, c.opts.Prefix, "", 1)
	index := strings.LastIndex(key, "/")
	if index > 0 {
		key = key[:index]
	}
	return key
}

func (c *Discovery) getServiceID(key string) string {
	id := key
	index := strings.LastIndex(id, "/")
	if index > 0 {
		id = id[index:]
	}
	return id
}

func (c *Discovery) loadAll(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opts.Timeout)
	defer cancel()
	key := c.opts.Prefix + name
	resp, err := c.cli.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		srv := &service.Service{}
		err := json.Unmarshal(kv.Value, srv)
		if err != nil {
			continue
		}

		key := c.getServiceName(string(kv.Key))
		one, ok := c.srvList[key]
		if !ok {
			one = service.NewServices()
			c.srvList[key] = one
		}
		one.Put(srv)
	}
	return nil
}

func (c *Discovery) startWatch(name string) error {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	key := c.opts.Prefix + name
	ch := c.cli.Watch(ctx, key, clientv3.WithPrefix())
	go func() {
		for resp := range ch {
			if resp.Canceled {
				return
			}
			_ = c.reload(&resp)
		}
	}()

	return nil
}

func (c *Discovery) reload(resp *clientv3.WatchResponse) error {
	if resp == nil {
		return nil
	}

	c.m.Lock()
	defer c.m.Unlock()

	for _, one := range resp.Events {
		key := string(one.Kv.Key)
		name := c.getServiceName(key)
		id := c.getServiceID(key)

		srvs, ok := c.srvList[name]
		if !ok {
			srvs = service.NewServices()
			c.srvList[name] = srvs
		}

		switch one.Type {
		case mvccpb.DELETE:
			srvs.Delete(id)
		case mvccpb.PUT:
			srv := &service.Service{}
			err := json.Unmarshal(one.Kv.Value, srv)
			if err != nil {
				continue
			}
			srvs.Put(srv)
		}

		if c.opts.Watch != nil {
			c.opts.Watch(one)
		}
	}

	return nil
}

// Get 获取服务信息
func (c *Discovery) Get(service string) *service.Service {
	c.m.RLock()
	defer c.m.RUnlock()
	srvs, ok := c.srvList[service]
	if !ok {
		return nil
	}
	_ = srvs

	return nil
}
