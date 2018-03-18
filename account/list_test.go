package account

import (
	"testing"

	"github.com/stvp/assert"
)

func TestList(t *testing.T) {
	accounts := List()
	assert.Equal(t, len(accounts), 2)
}
