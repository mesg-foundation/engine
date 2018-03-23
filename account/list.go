package account

import "github.com/ethereum/go-ethereum/common"

// List all available accounts on this computer
func List() (accountList []*Account) {
	// TODO add real list
	accountList = []*Account{
		&Account{Name: "Test1", Address: common.Address{0}},
		&Account{Name: "Test2", Address: common.Address{1}},
	}
	return
}
