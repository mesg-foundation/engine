package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	ownershippb "github.com/mesg-foundation/engine/ownership"
)

// OwnershipKeeper module interface.
type OwnershipKeeper interface {
	Set(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash, resource ownershippb.Ownership_Resource, resourceAddress sdk.AccAddress) (*ownershippb.Ownership, error)
}
