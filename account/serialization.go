package account

import (
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/mesg-foundation/core/config"
)

// Import a file and return a Account object
func Import(filePath string, password string, newPassword string) (acc accounts.Account, err error) {
	content, err := ioutil.ReadFile(filePath)
	acc, err = config.Store.Import(content, password, newPassword)
	return
}

// Export a given account to a file
func Export(acc accounts.Account, password string, newPassword string, path string) (err error) {
	content, err := config.Store.Export(acc, password, newPassword)
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(content)
	if err != nil {
		return
	}
	return
}
