package workflow

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/mesg-foundation/core/systemservices/-sources/resolver/proto/core"
	"github.com/mesg-foundation/core/x/xstrings"
	"github.com/sirupsen/logrus"
)

// TODO(ilgooz):
//  * separete VM to its own package?
//  * gracefull shutdown by waiting for continuing task executions.
//  * implement smarter service listening mechanism depending on newly running
//    and deleted workflows. do not listen services if no there are no dependend
//    workflows.
//  * pause the execution of a workflow if one of the services it depends on
//    is not responsive.
//  * do not return from Run if the dependend services are not yet started or
//    listening. use chans as watchers.
//  * log proper messages when a service cannot be listened anymore and implement
//    reconnection mechanisms.

// VM is a virtual machine that runs workflows.
type VM struct {
	// workflows that are running.
	workflows []*WorkflowDocument
	mw        sync.RWMutex

	// activeServices are the started and listening services.
	activeServices map[string]bool
	ma             sync.Mutex

	// core is core client.
	core coreClient
}

// newVM returns a new VM with given core client.
func newVM(core coreClient) *VM {
	return &VM{
		activeServices: make(map[string]bool),
		core:           core,
	}
}

// Run validates a workflow and then runs it if there are no validation errors.
// It starts all the services that workflow depends on and listens for their
// events and results.
func (v *VM) Run(workflow *WorkflowDocument) error {
	if err := v.Validate(workflow); err != nil {
		return err
	}

	inActiveServiceIDs := v.activateServiceIDs(v.getServiceIDs(workflow))
	go func() {
		if err := v.prepareServices(inActiveServiceIDs...); err != nil {
			log.Println(err)
		}
	}()

	v.addWorkflow(workflow)

	logrus.WithFields(logrus.Fields{
		workflowCreationIDKey: workflow.CreationID,
		workflowLogKey: WorkflowLog{
			RunStart:                      true,
			WorkflowID:                    workflow.ID,
			WorkflowCreationID:            workflow.CreationID,
			WorkflowName:                  workflow.Name,
			WorkflowDefinitionName:        workflow.Definition.Name,
			WorkflowDefinitionDescription: workflow.Definition.Description,
		}}).Info("workflow is running")

	return nil
}

// Validate validates a workflow by checking the existence of it's depended services.
// It checks for existence of all the events, results, tasks on the services and
// validates their data schema.
// It checks against infinite cycling executions.
func (v *VM) Validate(workflow *WorkflowDocument) error {
	return nil
}

// prepareServices prepares services by starting and listening events and results on them.
func (v *VM) prepareServices(ids ...string) error {
	ctx := context.Background()

	if err := v.startServices(ctx, ids...); err != nil {
		return err
	}

	if err := v.listenEvents(ctx, ids...); err != nil {
		return err
	}

	return v.listenResults(ctx, ids...)
}

// Terminate terminates a running workflow.
func (v *VM) Terminate(workflowID string) error {
	v.removeWorkflow(v.getWorkflow(workflowID))
	return nil
}

// TerminateAll terminates all running workflows.
func (v *VM) TerminateAll() error {
	if err := v.core.Close(); err != nil {
		return err
	}
	for _, workflow := range v.getWorkflows() {
		if err := v.Terminate(workflow.ID); err != nil {
			return err
		}
	}
	return nil
}

