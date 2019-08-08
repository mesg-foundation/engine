package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

func main() {
	root := os.ExpandEnv("$HOME/.mesg/tendermint")

	if err := os.MkdirAll(filepath.Join(root, "config"), 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "data"), 0755); err != nil {
		panic(err)
	}

	cfg := config.DefaultConfig()
	cfg.SetRoot(root)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		panic(err)
	}
	fmt.Printf("MESG_TENDERMINT_P2P_SEEDS=%s@%s:%s\n", nodeKey.ID(), "engine", "26656")

	me := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
	pk := me.GetPubKey().(ed25519.PubKeyEd25519)
	fmt.Printf("MESG_TENDERMINT_VALIDATORPUBKEY=%X\n", pk[:])
}
