package api

import (
	"context"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/sdk"
	"github.com/mesg-foundation/engine/workflow"
)

// WorkflowServer is the type to aggregate all Service APIs.
type WorkflowServer struct {
	sdk *sdk.SDK
}

// NewWorkflowServer creates a new WorkflowServer.
func NewWorkflowServer(sdk *sdk.SDK) *WorkflowServer {
	return &WorkflowServer{sdk: sdk}
}

// Create creates a new service from definition.
func (s *WorkflowServer) Create(ctx context.Context, req *api.CreateWorkflowRequest) (*api.CreateWorkflowResponse, error) {
	wf, err := fromProtoWorkflow(&types.Workflow{
		Key:     req.Key,
		Trigger: req.Trigger,
		Tasks:   req.Tasks,
	})
	if err != nil {
		return nil, err
	}
	srv, err := s.sdk.Workflow.Create(wf)
	if err != nil {
		return nil, err
	}
	return &api.CreateWorkflowResponse{Hash: srv.Hash.String()}, nil
}

// Delete deletes service by hash or sid.
func (s *WorkflowServer) Delete(ctx context.Context, request *api.DeleteWorkflowRequest) (*api.DeleteWorkflowResponse, error) {
	hash, err := hash.Decode(request.Hash)
	if err != nil {
		return nil, err
	}
	return &api.DeleteWorkflowResponse{}, s.sdk.Workflow.Delete(hash)
}

// Get returns service from given hash.
func (s *WorkflowServer) Get(ctx context.Context, req *api.GetWorkflowRequest) (*types.Workflow, error) {
	hash, err := hash.Decode(req.Hash)
	if err != nil {
		return nil, err
	}

	wf, err := s.sdk.Workflow.Get(hash)
	if err != nil {
		return nil, err
	}
	return toProtoWorkflow(wf), nil
}

// List returns all workflows.
func (s *WorkflowServer) List(ctx context.Context, req *api.ListWorkflowRequest) (*api.ListWorkflowResponse, error) {
	workflows, err := s.sdk.Workflow.List()
	if err != nil {
		return nil, err
	}
	wfs := toProtoWorkflows(workflows)
	return &api.ListWorkflowResponse{
		Workflows: wfs,
	}, nil
}

func fromProtoFilters(filters []*types.Workflow_Trigger_Filter) []*workflow.TriggerFilter {
	fs := make([]*workflow.TriggerFilter, len(filters))
	for i, filter := range filters {
		var predicate workflow.Predicate
		// switch filter.Predicate {
		if filter.Predicate == types.Workflow_Trigger_Filter_EQ {
			predicate = workflow.EQ
		}
		fs[i] = &workflow.TriggerFilter{
			Key:       filter.Key,
			Predicate: predicate,
			Value:     filter.Value,
		}
	}
	return fs
}

func fromProtoWorkflowTasks(tasks []*types.Workflow_Task) ([]*workflow.Task, error) {
	res := make([]*workflow.Task, len(tasks))
	for i, task := range tasks {
		instanceHash, err := hash.Decode(task.InstanceHash)
		if err != nil {
			return nil, err
		}
		res[i] = &workflow.Task{
			InstanceHash: instanceHash,
			TaskKey:      task.TaskKey,
		}
	}
	return res, nil
}

func fromProtoWorkflow(wf *types.Workflow) (*workflow.Workflow, error) {
	var triggerType workflow.TriggerType
	switch wf.Trigger.Type {
	case types.Workflow_Trigger_Result:
		triggerType = workflow.RESULT
	case types.Workflow_Trigger_Event:
		triggerType = workflow.EVENT
	}
	instanceHash, err := hash.Decode(wf.Trigger.InstanceHash)
	if err != nil {
		return nil, err
	}
	tasks, err := fromProtoWorkflowTasks(wf.Tasks)
	if err != nil {
		return nil, err
	}
	return &workflow.Workflow{
		Key: wf.Key,
		Trigger: &workflow.Trigger{
			Type:         triggerType,
			InstanceHash: instanceHash,
			Key:          wf.Trigger.Key,
			Filters:      fromProtoFilters(wf.Trigger.Filters),
		},
		Tasks: tasks,
	}, nil
}

func toProtoFilters(filters []*workflow.TriggerFilter) []*types.Workflow_Trigger_Filter {
	fs := make([]*types.Workflow_Trigger_Filter, len(filters))
	for i, filter := range filters {
		var predicate types.Workflow_Trigger_Filter_Predicate
		// switch filter.Predicate {
		if filter.Predicate == workflow.EQ {
			predicate = types.Workflow_Trigger_Filter_EQ
		}
		fs[i] = &types.Workflow_Trigger_Filter{
			Key:       filter.Key,
			Predicate: predicate,
			Value:     filter.Value.(string),
		}
	}
	return fs
}

func toProtoWorkflowTasks(tasks []*workflow.Task) []*types.Workflow_Task {
	res := make([]*types.Workflow_Task, len(tasks))
	for i, task := range tasks {
		res[i] = &types.Workflow_Task{
			InstanceHash: task.InstanceHash.String(),
			TaskKey:      task.TaskKey,
		}
	}
	return res
}

func toProtoWorkflow(wf *workflow.Workflow) *types.Workflow {
	var triggerType types.Workflow_Trigger_Type
	switch wf.Trigger.Type {
	case workflow.EVENT:
		triggerType = types.Workflow_Trigger_Event
	case workflow.RESULT:
		triggerType = types.Workflow_Trigger_Result
	}
	return &types.Workflow{
		Hash: wf.Hash.String(),
		Key:  wf.Key,
		Trigger: &types.Workflow_Trigger{
			Type:         triggerType,
			InstanceHash: wf.Trigger.InstanceHash.String(),
			Key:          wf.Trigger.Key,
			Filters:      toProtoFilters(wf.Trigger.Filters),
		},
		Tasks: toProtoWorkflowTasks(wf.Tasks),
	}
}

func toProtoWorkflows(workflows []*workflow.Workflow) []*types.Workflow {
	wfs := make([]*types.Workflow, len(workflows))
	for i, wf := range workflows {
		wfs[i] = toProtoWorkflow(wf)
	}
	return wfs
}
