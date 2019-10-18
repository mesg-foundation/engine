package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/x/xos"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/privval"
	"gopkg.in/yaml.v2"
)

const mnemonicEntropySize = 256

const (
	vno        = 1
	testdir    = "test"
	nameprefix = "testuser"
	chainid    = "mesg-testnet"
)

var (
	kbpath = filepath.Join(testdir, "cosmos")

	cdc = func() *codec.Codec {
		cdc := codec.New()
		codec.RegisterCrypto(cdc)
		sdktypes.RegisterCodec(cdc)
		stakingtypes.RegisterCodec(cdc)
		return cdc
	}()
)

func genTestData(t *testing.T) {
	kb, err := cosmos.NewKeybase(kbpath)
	if err != nil {
		t.Fatal("creating keybase error:", err)
	}

	msgs := []sdktypes.Msg{}
	for i := 0; i < vno; i++ {
		namepass := fmt.Sprintf("%s-%d", nameprefix, i)
		tmpath := filepath.Join(testdir, namepass, "tendermint")

		// read entropy seed straight from crypto.Rand and convert to mnemonic
		entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			t.Fatal("new entropy error:", err)
		}

		mnemonic, err := bip39.NewMnemonic(entropySeed)
		if err != nil {
			t.Fatal("new mnemonic error:", err)
		}

		acc, err := kb.CreateAccount(namepass, mnemonic, "", namepass, 0, 0)
		if err != nil {
			t.Fatal("create account error:", err)
		}

		if err := os.MkdirAll(filepath.Join(tmpath, "config"), 0755); err != nil {
			t.Fatal("mkdir tm config error:", err)
		}

		if err := os.MkdirAll(filepath.Join(tmpath, "data"), 0755); err != nil {
			t.Fatal("mkdir tm data error:", err)
		}

		cfg := tmconfig.DefaultConfig()
		cfg.SetRoot(tmpath)

		// nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
		// if err != nil {
		// 	t.Fatal("load/gen node key error:", err)
		// }

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
	b := cosmos.NewTxBuilder(cdc, 0, 0, kb, chainid)
	signedMsg, err := b.BuildSignMsg(msgs)
	if err != nil {
		t.Fatal("build sign msg error:", err)
	}

	// test to sign with only 1 validator
	stdTx := authtypes.NewStdTx(signedMsg.Msgs, signedMsg.Fee, []authtypes.StdSignature{}, signedMsg.Memo)
	for i := 0; i < vno; i++ {
		namepass := fmt.Sprintf("%s-%d", nameprefix, i)

		stdTx, err = b.SignStdTx(namepass, namepass, stdTx, true)
		if err != nil {
			t.Fatal("sign msg create validator error:", err)
		}

		// copy cosmos keybase for every user
		userkbpath := filepath.Join(testdir, namepass, "cosmos")
		if err := xos.Copy(kbpath, userkbpath); err != nil {
			t.Fatal("copy keybase error:", err)
		}

		// write genesis file for each user
		gentxfilepath := filepath.Join(testdir, namepass, "genesis-tx.json")
		if err := ioutil.WriteFile(gentxfilepath, cdc.MustMarshalJSON(stdTx), 0644); err != nil {
			t.Fatal("error during writing genesis tx file:", err)
		}

		// generate config file for each user
		c, err := config.Default()
		if err != nil {
			t.Fatal("default config error:", err)
		}
		c.Path = filepath.Join(testdir, namepass)
		c.Log.Level = "fatal"

		out, err := yaml.Marshal(c)
		if err != nil {
			t.Fatal("yaml unmarshal error:", err)
		}

		configpath := filepath.Join(testdir, namepass, "config.yml")
		if err := ioutil.WriteFile(configpath, out, 0644); err != nil {
			t.Fatal("error during writing config file:", err)
		}
	}
}
