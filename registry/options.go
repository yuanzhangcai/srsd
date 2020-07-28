package registry

import (
	"time"

	"github.com/yuanzhangcai/srsd/service"
)

var (
	defaultPrefix    = "/srsd/services/"
	defaultAddresses = []string{"127.0.0.1:2379"}
	defaultTimeout   = 5 * time.Second
	defaultTTL       = 10 * time.Second
)

// Option 设置服务注册参数
type Option func(*Options)

// Watch watch回调函数
type Watch func(event int32, host string)

// Options 服务注册参数
type Options struct {
	Addresses []string      // etcd地址
	Username  string        // etcd用户名
	Password  string        // etcd密码
	Prefix    string        //服务注册前缀
	Timeout   time.Duration // etcd超时时间
	TTL       time.Duration // 服务存活时间
}

// NewOptions 那建服务注册参数对象
func newOptions(opts ...Option) *Options {
	opt := &Options{
		Addresses: defaultAddresses,
		Prefix:    defaultPrefix,
		Timeout:   defaultTimeout,
		TTL:       defaultTTL,
	}

	for _, one := range opts {
		one(opt)
	}
	return opt
}

//

// CreateServiceKey 生成服务注册key
func (c *Options) CreateServiceKey(srv *service.Service) string {
	return c.Prefix + srv.Name + "/" + srv.ID
}

// Addresses 设置etcd地址
func Addresses(addresses []string) Option {
	return func(opt *Options) {
		opt.Addresses = addresses
	}
}

// Username 设置etcd用户名
func Username(userName string) Option {
	return func(opt *Options) {
		opt.Username = userName
	}
}

// Password 设置etcd密码
func Password(password string) Option {
	return func(opt *Options) {
		opt.Password = password
	}
}

// Prefix 设置服务发现前缀
func Prefix(prefix string) Option {
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	return func(opt *Options) {
		opt.Prefix = prefix
	}
}

// Timeout 设置etcd超时时间
func Timeout(timeout time.Duration) Option {
	return func(opt *Options) {
		opt.Timeout = timeout
	}
}

// TTL 设置服务存活时间
func TTL(ttl time.Duration) Option {
	return func(opt *Options) {
		opt.TTL = ttl
	}
}
