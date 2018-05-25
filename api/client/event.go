package client

// When creates a new workflow based on a service's event
func When(event *Event) (wf *Workflow) {
	wf = &Workflow{
		Event: event,
	}
	return wf
}
