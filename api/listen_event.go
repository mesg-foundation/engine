package api

// ListenEvent listens events matches with eventFilter on serviceID.
func (a *API) ListenEvent(serviceID, eventFilter string) (*EventListener, error) {
	l := newEventListener(a)
	return l, l.listen(serviceID, eventFilter)
}
