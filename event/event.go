package event

import (
	"github.com/mesg-foundation/engine/cosmos/address"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/tendermint/tendermint/crypto"
)

// New creates an event eventKey with eventData for service s.
func New(instanceHash address.InstAddress, eventKey string, eventData *types.Struct) *Event {
	e := &Event{
		InstanceHash: instanceHash,
		Key:          eventKey,
		Data:         eventData,
	}
	e.Hash = address.EventAddress(crypto.AddressHash([]byte(e.HashSerialize())))
	return e
}
