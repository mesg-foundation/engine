package event

// EngineEventType type to describe engine events
// Engine events doesn't need the same validation that normal events need
// These events are used internally only and cannot be emitted by services
type EngineEventType string

const (
	// EngineAPIExecution This event is triggered when the `execution/create` API is called
	EngineAPIExecution EngineEventType = "engine-api-execution"
)

// EngineEvent creates an engine event.
// Engine event does not have any instance hash and are prefixed by mesg:
func EngineEvent(evtType EngineEventType, data map[string]interface{}) *Event {
	return Create(nil, "mesg:"+string(evtType), data)
}
