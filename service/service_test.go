package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	srv := NewService()
	assert.NotNil(t, srv)
	assert.NotNil(t, srv.Metadata)
	assert.NotEmpty(t, srv.ID)
}

func TestGetRealIP(t *testing.T) {
	srv := NewService()
	srv.Host = ":4000"
	srv.PProf = "10.10.8.59:4001"
	srv.Metrics = "10.10.8.159:4002"
	err := srv.GetRealIP()
	assert.Nil(t, err)
	assert.NotEqual(t, ":4000", srv.Host)
	assert.Equal(t, "10.10.8.59:4001", srv.PProf)
	assert.Equal(t, "10.10.8.159:4002", srv.Metrics)
}
