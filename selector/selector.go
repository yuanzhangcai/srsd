package selector

import (
	"github.com/yuanzhangcai/srsd/service"
)

// Selector 服务选择器
type Selector interface {
	// Filter 选择过滤器
	Filter(name string, srvs []*service.Service) []*service.Service
}
