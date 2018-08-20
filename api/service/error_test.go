package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestMissingExecutionError(t *testing.T) {
	e := MissingExecutionError{ID: "test"}
	assert.Equal(t, "Execution test doesn't exists", e.Error())
}
