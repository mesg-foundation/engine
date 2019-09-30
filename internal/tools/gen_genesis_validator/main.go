package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func main() {
	name := flag.String("name", "", "")
	password := flag.String("password", "", "")
	chainid := flag.String("chain-id", "", "")
	pubkey := flag.String("pubkey", "", "")
	kbpath := flag.String("kbpath", "", "")
	flag.Parse()

	if *name == "" {
		log.Fatalln("flag name required")
	}

	if *password == "" {
		log.Fatalln("flag password required")
	}

	if *chainid == "" {
		log.Fatalln("flag chain-id required")
	}

	if *pubkey == "" {
		log.Fatalln("flag pubkey required")
	}

	if *kbpath == "" {
		log.Fatalln("flag kbpath required")
	}

	dec, err := hex.DecodeString(*pubkey)
	if err != nil {
		log.Fatalln("validator public key decode error:", err)
	}

	var validatorPubKey ed25519.PubKeyEd25519
	copy(validatorPubKey[:], dec)

	kb, err := cosmos.NewKeybase(*kbpath)
	if err != nil {
		log.Fatalln("creating keybase error:", err)
	}

	// fetch account
	account, err := kb.Get(*name)
	if err != nil {
		log.Fatalln("getting user from keybase error:", err)
	}

	msg := stakingtypes.NewMsgCreateValidator(
		sdktypes.ValAddress(account.GetAddress()),
		validatorPubKey,
		sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.TokensFromConsensusPower(100)),
		stakingtypes.Description{
			Moniker: *name,
			Details: "create-first-validator",
		},
		stakingtypes.NewCommissionRates(
			sdktypes.ZeroDec(),
			sdktypes.ZeroDec(),
			sdktypes.ZeroDec(),
		),
		sdktypes.NewInt(1),
	)

	msgs := []sdktypes.Msg{msg}
	accounts, err := kb.List()
	if err != nil {
		log.Fatal("list account error:", err)
	}

	for _, acc := range accounts {
		if acc.GetName() == *name {
			continue
		}
		msgs = append(msgs, stakingtypes.NewMsgCreateValidator(
			sdktypes.ValAddress(acc.GetAddress()),
			validatorPubKey,
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

	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdktypes.RegisterCodec(cdc)
	stakingtypes.RegisterCodec(cdc)

	tx, err := cosmos.NewTxBuilder(cdc, 0, 0, kb, *chainid).BuildAndSignStdTx(msgs, *name, *password)
	if err != nil {
		log.Fatalln("sign msg create validator error:", err)
	}

	fmt.Printf("MESG_COSMOS_GENESISVALIDATORTX='%s'\n", string(cdc.MustMarshalJSON(tx)))
}
