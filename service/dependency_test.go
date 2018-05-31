package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestExtractPortEmpty(t *testing.T) {
	ports := extractPorts(&Dependency{})
	assert.Equal(t, len(ports), 0)
}

func TestExtractPorts(t *testing.T) {
	ports := extractPorts(&Dependency{
		Ports: []string{
			"80",
			"3000:8080",
		},
	})
	assert.Equal(t, len(ports), 2)
	assert.Equal(t, ports[0].Target, uint32(80))
	assert.Equal(t, ports[0].Published, uint32(80))
	assert.Equal(t, ports[1].Target, uint32(8080))
	assert.Equal(t, ports[1].Published, uint32(3000))
}
