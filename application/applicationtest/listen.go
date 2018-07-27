package applicationtest

type EventListen struct {
	serviceID string
	event     string
}

func (l *EventListen) ServiceID() string {
	return l.serviceID
}

func (l *EventListen) EventFilter() string {
	return l.event
}

type ResultListen struct {
	serviceID string
	key       string
	task      string
}

func (l *ResultListen) ServiceID() string {
	return l.serviceID
}

func (l *ResultListen) KeyFilter() string {
	return l.key
}

func (l *ResultListen) TaskFilter() string {
	return l.task
}
