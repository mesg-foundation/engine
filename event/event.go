package event

import (
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/mesg-foundation/engine/hash"
)

// Create creates an event eventKey with eventData for service s.
func Create(instanceHash hash.Hash, eventKey string, eventData structpb.Struct) *Event {
	e := &Event{
		InstanceHash: instanceHash,
		Key:          eventKey,
		Data:         eventData,
	}
	e.Hash = hash.Dump(e)
	return e
}
