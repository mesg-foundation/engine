package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMissingExecutionError(t *testing.T) {
	e := MissingExecutionError{ID: "test"}
	require.Equal(t, `Execution "test" doesn't exists`, e.Error())
}
