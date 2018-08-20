package core

import (
	"testing"

	"github.com/stvp/assert"
)

func TestNotRunningServiceError(t *testing.T) {
	e := NotRunningServiceError{ServiceID: "test"}
	assert.Equal(t, "Service test is not running", e.Error())
}
