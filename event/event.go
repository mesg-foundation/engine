package event

import (
	"github.com/mesg-foundation/core/hash"
)

// Event stores all informations about Events.
type Event struct {
	Hash         hash.Hash
	InstanceHash hash.Hash
	Key          string
	Data         map[string]interface{}
}

// Create creates an event eventKey with eventData for service s.
func Create(instanceHash hash.Hash, eventKey string, eventData map[string]interface{}) *Event {
	e := &Event{
		InstanceHash: instanceHash,
		Key:          eventKey,
		Data:         eventData,
	}
	e.Hash = hash.Dump(e)
	return e
}
