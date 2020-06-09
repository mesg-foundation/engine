package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	ownershippb "github.com/mesg-foundation/engine/ownership"
)

// InstanceKeeper module interface.
type InstanceKeeper interface {
	FetchOrCreate(ctx sdk.Context, serviceHash hash.Hash, envHash hash.Hash) (*instance.Instance, error)
}

// OwnershipKeeper module interface.
type OwnershipKeeper interface {
	Set(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash, resource ownershippb.Ownership_Resource, resourceAddress sdk.AccAddress) (*ownershippb.Ownership, error)
	Delete(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash) error
}
