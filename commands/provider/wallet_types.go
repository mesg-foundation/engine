package provider

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
	Address string      `json:"address"`
	Crypto  interface{} `json:"crypto"`
	ID      string      `json:"id"`
	Version int         `json:"version"`
}

// Transaction represents created transaction.
type Transaction struct {
	ChainID  int64  `json:"chainID"`
	Nonce    uint64 `json:"nonce"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Gas      uint64 `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Data     string `json:"data"`
}
