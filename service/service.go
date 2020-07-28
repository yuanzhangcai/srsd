package service

import (
	"sync"

	"github.com/google/uuid"
)

// Service 服务注册信息
type Service struct {
	ID       string            `json:"id"`       // 服务唯一ID
	Name     string            `json:"name"`     // 服务名称
	Version  string            `json:"version"`  // 版本
	Host     string            `json:"host"`     // 服务地址
	PProf    string            `json:"pprof"`    // pprof地址
	Metrics  string            `json:"metrics"`  // prometheus指标曝露地址
	Metadata map[string]string `json:"metadata"` // 扩展信息
}

// NewService 创建Service对象
func NewService() *Service {
	return &Service{
		ID:       uuid.New().String(),
		Version:  "latest",
		Metadata: make(map[string]string),
	}
}

// Services 服务列表
type Services struct {
	sync.RWMutex
	Index uint
	List  []*Service
}

// NewServices 创建Services对象
func NewServices() *Services {
	return &Services{}
}

// Delete 删除指定id的服务注册信息
func (c *Services) Delete(id string) {
	if len(c.List) == 0 {
		return
	}

	c.Lock()
	defer c.Unlock()

	for i, one := range c.List {
		if one.ID == id {
			tmp := append([]*Service{}, c.List[0:i]...)
			tmp = append(tmp, c.List[i+1:]...)
			c.List = tmp
		}
	}
}

// Put 将服务注册信息加入服务列表
func (c *Services) Put(srv *Service) {
	c.Lock()
	defer c.Unlock()

	has := false
	for i, one := range c.List {
		if one.ID == srv.ID {
			c.List[i] = srv
			has = true
			break
		}
	}

	if !has {
		c.List = append(c.List, srv)
	}
}
