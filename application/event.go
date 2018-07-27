package application

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/api/core"
)

// Event is a MESG event.
type Event struct {
	Key  string
	data string
}

// Decode decodes event data into out.
func (e *Event) Decode(out interface{}) error {
	return json.Unmarshal([]byte(e.data), out)
}

// EventListener is a MESG event listener.
type EventListener struct {
	app *Application

	// event is the actual event to listen for.
	event string

	//eventServiceID is the service id of where event is emitted.
	eventServiceID string

	// task is the actual task that will be executed.
	task string

	// taskServiceID is the service id of target task.
	taskServiceID string

	// filterFunc is a func that returns a boolean value to check
	// if the task should be executed or not.
	filterFunc func(*Event) bool

	mapData Data

	// provideFunc is a func that returns input data of task.
	mapFunc func(*Event) Data

	// cancel cancels listening for upcoming events.
	cancel context.CancelFunc
}

// EventOption is the configuration func of EventListener.
type EventOption func(*EventListener)

// EventFilterOption returns a new option to filter events by name.
// Default is all(*).
func EventFilterOption(event string) EventOption {
	return func(l *EventListener) {
		l.event = event
	}
}

// WhenEvent creates an EventListener for serviceID.
func (a *Application) WhenEvent(serviceID string, options ...EventOption) *EventListener {
	l := &EventListener{
		app:            a,
		eventServiceID: serviceID,
		event:          "*",
	}
	for _, option := range options {
		option(l)
	}
	return l
}

// FilterFunc expects the returned value to be true to do task execution.
func (l *EventListener) FilterFunc(fn func(*Event) bool) *EventListener {
	l.filterFunc = fn
	return l
}

// Data is piped as the input data to task.
type Data interface{}

// Map sets data as the input data to task.
func (l *EventListener) Map(data Data) *EventListener {
	l.mapData = data
	return l
}

// MapFunc sets the returned data as the input data of task.
// You can dynamically produce input values for task over event data.
func (l *EventListener) MapFunc(fn func(*Event) Data) *EventListener {
	l.mapFunc = fn
	return l
}

// Execute executes task on serviceID.
func (l *EventListener) Execute(serviceID, task string) (*Stream, error) {
	l.taskServiceID = serviceID
	l.task = task
	stream := &Stream{
		Executions: make(chan *Execution, 0),
		Err:        make(chan error, 0),
	}
	if err := l.app.startServices(l.eventServiceID, serviceID); err != nil {
		return nil, err
	}
	cancel, err := l.listen(stream)
	if err != nil {
		return nil, err
	}
	stream.cancel = cancel
	return stream, nil
}

// Listen starts listening for events.
func (l *EventListener) listen(stream *Stream) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	resp, err := l.app.client.ListenEvent(ctx, &core.ListenEventRequest{
		ServiceID:   l.eventServiceID,
		EventFilter: l.event,
	})
	if err != nil {
		return cancel, err
	}
	go l.readStream(stream, resp)
	return cancel, nil
}

func (l *EventListener) readStream(stream *Stream, resp core.Core_ListenEventClient) {
	for {
		data, err := resp.Recv()
		if err != nil {
			stream.Err <- err
			return
		}
		event := &Event{
			Key:  data.EventKey,
			data: data.EventData,
		}
		go l.execute(stream, event)
	}
}

func (l *EventListener) execute(stream *Stream, event *Event) {
	if l.filterFunc != nil && !l.filterFunc(event) {
		return
	}

	var data Data
	switch {
	case l.mapData != nil:
		data = l.mapData
	case l.mapFunc != nil:
		data = l.mapFunc(event)
	default:
		if err := event.Decode(&data); err != nil {
			stream.Executions <- &Execution{
				Err: err,
			}
			return
		}
	}

	executionID, err := l.app.Execute(l.taskServiceID, l.task, data)
	stream.Executions <- &Execution{
		ID:  executionID,
		Err: err,
	}
}
