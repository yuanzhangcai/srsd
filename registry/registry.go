package registry

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/yuanzhangcai/srsd/service"
)

// Registry 服务注册主件
type Registry struct {
	opts *Options
	srv  *service.Service
	m    sync.Mutex
	cli  *clientv3.Client
	key  string
}

// NewRegistry 创建服务注册组件
func NewRegistry(srv *service.Service, opts ...Option) *Registry {
	opt := newOptions(opts...)
	return &Registry{
		opts: opt,
		srv:  srv,
		key:  opt.CreateServiceKey(srv),
	}
}

// Register 服务注册
func (c *Registry) Register() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.cli == nil {
		cli, err := c.createEtcdClient()
		if err != nil {
			return err
		}

		c.cli = cli
	}

	val, err := json.Marshal(c.srv)
	if err != nil {
		return err
	}

	lease := clientv3.NewLease(c.cli)
	ctx, cancel := context.WithTimeout(context.Background(), c.opts.Timeout)
	defer cancel()
	grant, err := lease.Grant(ctx, int64(c.opts.TTL/time.Second))
	if err != nil {
		return err
	}

	pCtx, pCancel := context.WithTimeout(context.Background(), c.opts.Timeout)
	defer pCancel()
	_, err = c.cli.Put(pCtx, c.key, string(val), clientv3.WithLease(grant.ID))
	if err != nil {
		return err
	}

	kCtx, kCancel := context.WithTimeout(context.Background(), c.opts.Timeout)
	defer kCancel()
	ch, err := c.cli.KeepAlive(kCtx, grant.ID)
	if err != nil {
		return err
	}

	go func() {
		for range ch {
		}

		for {
			<-time.After(c.opts.Timeout)
			err := c.Register()
			if err == nil {
				return
			}
		}
	}()

	return nil
}

// Deregister 服务删除
func (c *Registry) Deregister() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.cli != nil {
		ctx, cancel := context.WithTimeout(context.Background(), c.opts.Timeout)
		defer cancel()
		_, err := c.cli.Delete(ctx, c.key, clientv3.WithIgnoreLease())
		if err != nil {
			return err
		}

		err = c.cli.Close()
		if err != nil {
			return err
		}
		c.cli = nil
	}
	return nil
}

func (c *Registry) createEtcdClient() (*clientv3.Client, error) {
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