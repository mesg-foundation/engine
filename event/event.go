package event

import (
	"github.com/gogo/protobuf/types"
	"github.com/mesg-foundation/engine/hash"
)

// Create creates an event eventKey with eventData for service s.
func Create(instanceHash hash.Hash, eventKey string, eventData *types.Struct) *Event {
	e := &Event{
		InstanceHash: instanceHash,
		Key:          eventKey,
		Data:         eventData,
	}
	e.Hash = hash.Dump(e)
	return e
}
