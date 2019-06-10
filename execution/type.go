package execution

// Status stores the state of an execution
type Status int

// Status for an execution
// Created    => The execution is created but not yet processed
// InProgress => The execution is being processed
// Completed  => The execution is completed
const (
	Created Status = iota + 1
	InProgress
	Completed
	Failed
)

// Execution stores all informations about executions.
type Execution struct {
	Hash        []byte                 `hash:"-"`
	ParentHash  []byte                 `hash:"name:parentHash"`
	EventID     string                 `hash:"name:eventID"`
	Status      Status                 `hash:"-"`
	ServiceHash string                 `hash:"name:serviceHash"`
	TaskKey     string                 `hash:"name:taskKey"`
	Tags        []string               `hash:"name:tags"`
	Inputs      map[string]interface{} `hash:"name:inputs"`
	Outputs     map[string]interface{} `hash:"-"`
	Error       string                 `hash:"-"`
}
