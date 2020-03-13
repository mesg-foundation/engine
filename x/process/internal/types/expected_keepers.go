package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/mesg-foundation/engine/hash"
	instancepb "github.com/mesg-foundation/engine/instance"
	ownershippb "github.com/mesg-foundation/engine/ownership"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table params.KeyTable) params.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

// InstanceKeeper module interface.
type InstanceKeeper interface {
	Get(ctx sdk.Context, instanceHash hash.Hash) (*instancepb.Instance, error)
}

// OwnershipKeeper module interface.
type OwnershipKeeper interface {
	Delete(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash) error
	Set(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash, resource ownershippb.Ownership_Resource) (*ownershippb.Ownership, error)
}

// BankKeeper module interface.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
