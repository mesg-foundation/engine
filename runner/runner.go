package runner

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/hash"
	"github.com/tendermint/tendermint/crypto"
)

// New returns a new execution.
func New(owner string, instanceHash hash.Hash) (*Runner, error) {
	run := &Runner{
		Owner:        owner,
		InstanceHash: instanceHash,
	}
	run.Hash = hash.Dump(run)
	run.Address = sdk.AccAddress(crypto.AddressHash(run.Hash))
	return run, xvalidator.Struct(run)
}
