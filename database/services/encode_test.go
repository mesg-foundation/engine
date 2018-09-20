package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeError(t *testing.T) {
	hash := "IDToTest"
	_, err := decode(hash, []byte("oaiwdhhiodoihwaiohwa"))
	require.Error(t, err)
	require.IsType(t, &DecodeError{}, err)
}
