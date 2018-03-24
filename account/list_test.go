package account

import (
	"testing"

	"github.com/stvp/assert"
)

func TestList(t *testing.T) {
	acc, _ := Generate("xxx")
	accounts := List()
	assert.NotEqual(t, len(accounts), 0)
	Destroy(acc)
}
