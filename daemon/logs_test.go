package daemon

import (
	"testing"

	"github.com/stvp/assert"
)

func TestLogs(t *testing.T) {
	startForTest()
	reader, err := Logs()
	assert.Nil(t, err)
	assert.NotNil(t, reader)
}
