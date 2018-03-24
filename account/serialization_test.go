package account

import (
	"testing"

	"github.com/stvp/assert"
)

func TestImport(t *testing.T) {
	account, _ := Import("/...", "test")
	assert.NotEqual(t, account.Address, "")
}

func TestImportWithDefaultName(t *testing.T) {
	account, _ := Import("/...", "")
	assert.NotEqual(t, account.Address, "")
}
