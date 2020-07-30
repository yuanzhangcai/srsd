package registry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yuanzhangcai/srsd/service"
)

var testEtcdAddr = []string{"127.0.0.1:2379"}

func TestNewRegistry(t *testing.T) {
	srv := service.NewService()
	srv.Name = "zacyuan.com"
	srv.Host = "127.0.0.1:4444"

	reg := NewRegistry(srv,
		Addresses(testEtcdAddr),
		Username("zacyuan"),
		Password("12345678"),
		Prefix("/zacyuan/test"),
		Timeout(3*time.Second),
		TTL(60*time.Second),
	)
	assert.NotNil(t, reg)
	assert.Equal(t, testEtcdAddr, reg.opts.Addresses)
	assert.Equal(t, "zacyuan", reg.opts.Username)
	assert.Equal(t, "12345678", reg.opts.Password)
	assert.Equal(t, "/zacyuan/test/", reg.opts.Prefix)
	assert.Equal(t, 3*time.Second, reg.opts.Timeout)
	assert.Equal(t, 60*time.Second, reg.opts.TTL)
}

func TestStart(t *testing.T) {
	t.Run("Start no etcd address", func(t *testing.T) {
		srv := service.NewService()
		srv.Name = "zacyuan.com"

		srv.Host = "127.0.0.1:4444"
		reg := NewRegistry(srv, Addresses([]string{}))
		err := reg.Start()
		assert.NotNil(t, err)
	})

	t.Run("Start success", func(t *testing.T) {
		srv := service.NewService()
		srv.Name = "zacyuan.com"

		srv.Host = "127.0.0.1:4444"
		reg := NewRegistry(srv, Addresses(testEtcdAddr))
		err := reg.Start()
		assert.Nil(t, err)
	})

	t.Run("keepAlive", func(t *testing.T) {
		srv := service.NewService()
		srv.Name = "zacyuan.com"

		srv.Host = "127.0.0.1:4444"
		reg := NewRegistry(srv, Addresses(testEtcdAddr), TTL(2*time.Second))
		err := reg.Start()
		assert.Nil(t, err)
		time.Sleep(3 * time.Second)

		reg.m.Lock()
		key := reg.opts.CreateServiceKey(srv)
		value, err := reg.cli.Get(context.Background(), key)
		assert.Nil(t, err)
		assert.Less(t, int64(0), value.Count)
		reg.m.Unlock()
	})

	t.Run("restart", func(t *testing.T) {
		srv := service.NewService()
		srv.Name = "zacyuan.com"

		srv.Host = "127.0.0.1:4444"
		reg := NewRegistry(srv, Addresses(testEtcdAddr), TTL(2*time.Second))
		err := reg.Start()
		assert.Nil(t, err)

		reg.m.Lock()
		reg.cli.Close()
		reg.cli = nil
		reg.m.Unlock()

		time.Sleep(1 * time.Second)

		reg.m.Lock()
		key := reg.opts.CreateServiceKey(srv)
		value, err := reg.cli.Get(context.Background(), key)
		assert.Nil(t, err)
		assert.Less(t, int64(0), value.Count)
		reg.m.Unlock()
	})
}

func TestStop(t *testing.T) {
	srv := service.NewService()
	srv.Name = "zacyuan.com"

	srv.Host = "127.0.0.1:4444"
	reg := NewRegistry(srv, Addresses(testEtcdAddr), TTL(2*time.Second))

	t.Run("Stop no start", func(t *testing.T) {
		err := reg.Stop()
		assert.Nil(t, err)
	})

	t.Run("Stop success", func(t *testing.T) {
		err := reg.Start()
		assert.Nil(t, err)

		err = reg.Stop()
		assert.Nil(t, err)
	})

}
