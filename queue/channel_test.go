package queue

import (
	"testing"

	"github.com/stvp/assert"
)

func TestNamespace(t *testing.T) {
	ch := Channel{Kind: Events, Name: "*"}
	assert.Equal(t, ch.namespace(), "event.*")
}