// terminateSiblings terminates running workflows that depends on serviceID.
func (v *VM) terminateSiblings(serviceID string) error {
	for _, workflow := range v.workflows {
		if xstrings.SliceContains(v.getServiceIDs(workflow), serviceID) {
			if err := v.Terminate(workflow.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// executeAll executes do func for multiple ids concurrently and
// returns on the first error and cancels others.
// If there are no errors on the way, it'll wait for all executions to complete.
func (v *VM) executeAll(do func(ctx context.Context, id string) error,
	ctx context.Context, ids ...string) error {
	var (
		lendIDs = len(ids)
		errC    = make(chan error, lendIDs)
	)

	for _, id := range ids {
		go func(id string) {
			errC <- do(ctx, id)
		}(id)
	}

	for i := 0; i < lendIDs; i++ {
		if err := <-errC; err != nil {
			return err
		}
	}

	return nil
}

// startServices starts mesg services with ids.
// TODO(ilgooz) add timeout to ctx.
func (v *VM) startServices(ctx context.Context, ids ...string) error {
	return v.executeAll(func(ctx context.Context, id string) error {
		return v.startService(ctx, id)
	}, ctx, ids...)
}

// listenEvents listens for events on service ids.
func (v *VM) listenEvents(ctx context.Context, ids ...string) error {
	return v.executeAll(func(ctx context.Context, id string) error {
		return v.listenEvent(ctx, id)
	}, ctx, ids...)
}

// listenResults listens for results on service ids.
func (v *VM) listenResults(ctx context.Context, ids ...string) error {
	return v.executeAll(func(ctx context.Context, id string) error {
		return v.listenResult(ctx, id)
	}, ctx, ids...)
}

// startService starts the service with id.
func (v *VM) startService(ctx context.Context, id string) error {
	_, err := v.core.StartService(ctx, &core.StartServiceRequest{
		ServiceID: id,
	})
	return err
}

// Event keeps the event data received from core with some addition info.
type Event struct {
	ServiceID string
	EventKey  string
	EventData string
}

// listenEvent listens events from service id and executes tasks.
func (v *VM) listenEvent(ctx context.Context, id string) error {
	stream, err := v.core.ListenEvent(ctx, &core.ListenEventRequest{
		ServiceID: id,
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			data, err := stream.Recv()
			if err != nil {
				v.terminateSiblings(id)
				return
			}
			go v.executeEvent(Event{
				ServiceID: id,
				EventKey:  data.EventKey,
				EventData: data.EventData,
			})
		}
	}()
	return nil
}

// Result keeps the result data received from core with some addition info.
type Result struct {
	ServiceID     string
	TaskKey       string
	OutputKey     string
	OutputData    string
	ExecutionTags []string
}

// listenEvent listens results from service id and executes tasks.
func (v *VM) listenResult(ctx context.Context, id string) error {
	stream, err := v.core.ListenResult(ctx, &core.ListenResultRequest{
		ServiceID: id,
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			data, err := stream.Recv()
			if err != nil {
				v.terminateSiblings(id)
				return
			}
			go v.executeResult(Result{
				ServiceID:     id,
				TaskKey:       data.TaskKey,
				OutputKey:     data.OutputKey,
				OutputData:    data.OutputData,
				ExecutionTags: data.ExecutionTags,
			})
		}
	}()
	return nil
}

// getServiceIDs returns the dependend service ids of workflow.
func (v *VM) getServiceIDs(workflow *WorkflowDocument) []string {
	var ids []string
	for _, s := range workflow.Definition.Services {
		ids = append(ids, s.ID)
	}
	return ids
}

// getServiceID return the service id where service is named with name on workflow.
func (v *VM) getServiceID(d *WorkflowDocument, name string) string {
	for _, s := range d.Definition.Services {
		if s.Name == name {
			return s.ID
		}
	}
	panic("unreachable")
}

// activateServiceIDs actives the sourceServiceIDs to mark them as the services
// that running workflows depending on and it returns them as inActiveServiceIDs.
// if some of the services in sourceServiceIDs is already active they'll not
// be put to inActiveServiceIDs.
func (v *VM) activateServiceIDs(sourceServiceIDs []string) (inActiveServiceIDs []string) {
	v.ma.Lock()
	defer v.ma.Unlock()
	for _, serviceID := range sourceServiceIDs {
		if _, ok := v.activeServices[serviceID]; !ok {
			inActiveServiceIDs = append(inActiveServiceIDs, serviceID)
			v.activeServices[serviceID] = true
		}
	}
	return inActiveServiceIDs
}

// getWorkflows gets all the running workflows.
func (v *VM) getWorkflows() []*WorkflowDocument {
	v.mw.RLock()
	defer v.mw.RLock()
	return v.workflows
}

// getWorkflow gets a running workflow by it's id.
func (v *VM) getWorkflow(id string) *WorkflowDocument {
	v.mw.RLock()
	defer v.mw.RLock()
	for _, workflow := range v.workflows {
		if workflow.ID == id {
			return workflow
		}
	}
	panic("unreachable")
}

// addWorkflow adds a workflow to running workflows list.
func (v *VM) addWorkflow(workflow *WorkflowDocument) {
	v.mw.RLock()
	defer v.mw.RLock()
	v.workflows = append(v.workflows, workflow)
}

// removeWorkflow removes a workflow from running workflows list.
func (v *VM) removeWorkflow(workflow *WorkflowDocument) {
	v.mw.RLock()
	defer v.mw.RLock()
	for i, w := range v.workflows {
		if w == workflow {
			v.workflows = append(v.workflows[:i], v.workflows[i+1:]...)
			return
		}
	}
}

// executeEvent executes a task depending on event.
func (v *VM) executeEvent(event Event) {
	for _, workflow := range v.getWorkflows() {
		for _, def := range workflow.Definition.Events {
			// execute a task if service id's and event keys matches.
			if event.ServiceID == v.getServiceID(workflow, def.ServiceName) &&
				(def.EventKey == "*" || event.EventKey == def.EventKey) {
				go v.execute(workflow, def, event)
			}
		}
	}
}

// getExecutionData prepares the task inputs with the configs gathered from
// workflow definition and event data.
func (v *VM) getExecutionData(
	services []ServiceDefinition,
	configs []ConfigDefinition,
	def EventDefinition, event Event) (map[string]interface{}, error) {
	// parse event data.
	var eventData map[string]interface{}
	if err := json.Unmarshal([]byte(event.EventData), &eventData); err != nil {
		return nil, err
	}

	// return task outputs as described.
	mappings := def.Map
	if len(mappings) == 0 {
		// use event data as task output.
		return eventData, nil
	}

	parser := dataParser{
		configs: configs,
		data:    eventData,
	}

	data := make(map[string]interface{})
	for _, mapping := range mappings {
		value, err := parser.Parse(mapping.Value)
		if err != nil {
			return nil, err
		}
		data[mapping.Key] = value
	}

	return data, nil
}

func (v *VM) executeResult(result Result) {}

// execute prepares task inputs and executes the task.
func (v *VM) execute(workflow *WorkflowDocument, def EventDefinition, event Event) {
	executionData, err := v.getExecutionData(
		workflow.Definition.Services,
		workflow.Definition.Configs,
		def, event,
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			workflowCreationIDKey: workflow.CreationID,
			eventLogKey: EventLog{
				WorkflowID:         workflow.ID,
				WorkflowCreationID: workflow.CreationID,
				WorkflowName:       workflow.Name,
				ServiceName:        def.ServiceName,
				EventKey:           event.EventKey,
			},
		}).Error("error while parsing event data")
		return
	}

	logrus.WithFields(logrus.Fields{
		workflowCreationIDKey: workflow.CreationID,
		eventLogKey: EventLog{
			WorkflowID:         workflow.ID,
			WorkflowCreationID: workflow.CreationID,
			WorkflowName:       workflow.Name,
			ServiceName:        def.ServiceName,
			EventKey:           event.EventKey,
			ExecutionData:      executionData,
		},
	}).Info("event received")

	if err := v.executeTask(
		v.getServiceID(workflow, def.Execute.ServiceName),
		def.Execute.TaskKey,
		executionData,
	); err != nil {
		logrus.WithFields(logrus.Fields{
			workflowCreationIDKey: workflow.CreationID,
			executionLogKey: ExecutionLog{
				WorkflowID:         workflow.ID,
				WorkflowCreationID: workflow.CreationID,
				WorkflowName:       workflow.Name,
				ServiceName:        def.ServiceName,
				TaskKey:            def.Execute.TaskKey,
			},
		}).Info("error while executing task")
		return
	}

	logrus.WithFields(logrus.Fields{
		workflowCreationIDKey: workflow.CreationID,
		executionLogKey: ExecutionLog{
			WorkflowID:         workflow.ID,
			WorkflowCreationID: workflow.CreationID,
			WorkflowName:       workflow.Name,
			ServiceName:        def.ServiceName,
			TaskKey:            def.Execute.TaskKey,
		},
	}).Info("task executed")
}

// executeTask executes a task on serviceID for taskKey with inputData.
func (v *VM) executeTask(serviceID, taskKey string, inputData map[string]interface{}) error {
	inputDataBytes, err := json.Marshal(inputData)
	if err != nil {
		return err
	}
	_, err = v.core.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID: serviceID,
		TaskKey:   taskKey,
		InputData: string(inputDataBytes),
	})
	return err
}

