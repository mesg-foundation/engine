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
		Nodes:   req.Nodes,
		Edges:   req.Edges,
	})
	if err != nil {
		return nil, err
	}
	wf, err = s.sdk.Workflow.Create(wf)
	if err != nil {
		return nil, err
	}
	return &api.CreateWorkflowResponse{Hash: wf.Hash.String()}, nil
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

func fromProtoWorkflowNodes(nodes []*types.Workflow_Node) ([]workflow.Node, error) {
	res := make([]workflow.Node, len(nodes))
	for i, node := range nodes {
		instanceHash, err := hash.Decode(node.InstanceHash)
		if err != nil {
			return nil, err
		}
		res[i] = workflow.Node{
			Key:          node.Key,
			InstanceHash: instanceHash,
			TaskKey:      node.TaskKey,
		}
	}
	return res, nil
}

func fromProtoWorkflowEdges(edges []*types.Workflow_Edge) []workflow.Edge {
	res := make([]workflow.Edge, len(edges))
	for i, edge := range edges {
		res[i] = workflow.Edge{
			Src: edge.Src,
			Dst: edge.Dst,
		}
	}
	return res
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
	nodes, err := fromProtoWorkflowNodes(wf.Nodes)
	if err != nil {
		return nil, err
	}
	return &workflow.Workflow{
		Key: wf.Key,
		Trigger: workflow.Trigger{
			Type:         triggerType,
			InstanceHash: instanceHash,
			Key:          wf.Trigger.Key,
			InitialNode:  wf.Trigger.InitialNode,
			Filters:      fromProtoFilters(wf.Trigger.Filters),
		},
		Nodes: nodes,
		Edges: fromProtoWorkflowEdges(wf.Edges),
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

func toProtoWorkflowNodes(nodes []workflow.Node) []*types.Workflow_Node {
	res := make([]*types.Workflow_Node, len(nodes))
	for i, node := range nodes {
		res[i] = &types.Workflow_Node{
			Key:          node.Key,
			InstanceHash: node.InstanceHash.String(),
			TaskKey:      node.TaskKey,
		}
	}
	return res
}

func toProtoWorkflowEdges(edges []workflow.Edge) []*types.Workflow_Edge {
	res := make([]*types.Workflow_Edge, len(edges))
	for i, edge := range edges {
		res[i] = &types.Workflow_Edge{
			Src: edge.Src,
			Dst: edge.Dst,
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
			InitialNode:  wf.Trigger.InitialNode,
		},
		Nodes: toProtoWorkflowNodes(wf.Nodes),
		Edges: toProtoWorkflowEdges(wf.Edges),
	}
}

func toProtoWorkflows(workflows []*workflow.Workflow) []*types.Workflow {
	wfs := make([]*types.Workflow, len(workflows))
	for i, wf := range workflows {
		wfs[i] = toProtoWorkflow(wf)
	}
	return wfs
}
