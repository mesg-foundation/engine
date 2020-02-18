package event

import (
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// New creates an event eventKey with eventData for service s.
func New(instanceHash hash.Hash, eventKey string, eventData *types.Struct) *Event {
	e := &Event{
		InstanceHash: instanceHash,
		Key:          eventKey,
		Data:         eventData,
	}
	e.Hash = hash.Dump(e)
	return e
}
