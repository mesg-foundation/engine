package account

import (
	"testing"

	"github.com/stvp/assert"
)

func TestImport(t *testing.T) {
	account, _ := Import("/...", "test")
	assert.Equal(t, account.Name, "test")
	assert.NotEqual(t, account.Address, "")
}

func TestImportWithDefaultName(t *testing.T) {
	account, _ := Import("/...", "")
	assert.NotEqual(t, account.Name, "")
	assert.NotEqual(t, account.Address, "")
}

func TestExport(t *testing.T) {
	account := &Account{
		Name:     "xxx",
		Password: "xxx",
	}
	account.Generate()
	path, _ := account.Export()
	assert.NotEqual(t, path, "")
}
