package workflow

import (
	"github.com/mesg-foundation/core/sdk"
)

// Workflow ...
type Workflow struct {
	sdk *sdk.SDK
}

// New ...
func New(sdk *sdk.SDK) Workflow {
	return Workflow{sdk}
}

// Start ...
func (w Workflow) Start() {
	go w.listenExecutions()
	go w.listenEvents()
}

func (w Workflow) listenExecutions() {
	for {
		execution := <-w.sdk.Execution.Stream()
		for _, srv := range w.findServicesByTrigger(execution.ServiceHash, "executionFinished") {
			for _, wf := range w.findWorkflowsByTrigger(srv, execution.ServiceHash, "executionFinished") {
				w.CreateExecutionChain(wf, event)
			}
		}

		execution.Workflow
	}
}

func (w Workflow) listenEvents() {
	for {
		event := <-w.sdk.Event.Stream()
		for _, srv := range w.findServicesByTrigger(event.ServiceHash, event.Key) {
			for _, wf := range w.findWorkflowsByTrigger(srv, event.ServiceHash, event.Key) {
				w.CreateExecutionChain(wf, event)
			}
		}
	}
}

func findServicesByTrigger(serviceHash string, event string) []service.Service {
	return []
}

func findWorkflowsByTrigger(srv service.Service, serviceHash string, event string) []service.Workflow {

}