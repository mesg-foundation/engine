package api

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/engine/filter"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/sdk"
)

// ProcessServer is the type to aggregate all Service APIs.
type ProcessServer struct {
	sdk *sdk.SDK
}

// NewProcessServer creates a new ProcessServer.
func NewProcessServer(sdk *sdk.SDK) *ProcessServer {
	return &ProcessServer{sdk: sdk}
}

// Create creates a new service from definition.
func (s *ProcessServer) Create(ctx context.Context, req *api.CreateProcessRequest) (*api.CreateProcessResponse, error) {
	wf, err := fromProtoProcess(&types.Process{
		Key:   req.Key,
		Nodes: req.Nodes,
		Edges: req.Edges,
	})
	if err != nil {
		return nil, err
	}
	wf, err = s.sdk.Process.Create(wf)
	if err != nil {
		return nil, err
	}
	return &api.CreateProcessResponse{Hash: wf.Hash}, nil
}

// Delete deletes service by hash or sid.
func (s *ProcessServer) Delete(ctx context.Context, request *api.DeleteProcessRequest) (*api.DeleteProcessResponse, error) {
	return &api.DeleteProcessResponse{}, s.sdk.Process.Delete(request.Hash)
}

// Get returns service from given hash.
func (s *ProcessServer) Get(ctx context.Context, req *api.GetProcessRequest) (*types.Process, error) {
	wf, err := s.sdk.Process.Get(req.Hash)
	if err != nil {
		return nil, err
	}
	return toProtoProcess(wf), nil
}

// List returns all processes.
func (s *ProcessServer) List(ctx context.Context, req *api.ListProcessRequest) (*api.ListProcessResponse, error) {
	processes, err := s.sdk.Process.List()
	if err != nil {
		return nil, err
	}
	wfs := toProtoProcesses(processes)
	return &api.ListProcessResponse{
		Processes: wfs,
	}, nil
}

func fromProtoProcessNodes(nodes []*types.Process_Node) ([]process.Node, error) {
	res := make([]process.Node, len(nodes))
	for i, node := range nodes {
		switch n := node.Type.(type) {
		case *types.Process_Node_Event_:
			res[i] = process.Event{Key: n.Event.Key, InstanceHash: n.Event.InstanceHash, EventKey: n.Event.EventKey}
		case *types.Process_Node_Result_:
			res[i] = process.Result{Key: n.Result.Key, InstanceHash: n.Result.InstanceHash, TaskKey: n.Result.TaskKey}
		case *types.Process_Node_Task_:
			res[i] = process.Task{InstanceHash: n.Task.InstanceHash, TaskKey: n.Task.TaskKey, Key: n.Task.Key}
		case *types.Process_Node_Map_:
			outputs := make([]process.Output, len(n.Map.Outputs))
			for j, output := range n.Map.Outputs {
				out := process.Output{Key: output.Key}
				switch x := output.Value.(type) {
				case *types.Process_Node_Map_Output_Ref:
					out.Ref = &process.OutputReference{
						NodeKey: output.GetRef().NodeKey,
						Key:     output.GetRef().Key,
					}
				default:
					return nil, fmt.Errorf("output has unexpected type %T", x)
				}
				outputs[j] = out
			}
			res[i] = process.Map{Key: n.Map.Key, Outputs: outputs}
		case *types.Process_Node_Filter_:
			conditions := make([]filter.Condition, len(n.Filter.Conditions))
			for j, condition := range n.Filter.Conditions {
				cond := filter.Condition{Key: condition.Key, Value: condition.Value}
				switch condition.Predicate {
				case types.Process_Node_Filter_Condition_EQ:
					cond.Predicate = filter.EQ
				default:
					return nil, fmt.Errorf("predicate %q not supported", condition.Predicate)
				}
				conditions[j] = cond
			}
			res[i] = process.Filter{Key: n.Filter.Key, Filter: filter.Filter{Conditions: conditions}}
		default:
			return nil, fmt.Errorf("node has unexpected type %T", n)
		}
	}
	return res, nil
}

