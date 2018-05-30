package daemon

import (
	"testing"

	"github.com/stvp/assert"
)

func TestStop(t *testing.T) {
	startForTest()
	err := Stop()
	assert.Nil(t, err)
}
