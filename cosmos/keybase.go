package cosmos

import (
	clientkey "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/keyerror"
	bip39 "github.com/cosmos/go-bip39"
)

const mnemonicEntropySize = 256

// Keybase is a standard cosmos keybase.
type Keybase struct {
	keys.Keybase
}

// NewKeybase initializes a filesystem keybase at a particular dir.
func NewKeybase(dir string) (*Keybase, error) {
	kb, err := clientkey.NewKeyBaseFromDir(dir)
	if err != nil {
		return nil, err
	}
	return &Keybase{kb}, nil
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
