package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	instancepb "github.com/mesg-foundation/engine/instance"
	processpb "github.com/mesg-foundation/engine/process"
	runnerpb "github.com/mesg-foundation/engine/runner"
	servicepb "github.com/mesg-foundation/engine/service"
)

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
	List(ctx sdk.Context) ([]*runnerpb.Runner, error)
}

// ProcessKeeper module interface.
type ProcessKeeper interface {
	Get(ctx sdk.Context, hash hash.Hash) (*processpb.Process, error)
}

// CreditKeeper module interface.
type CreditKeeper interface {
	Sub(ctx sdk.Context, address sdk.AccAddress, amount sdk.Int) (sdk.Int, error)
}
