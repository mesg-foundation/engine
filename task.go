package mesg

// Taskable represents implementation of task executor.
type Taskable interface {
	// Key returns the task key.
	Key() string

	// Execute called when a task execution request arrived with the same key
	// returned by Name().
	// outputKey is the output key and data is the output data of task.
	Execute(execution *Execution) (outputKey string, outputData Data)
}

// defaultTask implements Task.
type defaultTask struct {
	name    string
	handler func(execution *Execution) (key string, data Data)
}

func (t *defaultTask) Key() string {
	return t.name
}

func (t *defaultTask) Execute(execution *Execution) (outputKey string, outputData Data) {
	return t.handler(execution)
}

// Task creates an executable task for given task key, handler called when a
// matching task execution request arrived.
func Task(key string, handler func(*Execution) (key string, data Data)) Taskable {
	return &defaultTask{key, handler}
}
