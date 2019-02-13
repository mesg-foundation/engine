package provider

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const walletServiceID = "ethwallet"

type ethwalletListOutputSuccess struct {
	Addresses []common.Address `json:"addresses"`
}

type ethwalletCreateInputs struct {
	Passphrase string `json:"passphrase"`
}

type ethwalletCreateOutputSuccess struct {
	Address common.Address `json:"address"`
}

type ethwalletDeleteInputs struct {
	Address    common.Address `json:"address"`
	Passphrase string         `json:"passphrase"`
}

type ethwalletDeleteOutputSuccess struct {
	Address common.Address `json:"address"`
}

type ethwalletExportInputs struct {
	Address    common.Address `json:"address"`
	Passphrase string         `json:"passphrase"`
}

type ethwalletImportInputs struct {
	Account    EncryptedKeyJSONV3 `json:"account"`
	Passphrase string             `json:"passphrase"`
}

type ethwalletImportOutputSuccess struct {
	Address common.Address `json:"address"`
}

type ethwalletSignInputs struct {
	Address     common.Address `json:"address"`
	Passphrase  string         `json:"passphrase"`
	Transaction *Transaction   `json:"transaction"`
}

type ethwalletSignOutputSuccess struct {
	SignedTransaction string `json:"signedTransaction"`
}

type EncryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	ID      string              `json:"id"`
	Version int                 `json:"version"`
}

// Transaction represents created transaction.
type Transaction struct {
	ChainID  int64          `json:"chainID"`
	Nonce    uint64         `json:"nonce"`
	To       common.Address `json:"to"`
	Value    string         `json:"value"`
	Gas      uint64         `json:"gas"`
	GasPrice string         `json:"gasPrice"`
	Data     hexutil.Bytes  `json:"data"`
}
