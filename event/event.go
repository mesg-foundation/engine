package event

import (
	"github.com/mesg-foundation/core/hash"
)

// Event stores all informations about Events.
type Event struct {
	InstanceHash hash.Hash
	Key          string
	Data         map[string]interface{}
}

// Create creates an event eventKey with eventData for service s.
func Create(instanceHash hash.Hash, eventKey string, eventData map[string]interface{}) *Event {
	return &Event{
		InstanceHash: instanceHash,
		Key:          eventKey,
		Data:         eventData,
	}
}
