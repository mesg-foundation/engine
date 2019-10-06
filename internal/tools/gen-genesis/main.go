package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

const mnemonicEntropySize = 256

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randompassword() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var (
	validators    = flag.String("validators", "alice", "list of validator names separated with a comma")
	chainid       = flag.String("chain-id", "mesg-chain", "chain id")
	kbpath        = flag.String("co-kbpath", ".", "cosmos key base path")
	tmpath        = flag.String("tm-path", ".", "tendermint config path")
	gentxfilepath = flag.String("gentx-filepath", "genesistx.json", "genesis transaction file path")
	peersfilepath = flag.String("peers-filepath", "peers", "peers file path")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	validatorNames := strings.Split(*validators, ",")
	validatorNumber := len(validatorNames)

	passwords := make([]string, validatorNumber)
	for i := 0; i < validatorNumber; i++ {
		passwords[i] = randompassword()
	}

	kb, err := cosmos.NewKeybase(*kbpath)
	if err != nil {
		log.Fatalln("creating keybase error:", err)
	}

	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdktypes.RegisterCodec(cdc)
	stakingtypes.RegisterCodec(cdc)

	msgs := []sdktypes.Msg{}
	peers := []string{}
	for i, valName := range validatorNames {
		// read entropy seed straight from crypto.Rand and convert to mnemonic
		entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			log.Fatalln("new entropy error:", err)
		}

		mnemonic, err := bip39.NewMnemonic(entropySeed)
		if err != nil {
			log.Fatalln("new mnemonic error:", err)
		}

		acc, err := kb.CreateAccount(valName, mnemonic, "", passwords[i], 0, 0)
		if err != nil {
			log.Fatalln("create account error:", err)
		}

		if err := os.MkdirAll(filepath.Join(*tmpath, valName, "config"), 0755); err != nil {
			log.Fatalln("mkdir tm config error:", err)
		}

		if err := os.MkdirAll(filepath.Join(*tmpath, valName, "data"), 0755); err != nil {
			log.Fatalln("mkdir tm data error:", err)
		}

		cfg := config.DefaultConfig()
		cfg.SetRoot(filepath.Join(*tmpath, valName))

		nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
		if err != nil {
			log.Fatalln("load/gen node key error:", err)
		}

		fmt.Printf("Validator #%d:\nNode ID: %s\nName: %s\nPassword: %s\nAddress: %s\nMnemonic: %s\nPeer address: %s@%s:26656\n\n", i+1, nodeKey.ID(), valName, passwords[i], acc.GetAddress(), mnemonic, nodeKey.ID(), valName)
		peers = append(peers, fmt.Sprintf("%s@%s:26656", nodeKey.ID(), valName))

		me := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
		msgs = append(msgs, stakingtypes.NewMsgCreateValidator(
			sdktypes.ValAddress(acc.GetAddress()),
			me.GetPubKey(),
			sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.TokensFromConsensusPower(100)),
			stakingtypes.Description{
				Moniker: acc.GetName(),
				Details: "init-validator",
			},
			stakingtypes.NewCommissionRates(
				sdktypes.ZeroDec(),
				sdktypes.ZeroDec(),
				sdktypes.ZeroDec(),
			),
			sdktypes.NewInt(1),
		))
	}

	// generate genesis transaction
	b := cosmos.NewTxBuilder(cdc, 0, 0, kb, *chainid)
	signedMsg, err := b.BuildSignMsg(msgs)
	if err != nil {
		log.Fatalln("build sign msg error:", err)
	}
	// test to sign with only 1 validator
	stdTx := authtypes.NewStdTx(signedMsg.Msgs, signedMsg.Fee, []authtypes.StdSignature{}, signedMsg.Memo)
	for i, valName := range validatorNames {
		stdTx, err = b.SignStdTx(valName, passwords[i], stdTx, true)
		if err != nil {
			log.Fatalln("sign msg create validator error:", err)
		}
	}
	validatorTx := string(cdc.MustMarshalJSON(stdTx))
	if err := ioutil.WriteFile(*gentxfilepath, []byte(validatorTx), 0644); err != nil {
		log.Fatalln("error during writing genesis tx file:", err)
	}

	// peers file
	if err := ioutil.WriteFile(*peersfilepath, []byte(strings.Join(peers, ",")), 0644); err != nil {
		log.Fatalln("error during writing peers file:", err)
	}
}
