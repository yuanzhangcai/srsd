package service

import (
	"github.com/google/uuid"
	"github.com/yuanzhangcai/srsd/utils"
)

// Service 服务注册信息
type Service struct {
	ID         string            `json:"id"`          // 服务唯一ID
	Name       string            `json:"name"`        // 服务名称
	Version    string            `json:"version"`     // 版本
	Host       string            `json:"host"`        // 服务地址
	PProf      string            `json:"pprof"`       // pprof地址
	Metrics    string            `json:"metrics"`     // prometheus指标曝露地址
	Metadata   map[string]string `json:"metadata"`    // 扩展信息
	CreateTime string            `json:"create_time"` // 服务注册时间
}

// NewService 创建Service对象
func NewService() *Service {
	return &Service{
		ID:       uuid.New().String(),
		Version:  "latest",
		Metadata: make(map[string]string),
	}
}

// GetRealIP 获取Host、Metrics、PProf的真实IP
func (c *Service) GetRealIP() error {
	var err error
	if c.Host != "" {
		c.Host, err = utils.GetRealAddr(c.Host)
		if err != nil {
			return err
		}
	}

	if c.Metrics != "" {
		c.Metrics, err = utils.GetRealAddr(c.Metrics)
		if err != nil {
			return err
		}
	}

	if c.PProf != "" {
		c.PProf, err = utils.GetRealAddr(c.PProf)
		if err != nil {
			return err
		}
	}

	return nil
}
