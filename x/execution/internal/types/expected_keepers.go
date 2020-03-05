package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	instancepb "github.com/mesg-foundation/engine/instance"
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
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	InputOutputCoins(ctx sdk.Context, inputs []bank.Input, outputs []bank.Output) error
}

// ServiceKeeper module interface.
type ServiceKeeper interface {
	Get(ctx sdk.Context, hash sdk.AccAddress) (*servicepb.Service, error)
}

// InstanceKeeper module interface.
type InstanceKeeper interface {
	Get(ctx sdk.Context, instanceHash sdk.AccAddress) (*instancepb.Instance, error)
}

// RunnerKeeper module interface.
type RunnerKeeper interface {
	Get(ctx sdk.Context, hash sdk.AccAddress) (*runnerpb.Runner, error)
}

// ProcessKeeper module interface.
type ProcessKeeper interface {
	Get(ctx sdk.Context, hash sdk.AccAddress) (*processpb.Process, error)
}
