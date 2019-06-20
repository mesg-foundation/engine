package eventsdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubTopic(t *testing.T) {
	serviceHash := []byte{0}
	require.Equal(t, subTopic(serviceHash), "1.Event")
}
