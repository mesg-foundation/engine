package ownership

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/hash"
)

// New returns a new ownership and validate it.
func New(owner string, resource Ownership_Resource, resourceHash hash.Hash, resourceAddress sdk.AccAddress) (*Ownership, error) {
	own := &Ownership{
		Owner:           owner,
		Resource:        resource,
		ResourceHash:    resourceHash,
		ResourceAddress: resourceAddress,
	}
	own.Hash = hash.Dump(own)
	return own, xvalidator.Struct(own)
}
