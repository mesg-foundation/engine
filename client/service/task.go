package service

// Taskable represents implementation of task executor.
type Taskable interface {
	// Key returns the task key.
	Key() string

	// Execute called when a task execution request arrived with the same key
	// returned by Name().
	Execute(execution *Execution) (interface{}, error)
}

// defaultTask implements Task.
type defaultTask struct {
	name    string
	handler func(execution *Execution) (interface{}, error)
}

func (t *defaultTask) Key() string {
	return t.name
}

func (t *defaultTask) Execute(execution *Execution) (interface{}, error) {
	return t.handler(execution)
}

// Task creates an executable task for given task key, handler called when a
// matching task execution request arrived.
func Task(key string, handler func(*Execution) (interface{}, error)) Taskable {
	return &defaultTask{key, handler}
}
