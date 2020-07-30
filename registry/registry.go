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
	opts    *Options
	srv     *service.Service
	m       sync.Mutex
	cli     *clientv3.Client
	key     string
	started bool
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

// Start 开启服务注册
func (c *Registry) Start() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.started {
		return nil
	}

	if c.cli == nil {
		cli, err := c.createEtcdClient()
		if err != nil {
			return err
		}

		c.cli = cli
	}

	c.srv.CreateTime = time.Now().Format("2006-01-02 15:04:05")
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

	return c.keepAlive(grant.ID)
}

func (c *Registry) keepAlive(grantID clientv3.LeaseID) error {
	ctx := context.Background()
	ch, err := c.cli.KeepAlive(ctx, grantID)
	if err != nil {
		return err
	}

	go func() {
		for range ch {
		}

		c.m.Lock()
		started := c.started
		c.m.Unlock()
		if !started { // 服务停止，无需重启
			return
		}

		// KeepAlive异常结束时，重启服务
		for {
			err := c.Stop()
			if err != nil {
				time.Sleep(c.opts.Timeout)
				continue
			}

			err = c.Start()
			if err == nil {
				return
			}
			time.Sleep(c.opts.Timeout)
		}
	}()

	c.started = true

	return nil
}

// Stop 停止服务注册
func (c *Registry) Stop() error {
	c.m.Lock()
	defer c.m.Unlock()

	if !c.started {
		return nil
	}

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
	c.started = false
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
