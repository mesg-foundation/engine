package tendermint

import (
	clientkey "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/go-bip39"
)

const mnemonicEntropySize = 256

type keybase struct {
	keys.Keybase
}

// NewFSKeybas initializes a filesystem keybase at a particular dir.
func NewFSKeybase(dir string) (*keybase, error) {
	kb, err := clientkey.NewKeyBaseFromDir(dir)
	if err != nil {
		return nil, err
	}
	return &keybase{kb}, nil
}

// CreateAccount creates an account.
func (kb *keybase) GenerateAccount(name, mnemonic, password string) (keys.Info, error) {
	// read entropy seed straight from crypto.Rand and convert to mnemonic
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return nil, err
	}

	if mnemonic == "" {
		mnemonic, err = bip39.NewMnemonic(entropySeed[:])
		if err != nil {
			return nil, err
		}
	}
	return kb.CreateAccount(name, mnemonic, "", password, 0, 0)
}
