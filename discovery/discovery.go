package discovery

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/yuanzhangcai/srsd/selector"
	"github.com/yuanzhangcai/srsd/service"
)

// Event 监听事件
type Event = clientv3.Event

// Discovery 服务发现组件
type Discovery struct {
	opts    *Options
	cli     *clientv3.Client
	m       sync.RWMutex
	cancel  map[string]context.CancelFunc
	srvList map[string][]*service.Service
}

// NewDiscovery 创建服务发现组件
func NewDiscovery(opts ...Option) *Discovery {
	opt := newOptions(opts...)

	return &Discovery{
		opts:    opt,
		srvList: make(map[string][]*service.Service),
		cancel:  make(map[string]context.CancelFunc),
	}
}

// Start 开启服务发现
func (c *Discovery) Start(key string) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.cli == nil {
		cli, err := c.createEtcdClient()
		if err != nil {
			return err
		}

		c.cli = cli
	}

	if _, ok := c.cancel[key]; ok {
		return nil
	}

	err := c.loadAll(key)
	if err != nil {
		return err
	}

	err = c.startWatch(key)
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
		id = id[index+1:]
	}
	return id
}

func (c *Discovery) loadAll(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.opts.Timeout)
	defer cancel()
	key = c.opts.Prefix + key
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
		c.putSrv(key, srv)
	}
	return nil
}

func (c *Discovery) putSrv(key string, srv *service.Service) {
	list, ok := c.srvList[key]
	if !ok {
		list = []*service.Service{}
	}

	has := false
	for i, one := range list {
		if one.ID == srv.ID {
			list[i] = srv
			has = true
			break
		}
	}

	if !has {
		list = append(list, srv)
	}

	c.srvList[key] = list
}

func (c *Discovery) delSrv(key, id string) {
	list, ok := c.srvList[key]
	if !ok {
		return
	}

	if len(list) == 0 {
		return
	}

	for i, one := range list {
		if one.ID == id {
			list = append(list[0:i], list[i+1:]...)
		}
	}
	c.srvList[key] = list
}

func (c *Discovery) startWatch(key string) error {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel[key] = cancel
	watchKey := c.opts.Prefix + key
	ch := c.cli.Watch(ctx, watchKey, clientv3.WithPrefix())
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

		switch one.Type {
		case mvccpb.DELETE:
			c.delSrv(name, id)
		case mvccpb.PUT:
			srv := &service.Service{}
			err := json.Unmarshal(one.Kv.Value, srv)
			if err != nil {
				continue
			}
			c.putSrv(name, srv)
		}

		if c.opts.Watch != nil {
			c.opts.Watch(one)
		}
	}

	return nil
}

// Select 获取服务信息
func (c *Discovery) Select(name string, selectors ...selector.Selector) *service.Service {
	c.m.RLock()
	defer c.m.RUnlock()
	var list []*service.Service
	ok := false
	if name != "" {
		list, ok = c.srvList[name]
		if !ok {
			return nil
		}
	} else {
		for _, one := range c.srvList {
			list = append(list, one...)
		}
	}

	if len(list) == 0 {
		return nil
	}

	if len(selectors) == 0 {
		selectors = c.opts.Selectors
	}

	for _, one := range selectors {
		list = one.Filter(name, list)
	}

	if len(list) > 0 {
		return list[0]
	}

	return nil
}

// GetAll 获取所有服务器信
func (c *Discovery) GetAll(name string) []*service.Service {
	c.m.RLock()
	defer c.m.RUnlock()

	var list []*service.Service
	ok := false
	if name != "" {
		list, ok = c.srvList[name]
		if !ok {
			return nil
		}
	} else {
		for _, one := range c.srvList {
			list = append(list, one...)
		}
	}

	return list
}

// Stop 停止服务发现
func (c *Discovery) Stop() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.cli != nil {
		c.cli.Close()
		c.cli = nil
		c.srvList = make(map[string][]*service.Service)
	}

	return nil
}
