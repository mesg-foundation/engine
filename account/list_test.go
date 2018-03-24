package account

import (
	"testing"

	"github.com/stvp/assert"
)

func TestList(t *testing.T) {
	generate("xxx", "name")
	accounts := List()
	assert.NotEqual(t, len(accounts), 0)
}