func fromProtoProcessEdges(edges []*types.Process_Edge) []process.Edge {
	res := make([]process.Edge, len(edges))
	for i, edge := range edges {
		res[i] = process.Edge{
			Src: edge.Src,
			Dst: edge.Dst,
		}
	}
	return res
}

func fromProtoProcess(wf *types.Process) (*process.Process, error) {
	nodes, err := fromProtoProcessNodes(wf.Nodes)
	if err != nil {
		return nil, err
	}
	return &process.Process{
		Key: wf.Key,
		Graph: process.Graph{
			Nodes: nodes,
			Edges: fromProtoProcessEdges(wf.Edges),
		},
	}, nil
}

func toProtoProcessNodes(nodes []process.Node) []*types.Process_Node {
	res := make([]*types.Process_Node, len(nodes))
	for i, node := range nodes {
		protoNode := types.Process_Node{}
		switch n := node.(type) {
		case *process.Result:
			protoNode.Type = &types.Process_Node_Result_{
				Result: &types.Process_Node_Result{
					Key:          n.Key,
					InstanceHash: n.InstanceHash,
					TaskKey:      n.TaskKey,
				},
			}
		case *process.Event:
			protoNode.Type = &types.Process_Node_Event_{
				Event: &types.Process_Node_Event{
					Key:          n.Key,
					InstanceHash: n.InstanceHash,
					EventKey:     n.EventKey,
				},
			}
		case *process.Task:
			protoNode.Type = &types.Process_Node_Task_{
				Task: &types.Process_Node_Task{
					Key:          n.Key,
					InstanceHash: n.InstanceHash,
					TaskKey:      n.TaskKey,
				},
			}
		case *process.Map:
			outputs := make([]*types.Process_Node_Map_Output, len(n.Outputs))
			for j, output := range n.Outputs {
				out := &types.Process_Node_Map_Output{Key: output.Key}
				if output.Ref != nil {
					out.Value = &types.Process_Node_Map_Output_Ref{
						Ref: &types.Process_Node_Map_Output_Reference{
							NodeKey: output.Ref.NodeKey,
							Key:     output.Ref.Key,
						},
					}
				}
				outputs[j] = out
			}

			protoNode.Type = &types.Process_Node_Map_{
				Map: &types.Process_Node_Map{
					Key:     n.Key,
					Outputs: outputs,
				},
			}
		case *process.Filter:
			conditions := make([]*types.Process_Node_Filter_Condition, len(n.Conditions))
			for j, condition := range n.Conditions {
				cond := &types.Process_Node_Filter_Condition{Key: condition.Key, Value: condition.Value}
				if condition.Predicate == filter.EQ {
					cond.Predicate = types.Process_Node_Filter_Condition_EQ
				}
				conditions[j] = cond
			}

			protoNode.Type = &types.Process_Node_Filter_{
				Filter: &types.Process_Node_Filter{
					Key:        n.Key,
					Conditions: conditions,
				},
			}
		}
		res[i] = &protoNode
	}
	return res
}

func toProtoProcessEdges(edges []process.Edge) []*types.Process_Edge {
	res := make([]*types.Process_Edge, len(edges))
	for i, edge := range edges {
		res[i] = &types.Process_Edge{
			Src: edge.Src,
			Dst: edge.Dst,
		}
	}
	return res
}

func toProtoProcess(wf *process.Process) *types.Process {
	return &types.Process{
		Hash:  wf.Hash,
		Key:   wf.Key,
		Nodes: toProtoProcessNodes(wf.Nodes),
		Edges: toProtoProcessEdges(wf.Edges),
	}
}

func toProtoProcesses(processes []*process.Process) []*types.Process {
	wfs := make([]*types.Process, len(processes))
	for i, wf := range processes {
		wfs[i] = toProtoProcess(wf)
	}
	return wfs
}
