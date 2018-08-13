package dockertest

import (
	"testing"

	"github.com/stvp/assert"
)

// TestNotFoundErr makes sure NotFoundErr to implement docker client's notFound interface.
func TestNotFoundErr(t *testing.T) {
	err := NotFoundErr{}
	assert.True(t, err.NotFound())
	assert.Equal(t, "not found", err.Error())
}
