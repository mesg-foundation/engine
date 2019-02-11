package ethwallet

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// From https://github.com/ethereum/go-ethereum/blob/master/accounts/keystore/key.go#L67
type encryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	ID      string              `json:"id"`
	Version int                 `json:"version"`
}
