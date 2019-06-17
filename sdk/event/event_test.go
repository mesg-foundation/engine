package eventsdk

import (
	"testing"

	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stretchr/testify/require"
)

func TestSubTopic(t *testing.T) {
	serviceHash := "1"
	require.Equal(t, subTopic(serviceHash), hash.Calculate([]string{serviceHash, topic}))
}
