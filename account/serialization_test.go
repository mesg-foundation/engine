package account

import (
	"os"
	"testing"

	"github.com/stvp/assert"
)

func TestImport(t *testing.T) {
	acc, _ := Generate("test")
	account, err := Import(acc.URL.Path, "test", "test")
	assert.Nil(t, err)
	assert.NotEqual(t, account.Address.String(), "")
	Destroy(acc)
}

func TestImportWrongPassword(t *testing.T) {
	acc, _ := Generate("test")
	_, err := Import(acc.URL.Path, "wrongpassword", "test")
	assert.NotNil(t, err)
	Destroy(acc)
}

func TestExport(t *testing.T) {
	acc, _ := Generate("test")
	err := Export(acc, "test", "test", "./TestExport")
	assert.Nil(t, err)
	Destroy(acc)
	os.Remove("./TestExport")
}

func TestExportWrongPath(t *testing.T) {
	acc, _ := Generate("test")
	err := Export(acc, "test", "test", "./")
	assert.NotNil(t, err)
	Destroy(acc)
}
