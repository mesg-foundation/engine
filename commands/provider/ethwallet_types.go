package provider

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const walletServiceID = "ethwallet"

type ethwalletListOutputSuccess struct {
	Addresses []string `json:"addresses"`
}

type ethwalletCreateInputs struct {
	Passphrase string `json:"passphrase"`
}

type ethwalletCreateOutputSuccess struct {
	Address string `json:"address"`
}

type ethwalletDeleteInputs struct {
	Address    string `json:"address"`
	Passphrase string `json:"passphrase"`
}

type ethwalletDeleteOutputSuccess struct {
	Address string `json:"address"`
}

type ethwalletExportInputs struct {
	Address    string `json:"address"`
	Passphrase string `json:"passphrase"`
}

type ethwalletImportInputs struct {
	Account    EncryptedKeyJSONV3 `json:"account"`
	Passphrase string             `json:"passphrase"`
}

type ethwalletImportFromPrivateKeyInputs struct {
	PrivateKey string `json:"privateKey"`
	Passphrase string `json:"passphrase"`
}

type ethwalletImportOutputSuccess struct {
	Address string `json:"address"`
}

type ethwalletSignInputs struct {
	Address     string       `json:"address"`
	Passphrase  string       `json:"passphrase"`
	Transaction *Transaction `json:"transaction"`
}

type ethwalletSignOutputSuccess struct {
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
