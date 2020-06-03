package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper module interface.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}
