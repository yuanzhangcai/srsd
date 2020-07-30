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
