package cosmos

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ForceCheckTxProxyDecorator defines a proxy decorator that force the ctx.IsCheckTx to false.
type ForceCheckTxProxyDecorator struct {
	d sdk.AnteDecorator
}

// NewForceCheckTxProxyDecorator returns a proxy decorator that force the ctx.IsCheckTx to false.
func NewForceCheckTxProxyDecorator(d sdk.AnteDecorator) ForceCheckTxProxyDecorator {
	return ForceCheckTxProxyDecorator{
		d: d,
	}
}

// AnteHandle creates a proxy on top of the passed AnteDecorator AnteHandle.
func (d ForceCheckTxProxyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
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
