package pubsub

import (
	"testing"

	"github.com/stvp/assert"
)

func TestSubscribe(t *testing.T) {
	key := "TestSubscribe"
	res := Subscribe(key)
	assert.NotNil(t, res)
	assert.Equal(t, len(listeners[key]), 1)
}

func TestSubscribeMultipleTimes(t *testing.T) {
	key := "TestSubscribeMultipleTimes"
	Subscribe(key)
	Subscribe(key)
	assert.Equal(t, len(listeners[key]), 2)
}
