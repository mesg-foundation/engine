package eventsdk

import (
	"testing"

	"github.com/mesg-foundation/core/hash"
	"github.com/stretchr/testify/require"
)

func TestSubTopic(t *testing.T) {
	require.Equal(t, subTopic(hash.Hash{0}), "1.Event")
}
