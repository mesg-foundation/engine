package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

const mnemonicEntropySize = 256

var (
	names     = []string{"a", "b", "c", "d", "e"}
	passwords = []string{"a", "b", "c", "d", "e"}
)

var (
	vno     = flag.Int("vno", 2, "validator numbers")
	chainid = flag.String("chain-id", "test-net", "chain id")
	kbpath  = flag.String("co-kbpath", ".", "cosmos key base path")
	tmpath  = flag.String("tm-path", ".", "tendermint config path")
)

func main() {
	flag.Parse()

	kb, err := cosmos.NewKeybase(*kbpath)
	if err != nil {
		log.Fatalln("creating keybase error:", err)
	}

	// if account exists do not create it
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdktypes.RegisterCodec(cdc)
	stakingtypes.RegisterCodec(cdc)

	msgs := []sdktypes.Msg{}
	for i := 0; i < *vno; i++ {
		// read entropy seed straight from crypto.Rand and convert to mnemonic
		entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			log.Fatalln(err)
		}

		mnemonic, err := bip39.NewMnemonic(entropySeed)
		if err != nil {
			log.Fatalln(err)
		}

		acc, err := kb.CreateAccount(names[i], mnemonic, "", passwords[i], 0, 0)
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
		if err := os.MkdirAll(filepath.Join(*tmpath, names[i], "config"), 0755); err != nil {
			panic(err)
		}

		if err := os.MkdirAll(filepath.Join(*tmpath, names[i], "data"), 0755); err != nil {
			panic(err)
		}

		cfg := config.DefaultConfig()
		cfg.SetRoot(filepath.Join(*tmpath, names[i]))

		nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
		if err != nil {
			panic(err)
		}
		fmt.Printf("NODE_PUBKEY=%s@localhost:26656\n", nodeKey.ID())

		me := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())
		pk := me.GetPubKey().(ed25519.PubKeyEd25519)
		fmt.Printf("VALIDATOR_PUBKEY=%X\n", pk[:])
		msgs = append(msgs, stakingtypes.NewMsgCreateValidator(
			sdktypes.ValAddress(acc.GetAddress()),
			pk,
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

	b := cosmos.NewTxBuilder(cdc, 0, 0, kb, *chainid)
	signedMsg, err := b.BuildSignMsg(msgs)
	if err != nil {
		panic(err)
	}
	stdTx := authtypes.NewStdTx(signedMsg.Msgs, signedMsg.Fee, []authtypes.StdSignature{}, signedMsg.Memo)
	for i := 0; i < *vno; i++ {
		stdTx, err = b.SignStdTx(names[i], passwords[i], stdTx, true)
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		log.Fatalln("sign msg create validator error:", err)
	}

	fmt.Printf("MESG_COSMOS_GENESISVALIDATORTX='%s'\n", string(cdc.MustMarshalJSON(stdTx)))
}
