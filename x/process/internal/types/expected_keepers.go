package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	instancepb "github.com/mesg-foundation/engine/instance"
	ownershippb "github.com/mesg-foundation/engine/ownership"
)

// InstanceKeeper module interface.
type InstanceKeeper interface {
	Get(ctx sdk.Context, instanceHash hash.Hash) (*instancepb.Instance, error)
}

// OwnershipKeeper module interface.
type OwnershipKeeper interface {
	Delete(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash) error
	Set(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash, resource ownershippb.Ownership_Resource, resourceAddress sdk.AccAddress) (*ownershippb.Ownership, error)
}
