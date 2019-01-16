package ethwallet

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// From https://github.com/ethereum/go-ethereum/blob/master/accounts/keystore/key.go#L67
type encryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	Id      string              `json:"id"`
	Version int                 `json:"version"`
}

type outputError struct {
	Message string `json:"message"`
}

func OutputError(message string) (string, interface{}) {
	return "error", outputError{
		Message: message,
	}
}
