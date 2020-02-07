package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type forceCheckTxProxyDecorator struct {
	d sdk.AnteDecorator
}

func newForceCheckTxProxyDecorator(d sdk.AnteDecorator) forceCheckTxProxyDecorator {
	return forceCheckTxProxyDecorator{
		d: d,
	}
}

func (d forceCheckTxProxyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	ctx = ctx.WithIsCheckTx(false).WithValue("mesgHackIsCheckTx", ctx.IsCheckTx())
	n := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		checkTx, ok := ctx.Value("mesgHackIsCheckTx").(bool)
		if !ok {
			return ctx, fmt.Errorf("mesgHackIsCheckTx is not a bool")
		}
		return next(ctx.WithIsCheckTx(checkTx), tx, simulate)
	}
	return d.d.AnteHandle(ctx, tx, simulate, n)
}
