package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const mnemonicEntropySize = 256

func main() {
	name := flag.String("name", "", "")
	password := flag.String("password", "", "")
	flag.Parse()

	if *name == "" {
		log.Fatalln("name required")
	}

	if *password == "" {
		log.Fatalln("password required")
	}

	// read entropy seed straight from crypto.Rand and convert to mnemonic
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		log.Fatalln(err)
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		log.Fatalln(err)
	}

	acc, err := keys.NewInMemory().CreateAccount(*name, mnemonic, "", *password, 0, 0)
	if err != nil {
		log.Fatalln(err)
	}
	pubkey := acc.GetPubKey().(secp256k1.PubKeySecp256k1)

	fmt.Printf(`MESG_COSMOS_ACCOUNT_PUBKEY=%X
MESG_COSMOS_ACCOUNT_ADDRESS=%s
MESG_COSMOS_ACCOUNT_MNEMONIC="%s"
`,
		pubkey[:],
		acc.GetAddress(),
		mnemonic,
	)
}
