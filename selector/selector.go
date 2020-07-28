package selector

import (
	"github.com/yuanzhangcai/srsd/service"
)

// Selector 服务选择器
type Selector interface {
	// Filter 选择过滤器
	Filter(srvs []*service.Service) []*service.Service
	// Get 获取当个服务器
	Get(srvs []*service.Service) *service.Service
}