// list of field keys for log messages.
var (
	workflowCreationIDKey = "workflowCreationID"
	workflowLogKey        = "workflow"
	eventLogKey           = "event"
	executionLogKey       = "execution"
)

// WorkflowLog keeps the workflow log data format.
type WorkflowLog struct {
	RunStart                      bool   `json:"runStart"`
	Deleted                       bool   `json:"deleted"`
	WorkflowID                    string `json:"workflowID"`
	WorkflowCreationID            string `json:"workflowCreationID"`
	WorkflowName                  string `json:"workflowName"`
	WorkflowDefinitionName        string `json:"workflowDefinitionName"`
	WorkflowDefinitionDescription string `json:"workflowDefinitionDescription"`
}

// EventLog keeps the workflow log data format.
type EventLog struct {
	WorkflowID         string      `json:"workflowID"`
	WorkflowCreationID string      `json:"workflowCreationID"`
	WorkflowName       string      `json:"workflowName"`
	ServiceName        string      `json:"serviceName"`
	EventKey           string      `json:"eventKey"`
	ExecutionData      interface{} `json:"executionData"`
}

// ExecutionLog keeps the workflow log data format.
type ExecutionLog struct {
	WorkflowID         string `json:"workflowID"`
	WorkflowCreationID string `json:"workflowCreationID"`
	WorkflowName       string `json:"workflowName"`
	ServiceName        string `json:"serviceName"`
	TaskKey            string `json:"taskKey"`
}
