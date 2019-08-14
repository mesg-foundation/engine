package tendermint

import (
	clientkey "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
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

// GenerateAccount creates an account.
func (kb *Keybase) GenerateAccount(name, mnemonic, password string) (keys.Info, error) {
	// read entropy seed straight from crypto.Rand and convert to mnemonic
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return nil, err
	}

	if mnemonic == "" {
		mnemonic, err = bip39.NewMnemonic(entropySeed)
		if err != nil {
			return nil, err
		}
	}
	return kb.CreateAccount(name, mnemonic, "", password, 0, 0)
}
