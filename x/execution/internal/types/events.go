package types

// module event types and attributes
const (
	EventType = "execution"

	AttributeKeyHash     = "hash"
	AttributeKeyAddress  = "address"
	AttributeKeyExecutor = "executor"
	AttributeKeyProcess  = "process"
	AttributeKeyInstance = "instance"

	AttributeActionProposed  = "proposed"
	AttributeActionCreated   = "created"
	AttributeActionCompleted = "completed"
	AttributeActionFailed    = "failed"
)
