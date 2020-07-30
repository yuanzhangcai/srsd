package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuanzhangcai/srsd/service"
)

func TestNewRound(t *testing.T) {
	sel := NewRound()
	assert.NotNil(t, sel)
	assert.NotNil(t, sel.index)
}

func TestRoundFilter(t *testing.T) {
	sel := NewRound()
	assert.NotNil(t, sel)
	var srvs []*service.Service
	for i := 0; i < 2; i++ {
		srvs = append(srvs, service.NewService())
	}

	tmp := sel.Filter("zacyuan.com", srvs)
	assert.Equal(t, 1, len(tmp))

	tmp2 := sel.Filter("zacyuan.com", srvs)
	assert.Equal(t, 1, len(tmp2))

	assert.NotEqual(t, tmp, tmp2)
}
