package discovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yuanzhangcai/srsd/registry"
	"github.com/yuanzhangcai/srsd/selector"
	"github.com/yuanzhangcai/srsd/service"
)

var testEtcdAddr = []string{"127.0.0.1:2379"}

func TestNewDiscovery(t *testing.T) {
	dis := NewDiscovery(
		Addresses(testEtcdAddr),
		Username("zacyuan"),
		Password("12345678"),
		Prefix("/zacyuan/test"),
		Timeout(3*time.Second),
		Selectors(selector.NewRandom()),
	)
	assert.NotNil(t, dis)
	assert.NotNil(t, dis.srvList)
	assert.Equal(t, testEtcdAddr, dis.opts.Addresses)
	assert.Equal(t, "zacyuan", dis.opts.Username)
	assert.Equal(t, "12345678", dis.opts.Password)
	assert.Equal(t, "/zacyuan/test/", dis.opts.Prefix)
	assert.Equal(t, 3*time.Second, dis.opts.Timeout)
	assert.Less(t, 0, len(dis.opts.Selectors))
}

func TestStart(t *testing.T) {
	info := service.NewService()
	info.Name = "zacyuan.com"
	info.Host = "127.0.0.1:4001"
	reg1 := registry.NewRegistry(info, registry.Addresses(testEtcdAddr), registry.TTL(10*time.Second))
	_ = reg1.Start()

	info = service.NewService()
	info.Name = "zacyuan.com"
	info.Host = "127.0.0.1:4002"
	reg2 := registry.NewRegistry(info, registry.Addresses(testEtcdAddr), registry.TTL(10*time.Second))
	_ = reg2.Start()

	t.Run("start etcd error", func(t *testing.T) {
		dis := NewDiscovery(
			Addresses([]string{}),
		)

		err := dis.Start("")
		assert.NotNil(t, err)
	})

	dis := NewDiscovery(Addresses(testEtcdAddr))

	t.Run("start success", func(t *testing.T) {
		err := dis.Start("")
		assert.Nil(t, err)
	})

	t.Run("start success again", func(t *testing.T) {
		err := dis.Start("")
		assert.Nil(t, err)
	})

	t.Run("start modify svr info", func(t *testing.T) {
		info.Name = "zacyuan.com"
		info.Host = "127.0.0.1:4003"
		reg := registry.NewRegistry(info, registry.Addresses(testEtcdAddr), registry.TTL(10*time.Second))
		_ = reg.Start()
		time.Sleep(1 * time.Second)
		_ = reg.Stop()
		time.Sleep(1 * time.Second)
	})

	t.Run("Select success", func(t *testing.T) {
		srv := dis.Select("zacyuan.com")
		assert.NotNil(t, srv)
	})

	t.Run("Select success", func(t *testing.T) {
		srv := dis.Select("")
		assert.NotNil(t, srv)
	})

	t.Run("Select no data", func(t *testing.T) {
		srv := dis.Select("zacyuan.com.xyz")
		assert.Nil(t, srv)
	})

	t.Run("GetAll success", func(t *testing.T) {
		srvs := dis.GetAll("zacyuan.com")
		assert.Less(t, 0, len(srvs))
	})

	t.Run("GetAll success no name", func(t *testing.T) {
		srvs := dis.GetAll("")
		assert.Less(t, 0, len(srvs))
	})

	t.Run("Stop success", func(t *testing.T) {
		err := dis.Stop()
		assert.Nil(t, err)
		assert.Nil(t, dis.cli)
	})

	t.Run("Stop again", func(t *testing.T) {
		err := dis.Stop()
		assert.Nil(t, err)
	})

}

func TestGetServiceName(t *testing.T) {
	key := "/srsd/services/zacyuan.com/aaaa"
	dis := NewDiscovery(Addresses(testEtcdAddr))
	name := dis.getServiceName(key)
	assert.Equal(t, "zacyuan.com", name)

	key = "/srsd/services/srsd/services/zacyuan.com/aaaa"
	name = dis.getServiceName(key)
	assert.Equal(t, "srsd/services/zacyuan.com", name)

	key = "/srsd/services/zacyuan.com"
	name = dis.getServiceName(key)
	assert.Equal(t, "zacyuan.com", name)
}

func TestGetServiceID(t *testing.T) {
	key := "/srsd/services/zacyuan.com/aaaa"
	dis := NewDiscovery(Addresses(testEtcdAddr))
	id := dis.getServiceID(key)
	assert.Equal(t, "aaaa", id)
}
