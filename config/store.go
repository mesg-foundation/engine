package config

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// Store is the object that let you control all your store of data including the accounts
var Store *keystore.KeyStore

func init() {
	Store = keystore.NewKeyStore(AccountDirectory, keystore.StandardScryptN, keystore.StandardScryptP)
}
