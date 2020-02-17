package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/mesg-foundation/engine/hash"
	instancepb "github.com/mesg-foundation/engine/instance"
	ownershippb "github.com/mesg-foundation/engine/ownership"
	processpb "github.com/mesg-foundation/engine/process"
	runnerpb "github.com/mesg-foundation/engine/runner"
	servicepb "github.com/mesg-foundation/engine/service"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table params.KeyTable) params.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

// BankKeeper module interface.
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// ServiceKeeper module interface.
type ServiceKeeper interface {
	Get(ctx sdk.Context, hash hash.Hash) (*servicepb.Service, error)
}

// InstanceKeeper module interface.
type InstanceKeeper interface {
	Get(ctx sdk.Context, instanceHash hash.Hash) (*instancepb.Instance, error)
}

// RunnerKeeper module interface.
type RunnerKeeper interface {
	Get(ctx sdk.Context, hash hash.Hash) (*runnerpb.Runner, error)
}

// ProcessKeeper module interface.
type ProcessKeeper interface {
	Get(ctx sdk.Context, hash hash.Hash) (*processpb.Process, error)
}

// OwnershipKeeper module interface.
type OwnershipKeeper interface {
	GetOwners(ctx sdk.Context, resourceHash hash.Hash) ([]*ownershippb.Ownership, error)
}
