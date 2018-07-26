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
