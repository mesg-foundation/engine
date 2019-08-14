package tendermint

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sirupsen/logrus"
)

type Signer struct {
	cdc     *codec.Codec
	kb      *Keybase
	chainID string
}

func NewSigner(cdc *codec.Codec, kb *Keybase, chainID string) *Signer {
	return &Signer{
		cdc:     cdc,
		kb:      kb,
		chainID: chainID,
	}
}

func (s *Signer) signTransaction(msg sdktypes.Msg, accountName, accountPassword string) (authtypes.StdTx, error) {
	fees := authtypes.NewStdFee(flags.DefaultGasLimit, sdktypes.NewCoins())
	gasPrices := sdktypes.DecCoins{}
	stdTx := authtypes.NewStdTx([]sdktypes.Msg{msg}, fees, []authtypes.StdSignature{}, "")

	txBldr := authtypes.NewTxBuilder(
		authutils.GetTxEncoder(s.cdc),
		0,
		0,
		flags.DefaultGasLimit,
		flags.DefaultGasAdjustment,
		true,
		s.chainID,
		"",
		sdktypes.NewCoins(),
		gasPrices,
	).WithKeybase(s.kb)

	return txBldr.SignStdTx(accountName, accountPassword, stdTx, false)
}

func (s *Signer) signTransaction2(msg sdktypes.Msg, accountName, accountPassword string) ([]byte, error) {
	logrus.Warning("chainId", s.chainID)
	txBldr := authtypes.NewTxBuilder(
		authutils.GetTxEncoder(s.cdc),
		100,
		0,
		flags.DefaultGasLimit,
		flags.DefaultGasAdjustment,
		true,
		s.chainID,
		"",
		sdktypes.NewCoins(),
		sdktypes.DecCoins{},
	).WithKeybase(s.kb)

	return txBldr.BuildAndSign(accountName, accountPassword, []sdktypes.Msg{msg})
}
