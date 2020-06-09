package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	servicepb "github.com/mesg-foundation/engine/service"
)

// ServiceKeeper module interface.
type ServiceKeeper interface {
	Get(ctx sdk.Context, hash hash.Hash) (*servicepb.Service, error)
}
