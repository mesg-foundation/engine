package api

import (
	"context"
	"fmt"

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
		Key:   req.Key,
		Nodes: req.Nodes,
		Edges: req.Edges,
	})
	if err != nil {
		return nil, err
	}
	wf, err = s.sdk.Workflow.Create(wf)
	if err != nil {
		return nil, err
	}
	return &api.CreateWorkflowResponse{Hash: wf.Hash}, nil
}

// Delete deletes service by hash or sid.
func (s *WorkflowServer) Delete(ctx context.Context, request *api.DeleteWorkflowRequest) (*api.DeleteWorkflowResponse, error) {
	return &api.DeleteWorkflowResponse{}, s.sdk.Workflow.Delete(request.Hash)
}

// Get returns service from given hash.
func (s *WorkflowServer) Get(ctx context.Context, req *api.GetWorkflowRequest) (*types.Workflow, error) {
	wf, err := s.sdk.Workflow.Get(req.Hash)
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

<<<<<<< HEAD
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

func fromProtoWorkflowNodes(nodes []*types.Workflow_Node) []workflow.Node {
	res := make([]workflow.Node, len(nodes))
	for i, node := range nodes {
		res[i] = workflow.Node{
			Key:          node.Key,
			InstanceHash: node.InstanceHash,
			TaskKey:      node.TaskKey,
=======
func fromProtoWorkflowNodes(nodes []*types.Workflow_Node) ([]workflow.Node, error) {
	res := make([]workflow.Node, len(nodes))
	for i, node := range nodes {
		switch n := node.Type.(type) {
		case *types.Workflow_Node_Event_:
			hash, err := hash.Decode(n.Event.InstanceHash)
			if err != nil {
				return nil, err
			}
			res[i] = workflow.Event{Key: n.Event.Key, InstanceHash: hash, EventKey: n.Event.EventKey}
		case *types.Workflow_Node_Result_:
			hash, err := hash.Decode(n.Result.InstanceHash)
			if err != nil {
				return nil, err
			}
			res[i] = workflow.Result{Key: n.Result.Key, InstanceHash: hash, TaskKey: n.Result.TaskKey}
		case *types.Workflow_Node_Task_:
			hash, err := hash.Decode(n.Task.InstanceHash)
			if err != nil {
				return nil, err
			}
			res[i] = workflow.Task{InstanceHash: hash, TaskKey: n.Task.TaskKey, Key: n.Task.Key}
		case *types.Workflow_Node_Map_:
			outputs := make([]workflow.Output, len(n.Map.Outputs))
			for j, output := range n.Map.Outputs {
				out := workflow.Output{Key: output.Key}
				switch x := output.Value.(type) {
				case *types.Workflow_Node_Map_Output_Ref:
					out.Ref = &workflow.OutputReference{
						NodeKey: output.GetRef().NodeKey,
						Key:     output.GetRef().Key,
					}
				default:
					return nil, fmt.Errorf("output has unexpected type %T", x)
				}
				outputs[j] = out
			}
			res[i] = workflow.Map{Key: n.Map.Key, Outputs: outputs}
		default:
			return nil, fmt.Errorf("node has unexpected type %T", n)
>>>>>>> dev
		}
	}
	return res
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
<<<<<<< HEAD
	nodes := fromProtoWorkflowNodes(wf.Nodes)
	edges, err := fromProtoWorkflowEdges(wf.Edges)
	if err != nil {
		return nil, err
	}
	trigger := workflow.Trigger{
		InstanceHash: wf.Trigger.InstanceHash,
		NodeKey:      wf.Trigger.NodeKey,
		Filters:      fromProtoFilters(wf.Trigger.Filters),
	}
	switch x := wf.Trigger.Key.(type) {
	case *types.Workflow_Trigger_EventKey:
		trigger.EventKey = x.EventKey
	case *types.Workflow_Trigger_TaskKey:
		trigger.TaskKey = x.TaskKey
	default:
		return nil, fmt.Errorf("workflow trigger key has unexpected type %T", x)
	}
=======
	nodes, err := fromProtoWorkflowNodes(wf.Nodes)
	if err != nil {
		return nil, err
	}
>>>>>>> dev
	return &workflow.Workflow{
		Key: wf.Key,
		Graph: workflow.Graph{
			Nodes: nodes,
			Edges: fromProtoWorkflowEdges(wf.Edges),
		},
	}, nil
}

func toProtoWorkflowNodes(nodes []workflow.Node) []*types.Workflow_Node {
	res := make([]*types.Workflow_Node, len(nodes))
	for i, node := range nodes {
<<<<<<< HEAD
		res[i] = &types.Workflow_Node{
			Key:          node.Key,
			InstanceHash: node.InstanceHash,
			TaskKey:      node.TaskKey,
=======
		protoNode := types.Workflow_Node{}
		switch n := node.(type) {
		case *workflow.Result:
			protoNode.Type = &types.Workflow_Node_Result_{
				Result: &types.Workflow_Node_Result{
					Key:          n.Key,
					InstanceHash: n.InstanceHash.String(),
					TaskKey:      n.TaskKey,
				},
			}
		case *workflow.Event:
			protoNode.Type = &types.Workflow_Node_Event_{
				Event: &types.Workflow_Node_Event{
					Key:          n.Key,
					InstanceHash: n.InstanceHash.String(),
					EventKey:     n.EventKey,
				},
			}
		case *workflow.Task:
			protoNode.Type = &types.Workflow_Node_Task_{
				Task: &types.Workflow_Node_Task{
					Key:          n.Key,
					InstanceHash: n.InstanceHash.String(),
					TaskKey:      n.TaskKey,
				},
			}
		case *workflow.Map:
			outputs := make([]*types.Workflow_Node_Map_Output, len(n.Outputs))
			for j, output := range n.Outputs {
				out := &types.Workflow_Node_Map_Output{Key: output.Key}
				if output.Ref != nil {
					out.Value = &types.Workflow_Node_Map_Output_Ref{
						Ref: &types.Workflow_Node_Map_Output_Reference{
							NodeKey: output.Ref.NodeKey,
							Key:     output.Ref.Key,
						},
					}
				}
				outputs[j] = out
			}

			protoNode.Type = &types.Workflow_Node_Map_{
				Map: &types.Workflow_Node_Map{
					Key:     n.Key,
					Outputs: outputs,
				},
			}
>>>>>>> dev
		}
		res[i] = &protoNode
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
<<<<<<< HEAD
	trigger := &types.Workflow_Trigger{
		InstanceHash: wf.Trigger.InstanceHash,
		Filters:      toProtoFilters(wf.Trigger.Filters),
		NodeKey:      wf.Trigger.NodeKey,
	}
	if wf.Trigger.TaskKey != "" {
		trigger.Key = &types.Workflow_Trigger_TaskKey{TaskKey: wf.Trigger.TaskKey}
	}
	if wf.Trigger.EventKey != "" {
		trigger.Key = &types.Workflow_Trigger_EventKey{EventKey: wf.Trigger.EventKey}
	}
	return &types.Workflow{
		Hash:    wf.Hash,
		Key:     wf.Key,
		Trigger: trigger,
		Nodes:   toProtoWorkflowNodes(wf.Nodes),
		Edges:   toProtoWorkflowEdges(wf.Edges),
=======
	return &types.Workflow{
		Hash:  wf.Hash.String(),
		Key:   wf.Key,
		Nodes: toProtoWorkflowNodes(wf.Nodes),
		Edges: toProtoWorkflowEdges(wf.Edges),
>>>>>>> dev
	}
}

func toProtoWorkflows(workflows []*workflow.Workflow) []*types.Workflow {
	wfs := make([]*types.Workflow, len(workflows))
	for i, wf := range workflows {
		wfs[i] = toProtoWorkflow(wf)
	}
	return wfs
}
