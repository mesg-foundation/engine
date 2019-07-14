package event

import (
	"github.com/mesg-foundation/engine/hash"
)

// a 44 length constant hash
const engineEventHash = "EngineEventEngineEventEngineEventEngineEvent"

// EngineEventType type to describe engine events
// Engine events doesn't need the same validation that normal events need
// These events are used internally only and cannot be emitted by services
type EngineEventType string

const (
	// EngineAPIExecution This event is triggered when the `execution/create` API is called
	EngineAPIExecution EngineEventType = "engine-api-execution"
)

// EngineEvent creates an engine event
func EngineEvent(evtType EngineEventType, data map[string]interface{}) *Event {
	instanceHash, err := hash.Decode(engineEventHash)
	if err != nil {
		// This panics because this is a developer mistake
		panic("engineEventHash should be codable")
	}
	return Create(instanceHash, string(evtType), data)
}
