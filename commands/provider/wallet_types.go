package provider

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const walletServiceID = "ethwallet"

type walletListOutputSuccess struct {
	Addresses []string `json:"addresses"`
}

type walletCreateInputs struct {
	Passphrase string `json:"passphrase"`
}

type walletCreateOutputSuccess struct {
	Address string `json:"address"`
}

type walletDeleteInputs struct {
	Address    string `json:"address"`
	Passphrase string `json:"passphrase"`
}

type walletDeleteOutputSuccess struct {
	Address string `json:"address"`
}

type walletExportInputs struct {
	Address    string `json:"address"`
	Passphrase string `json:"passphrase"`
}

type walletImportInputs struct {
	Account    EncryptedKeyJSONV3 `json:"account"`
	Passphrase string             `json:"passphrase"`
}

type walletImportFromPrivateKeyInputs struct {
	PrivateKey string `json:"privateKey"`
	Passphrase string `json:"passphrase"`
}

type walletImportOutputSuccess struct {
	Address string `json:"address"`
}

type walletSignInputs struct {
	Address     string       `json:"address"`
	Passphrase  string       `json:"passphrase"`
	Transaction *Transaction `json:"transaction"`
}

type walletSignOutputSuccess struct {
	SignedTransaction string `json:"signedTransaction"`
}

type EncryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"` // TODO: remove all type from go-ethereum
	ID      string              `json:"id"`
	Version int                 `json:"version"`
}

// Transaction represents created transaction.
type Transaction struct {
	ChainID  int64         `json:"chainID"`
	Nonce    uint64        `json:"nonce"`
	To       string        `json:"to"`
	Value    string        `json:"value"`
	Gas      uint64        `json:"gas"`
	GasPrice string        `json:"gasPrice"`
	Data     hexutil.Bytes `json:"data"` // TODO: remove all type from go-ethereum
}
