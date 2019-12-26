package cosmos

import (
	"sync"

	clientkey "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/keyerror"
	"github.com/cosmos/cosmos-sdk/types"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto"
)

const (
	// AccNumber is the account number of the account in keybase.
	AccNumber = 0
	// AccIndex is the account index of the account in keybase.
	AccIndex = 0

	mnemonicEntropySize = 256
)

// Keybase is a standard cosmos keybase.
type Keybase struct {
	kb keys.Keybase

	mx sync.Mutex
}

// NewKeybase initializes a filesystem keybase at a particular dir.
func NewKeybase(dir string) (*Keybase, error) {
	kb, err := clientkey.NewKeyBaseFromDir(dir)
	if err != nil {
		return nil, err
	}
	return &Keybase{kb: kb}, nil
}

// NewMnemonic returns a new mnemonic phrase.
func (kb *Keybase) NewMnemonic() (string, error) {
	// read entropy seed straight from crypto.Rand and convert to mnemonic
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropySeed)
}

// Exist checks if the account exists.
func (kb *Keybase) Exist(name string) (bool, error) {
	_, err := kb.Get(name)
	if keyerror.IsErrKeyNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// List is a lock protected version of keys.List
func (kb *Keybase) List() ([]keys.Info, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.List()
}

// Get is a lock protected version of keys.Get
func (kb *Keybase) Get(name string) (keys.Info, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.Get(name)
}

// GetByAddress is a lock protected version of keys.GetByAddress
func (kb *Keybase) GetByAddress(address types.AccAddress) (keys.Info, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.GetByAddress(address)
}

// Delete is a lock protected version of keys.Delete
func (kb *Keybase) Delete(name, passphrase string, skipPass bool) error {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.Delete(name, passphrase, skipPass)
}

// Sign is a lock protected version of keys.Sign
func (kb *Keybase) Sign(name, passphrase string, msg []byte) ([]byte, crypto.PubKey, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.Sign(name, passphrase, msg)
}

// CreateMnemonic is a lock protected version of keys.CreateMnemonic
func (kb *Keybase) CreateMnemonic(name string, language keys.Language, passwd string, algo keys.SigningAlgo) (keys.Info, string, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.CreateMnemonic(name, language, passwd, algo)
}

// CreateAccount is a lock protected version of keys.CreateAccount
func (kb *Keybase) CreateAccount(name, mnemonic, bip39Passwd, encryptPasswd string, account, index uint32) (keys.Info, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.CreateAccount(name, mnemonic, bip39Passwd, encryptPasswd, account, index)
}

// Derive is a lock protected version of keys.Derive
func (kb *Keybase) Derive(name, mnemonic, bip39Passwd, encryptPasswd string, params hd.BIP44Params) (keys.Info, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.Derive(name, mnemonic, bip39Passwd, encryptPasswd, params)
}

// CreateLedger is a lock protected version of keys.CreateLedger
func (kb *Keybase) CreateLedger(name string, algo keys.SigningAlgo, hrp string, account, index uint32) (keys.Info, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.CreateLedger(name, algo, hrp, account, index)
}

// CreateOffline is a lock protected version of keys.CreateOffline
func (kb *Keybase) CreateOffline(name string, pubkey crypto.PubKey) (keys.Info, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.CreateOffline(name, pubkey)
}

// CreateMulti is a lock protected version of keys.CreateMulti
func (kb *Keybase) CreateMulti(name string, pubkey crypto.PubKey) (keys.Info, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.CreateMulti(name, pubkey)
}

// Update is a lock protected version of keys.Update
func (kb *Keybase) Update(name, oldpass string, getNewpass func() (string, error)) error {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.Update(name, oldpass, getNewpass)
}

// Import is a lock protected version of keys.Import
func (kb *Keybase) Import(name, armor string) error {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.Import(name, armor)
}

// ImportPrivKey is a lock protected version of keys.ImportPrivKey
func (kb *Keybase) ImportPrivKey(name, armor, passphrase string) error {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.ImportPrivKey(name, armor, passphrase)
}

// ImportPubKey is a lock protected version of keys.ImportPubKey
func (kb *Keybase) ImportPubKey(name, armor string) (err error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.ImportPubKey(name, armor)
}

// Export is a lock protected version of keys.Export
func (kb *Keybase) Export(name string) (armor string, err error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.Export(name)
}

// ExportPubKey is a lock protected version of keys.ExportPubKey
func (kb *Keybase) ExportPubKey(name string) (armor string, err error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.ExportPubKey(name)
}

// ExportPrivKey is a lock protected version of keys.ExportPrivKey
func (kb *Keybase) ExportPrivKey(name, decryptPassphrase, encryptPassphrase string) (armor string, err error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.ExportPrivKey(name, decryptPassphrase, encryptPassphrase)
}

// ExportPrivateKeyObject is a lock protected version of keys.ExportPrivateKeyObject
func (kb *Keybase) ExportPrivateKeyObject(name string, passphrase string) (crypto.PrivKey, error) {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	return kb.kb.ExportPrivateKeyObject(name, passphrase)
}

// CloseDB is a lock protected version of keys.CloseDB
func (kb *Keybase) CloseDB() {
	kb.mx.Lock()
	defer kb.mx.Unlock()
	kb.kb.CloseDB()
}
