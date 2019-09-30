package main

import (
	"flag"
	"fmt"
	"log"

	bip39 "github.com/cosmos/go-bip39"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const mnemonicEntropySize = 256

func main() {
	name := flag.String("name", "", "")
	password := flag.String("password", "", "")
	kbpath := flag.String("kbpath", "", "")
	flag.Parse()

	if *name == "" {
		log.Fatalln("name required")
	}

	if *password == "" {
		log.Fatalln("password required")
	}

	if *kbpath == "" {
		log.Fatalln("kbpath required")
	}

	kb, err := cosmos.NewKeybase(*kbpath)
	if err != nil {
		log.Fatalln("creating keybase error:", err)
	}

	// if account exists do not create it
	if info, _ := kb.Get(*name); info != nil {
		fmt.Printf(`MESG_COSMOS_ACCOUNT_PUBKEY=%X
MESG_COSMOS_ACCOUNT_ADDRESS=%s
MESG_COSMOS_ACCOUNT_MNEMONIC="%s"
`,
			info.GetPubKey().Address()[:],
			info.GetAddress(),
			// TODO: How can we restore mnemonic?
			"",
		)
		return
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

	acc, err := kb.CreateAccount(*name, mnemonic, "", *password, 0, 0)
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
