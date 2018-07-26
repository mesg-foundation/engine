package execution

// NotInQueueError is an error when trying to access an execution that doesn't exists
type NotInQueueError struct {
	ID    string
	Queue string
}

func (e *NotInQueueError) Error() string {
	return "Execution '" + e.ID + "' not found in queue '" + e.Queue + "'"
}
