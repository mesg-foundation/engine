package tendermint

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// TxBuilder implements a transaction context created in SDK modules.
type txBuilder struct {
	authtypes.TxBuilder
}

// NewTxBuilder returns a new initialized TxBuilder.
func NewTxBuilder(cdc *codec.Codec, accNumber, accSeq uint64, kb *Keybase, chainID string) txBuilder {
	return txBuilder{
		authtypes.NewTxBuilder(
			authutils.GetTxEncoder(cdc),
			accNumber,
			accSeq,
			flags.DefaultGasLimit,
			flags.DefaultGasAdjustment,
			true,
			chainID,
			"",
			sdktypes.NewCoins(),
			sdktypes.DecCoins{},
		).WithKeybase(kb),
	}
}

// DefaultSignStdTx appends a signature to a StdTx with default gas limit and not fees.
func (b txBuilder) DefaultSignStdTx(msg sdktypes.Msg, accountName, accountPassword string) (authtypes.StdTx, error) {
	fees := authtypes.NewStdFee(flags.DefaultGasLimit, sdktypes.NewCoins())
	stdTx := authtypes.NewStdTx([]sdktypes.Msg{msg}, fees, []authtypes.StdSignature{}, "")
	return b.SignStdTx(accountName, accountPassword, stdTx, false)
}
