package mesg

// Taskable represents implementation of task executor.
type Taskable interface {
	// Name returns the task key.
	Name() (key string)

	// Execute called when a task execution request arrived with the same key
	// returned by Name().
	// key is the output key and data is the output data of task.
	Execute(execution *Execution) (key string, data Data)
}

// defaultTask implements Task.
type defaultTask struct {
	name    string
	handler func(execution *Execution) (key string, data Data)
}

func (t *defaultTask) Name() string {
	return t.name
}

func (t *defaultTask) Execute(execution *Execution) (key string, data Data) {
	return t.handler(execution)
}

// Task creates an executable task for given task name, handler called when a
// matching task execution request arrived.
func Task(name string, handler func(*Execution) (key string, data Data)) Taskable {
	return &defaultTask{name, handler}
}
