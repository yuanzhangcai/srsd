package selector

import (
	"sync"

	"github.com/yuanzhangcai/srsd/service"
)

// Round 轮询选择器
type Round struct {
	m     sync.Mutex
	index map[string]uint
}

// NewRound 创建论询选择器
func NewRound() *Round {
	return &Round{
		index: make(map[string]uint),
	}
}

// Filter 轮询过滤器
func (c *Round) Filter(name string, srvs []*service.Service) []*service.Service {
	c.m.Lock()
	defer c.m.Unlock()
	index := c.index[name] % uint(len(srvs))
	c.index[name]++
	return []*service.Service{srvs[index]}
}
