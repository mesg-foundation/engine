package event

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/tendermint/tendermint/crypto"
)

// New creates an event eventKey with eventData for service s.
func New(instanceHash sdk.AccAddress, eventKey string, eventData *types.Struct) *Event {
	e := &Event{
		InstanceHash: instanceHash,
		Key:          eventKey,
		Data:         eventData,
	}
	e.Hash = sdk.AccAddress(crypto.AddressHash([]byte(e.HashSerialize())))
	return e
}
