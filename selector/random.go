package selector

import (
	"math/rand"
	"time"

	"github.com/yuanzhangcai/srsd/service"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random 随机选择器
type Random struct {
}

// NewRandom 创建随机选择器
func NewRandom() *Random {
	return &Random{}
}

// Filter 随机过滤器
func (c *Random) Filter(name string, srvs []*service.Service) []*service.Service {
	index := rand.Int() % len(srvs)
	return []*service.Service{srvs[index]}
}
