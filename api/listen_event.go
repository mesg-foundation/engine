package api

// ListenEventFilter is a filter func for filtering events.
type ListenEventFilter func(*EventListener)

// ListenEventKeyFilter returns an eventKey filter.
func ListenEventKeyFilter(eventKey string) ListenEventFilter {
	return func(ln *EventListener) {
		ln.eventKey = eventKey
	}
}

// ListenEvent listens events matches with eventFilter on serviceID.
func (a *API) ListenEvent(serviceID string, filters ...ListenEventFilter) (*EventListener, error) {
	l := newEventListener(a, filters...)
	return l, l.listen(serviceID)
}
