package main

import (
	"github.com/mesg-foundation/engine/systemservices/ethwallet/ethwallet"
)

func main() {
	ethwallet, err := ethwallet.New()
	if err != nil {
		panic(err)
	}
	if err = ethwallet.Listen(); err != nil {
		panic(err)
	}
}
