package applicationtest

// EventListen holds information about an event listen request.
type EventListen struct {
	serviceID string
	event     string
}

// ServiceID returns the id of service that events are emitted from.
func (l *EventListen) ServiceID() string {
	return l.serviceID
}

// EventFilter returns the event name.
func (l *EventListen) EventFilter() string {
	return l.event
}

// ResultListen holds information about a result listen request.
type ResultListen struct {
	serviceID string
	key       string
	task      string
}

// ServiceID returns the id of service that results are emitted from.
func (l *ResultListen) ServiceID() string {
	return l.serviceID
}

// KeyFilter returns the output key name.
func (l *ResultListen) KeyFilter() string {
	return l.key
}

// TaskFilter returns the task key name.
func (l *ResultListen) TaskFilter() string {
	return l.task
}
