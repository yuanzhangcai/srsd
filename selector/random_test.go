package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuanzhangcai/srsd/service"
)

func TestNewRandom(t *testing.T) {
	sel := NewRandom()
	assert.NotNil(t, sel)
}

func TestRandomFilter(t *testing.T) {
	sel := NewRandom()
	assert.NotNil(t, sel)
	var srvs []*service.Service
	for i := 0; i < 100; i++ {
		srvs = append(srvs, service.NewService())
	}

	tmp := sel.Filter("", srvs)
	assert.Equal(t, 1, len(tmp))

	tmp2 := sel.Filter("", srvs)
	assert.Equal(t, 1, len(tmp2))

	assert.NotEqual(t, tmp, tmp2)
}
