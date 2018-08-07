package service

// EventNotFoundError is an error when an event cannot be found in a service.
type EventNotFoundError struct {
	Service  *Service
	EventKey string
}

func (e *EventNotFoundError) Error() string {
	return "Event '" + e.EventKey + "' not found in service '" + e.Service.Name + "'"
}

// InvalidEventDataError is an error when the data of an event are not valid.
type InvalidEventDataError struct {
	Event *Event
	Key   string
	Data  map[string]interface{}
}

func (e *InvalidEventDataError) Error() string {
	errorString := "Data of event '" + e.Key + "' is invalid"
	for _, warning := range e.Event.Validate(e.Data) {
		errorString = errorString + ". " + warning.String()
	}
	return errorString
}

// TaskNotFoundError is an error when a task cannot be found in a service.
type TaskNotFoundError struct {
	Service *Service
	TaskKey string
}

func (e *TaskNotFoundError) Error() string {
	return "Task '" + e.TaskKey + "' not found in service '" + e.Service.Name + "'"
}

// InvalidTaskInputError is an error when the inputs of a task are not valid.
type InvalidTaskInputError struct {
	Task   *Task
	Key    string
	Inputs map[string]interface{}
}

func (e *InvalidTaskInputError) Error() string {
	errorString := "Inputs of task '" + e.Key + "' are invalid"
	for _, warning := range e.Task.Validate(e.Inputs) {
		errorString = errorString + ". " + warning.String()
	}
	return errorString
}

// OutputNotFoundError is an error when a service doesn't contain a specific output.
type OutputNotFoundError struct {
	Service   *Service
	OutputKey string
}

func (e *OutputNotFoundError) Error() string {
	return "Output '" + e.OutputKey + "' not found in service '" + e.Service.Name + "'"
}

// InvalidOutputDataError is an error when the outputs for one task result are not valid.
type InvalidOutputDataError struct {
	Output *Output
	Key    string
	Data   map[string]interface{}
}

func (e *InvalidOutputDataError) Error() string {
	errorString := "Outputs of task '" + e.Key + "' are invalid"
	for _, warning := range e.Output.Validate(e.Data) {
		errorString = errorString + ". " + warning.String()
	}
	return errorString
}
