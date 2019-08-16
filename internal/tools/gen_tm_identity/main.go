package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

func main() {
	path := flag.String("path", "", "")
	name := flag.String("name", "", "")
	port := flag.String("port", "", "")

	flag.Parse()

	if path == nil || *path == "" {
		panic("path is required")
	}

	if name == nil || *name == "" {
		panic("name is required")
	}

	if port == nil || *port == "" {
		panic("port is required")
	}

	if err := os.MkdirAll(filepath.Join(*path, "config"), 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(*path, "data"), 0755); err != nil {
		panic(err)
	}

	cfg := config.DefaultConfig()
	cfg.SetRoot(*path)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		panic(err)
	}
	fmt.Printf("NODE_PUBKEY=%s@%s:%s\n", nodeKey.ID(), *name, *port)

	me := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
	pk := me.GetPubKey().(ed25519.PubKeyEd25519)
	fmt.Printf("VALIDATOR_PUBKEY=%X\n", pk[:])
}
