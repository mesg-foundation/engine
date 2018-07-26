package application

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/api/core"
)

// Result is MESG result event.
type Result struct {
	TaskKey     string
	OutputKey   string
	data        string
	executionID string
}

// Decode decodes result data into out.
func (e *Result) Decode(out interface{}) error {
	return json.Unmarshal([]byte(e.data), out)
}

// ResultListener is a MESG result event listener.
type ResultListener struct {
	app *Application

	// resultTask is the actual event to listen for.
	resultTask string

	//resultServiceID is the service id of where result is emitted.
	resultServiceID string

	// outputKey is the output key to listen for.
	outputKey string

	// task is the actual task that will be executed.
	task string

	// taskServiceID is the service id of target task.
	taskServiceID string

	// filterFunc is a func that returns a boolean value to check
	// if the task should be executed or not.
	filterFunc func(*Result) bool

	mapData Data

	// mapFunc is a func that returns input data of task.
	mapFunc func(*Result) Data

	// cancel cancels listening for upcoming events.
	cancel context.CancelFunc
}

type ResultOption func(*ResultListener)

func TaskFilterOption(task string) ResultOption {
	return func(l *ResultListener) {
		l.resultTask = task
	}
}

func OutputKeyFilterOption(key string) ResultOption {
	return func(l *ResultListener) {
		l.outputKey = key
	}
}

// WhenResult creates a ResultListener for serviceID.
func (a *Application) WhenResult(serviceID string, options ...ResultOption) *ResultListener {
	l := &ResultListener{
		app:             a,
		resultServiceID: serviceID,
		resultTask:      "*",
		outputKey:       "*",
	}
	for _, option := range options {
		option(l)
	}
	return l
}

func (l *ResultListener) FilterFunc(fn func(*Result) bool) *ResultListener {
	l.filterFunc = fn
	return l
}

func (l *ResultListener) Map(data Data) *ResultListener {
	l.mapData = data
	return l
}

func (l *ResultListener) MapFunc(fn func(*Result) Data) *ResultListener {
	l.mapFunc = fn
	return l
}

func (l *ResultListener) Execute(serviceID, task string) (*Stream, error) {
	l.taskServiceID = serviceID
	l.task = task
	stream := &Stream{
		Executions: make(chan *Execution, 0),
		Err:        make(chan error, 0),
	}
	cancel, err := l.listen(stream)
	if err != nil {
		return nil, err
	}
	stream.cancel = cancel
	return stream, nil
}

func (l *ResultListener) listen(stream *Stream) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	resp, err := l.app.client.ListenResult(ctx, &core.ListenResultRequest{
		ServiceID:    l.resultServiceID,
		TaskFilter:   l.resultTask,
		OutputFilter: l.outputKey,
	})
	if err != nil {
		return cancel, err
	}
	go l.readStream(stream, resp)
	return cancel, nil
}

func (l *ResultListener) readStream(stream *Stream, resp core.Core_ListenResultClient) {
	for {
		data, err := resp.Recv()
		if err != nil {
			stream.Err <- err
			return
		}
		result := &Result{
			TaskKey:     data.TaskKey,
			OutputKey:   data.OutputKey,
			data:        data.OutputData,
			executionID: data.ExecutionID,
		}
		go l.execute(stream, result)
	}
}

func (l *ResultListener) execute(stream *Stream, result *Result) {
	if l.filterFunc != nil && !l.filterFunc(result) {
		return
	}

	var data Data
	switch {
	case l.mapData != nil:
		data = l.mapData
	case l.mapFunc != nil:
		data = l.mapFunc(result)
	default:
		if err := result.Decode(&data); err != nil {
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
