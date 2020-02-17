package reciver

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

var Reciver sdk.AccAddress

func InitReciver(kb keys.Keybase) {
	// read entropy seed straight from crypto.Rand and convert to mnemonic
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		panic(err)
	}
	mem, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		panic(err)
	}
	info, err := kb.CreateAccount("reciver", mem, "", "reciver", keys.CreateHDPath(1, 1).String(), keys.Secp256k1)
	if err != nil {
		panic(err)
	}
	Reciver = info.GetAddress()
}
