package process

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateProcess(t *testing.T) {
	trigger := &Process_Node{
		Key: "trigger:result",
		Type: &Process_Node_Result_{
			Result: &Process_Node_Result{
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			},
		},
	}

	nodes := []*Process_Node{
		trigger,
		{
			Key: "nodeKey1",
			Type: &Process_Node_Task_{&Process_Node_Task{
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}},
		},
		{
			Key: "nodeKey2",
			Type: &Process_Node_Task_{&Process_Node_Task{
				InstanceHash: hash.Int(3),
				TaskKey:      "-",
			}},
		},
	}

	edges := []*Process_Edge{
		{Src: trigger.Key, Dst: "nodeKey1"},
	}

	var tests = []struct {
		w     *Process
		valid bool
		err   string
	}{
		{w: &Process{
			Hash: hash.Int(1),
			Name: "invalid-struct",
		}, err: "should contain exactly one trigger"},
		{w: &Process{
			Nodes: []*Process_Node{
				{
					Type: &Process_Node_Result_{&Process_Node_Result{InstanceHash: hash.Int(1)}},
				},
			},
			Hash: hash.Int(1),
			Name: "missing-key",
		}, err: "Key is a required field. TaskKey is a required field"},
		{w: &Process{
			Hash:  hash.Int(1),
			Name:  "edge-src-missing-node",
			Nodes: nodes,
			Edges: append(edges, &Process_Edge{Src: "-", Dst: "nodeKey2"}),
		}, err: "node \"-\" not found"},
		{w: &Process{
			Hash:  hash.Int(1),
			Name:  "edge-dst-missing-node",
			Nodes: nodes,
			Edges: append(edges, &Process_Edge{Src: "nodeKey1", Dst: "-"}),
		}, err: "node \"-\" not found"},
		{w: &Process{
			Hash:  hash.Int(1),
			Name:  "cyclic-graph",
			Nodes: nodes,
			Edges: append(edges,
				&Process_Edge{Src: "nodeKey1", Dst: "nodeKey2"},
				&Process_Edge{Src: "nodeKey2", Dst: "nodeKey1"},
			),
		}, err: "process should not contain any cycles"},
		{w: &Process{
			Hash: hash.Int(1),
			Name: "non-connected-graph",
			Nodes: append(nodes, &Process_Node{
				Key: "nodeKey3",
				Type: &Process_Node_Task_{&Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Key: "nodeKey4",
				Type: &Process_Node_Task_{&Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}),
			Edges: append(edges,
				&Process_Edge{Src: "nodeKey1", Dst: "nodeKey2"},
				&Process_Edge{Src: "nodeKey3", Dst: "nodeKey4"},
			),
		}, err: "process should be a connected graph"},
		{w: &Process{
			Hash: hash.Int(1),
			Name: "multiple-parent-graph",
			Nodes: append(nodes, &Process_Node{
				Key: "nodeKey3",
				Type: &Process_Node_Task_{&Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Key: "nodeKey4",
				Type: &Process_Node_Task_{&Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}),
			Edges: append(edges,
				&Process_Edge{Src: "nodeKey1", Dst: "nodeKey2"},
				&Process_Edge{Src: "nodeKey1", Dst: "nodeKey3"},
				&Process_Edge{Src: "nodeKey2", Dst: "nodeKey4"},
				&Process_Edge{Src: "nodeKey3", Dst: "nodeKey4"},
			),
		}, err: "process should contain nodes with one parent maximum"},
		{w: &Process{
			Hash: hash.Int(1),
			Name: "multiple-parent-graph",
			Nodes: append(nodes, &Process_Node{
				Key: "nodeKey3",
				Type: &Process_Node_Task_{&Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Key: "nodeKey4",
				Type: &Process_Node_Task_{&Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Key: "nodeKey5",
				Type: &Process_Node_Task_{&Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Key: "nodeKey6",
				Type: &Process_Node_Task_{&Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Key: "nodeKey7",
				Type: &Process_Node_Task_{&Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}),
			Edges: append(edges,
				&Process_Edge{Src: "nodeKey1", Dst: "nodeKey2"},
				&Process_Edge{Src: "nodeKey2", Dst: "nodeKey3"},
				&Process_Edge{Src: "nodeKey2", Dst: "nodeKey4"},
				&Process_Edge{Src: "nodeKey3", Dst: "nodeKey5"},
				&Process_Edge{Src: "nodeKey4", Dst: "nodeKey6"},
				&Process_Edge{Src: "nodeKey4", Dst: "nodeKey7"},
			),
		}, valid: true},
		{w: &Process{
			Hash: hash.Int(1),
			Name: "inputs-with-invalid-node",
			Nodes: append(nodes, &Process_Node{
				Key: "mapping",
				Type: &Process_Node_Map_{&Process_Node_Map{
					Outputs: map[string]*Process_Node_Map_Output{
						"key": {
							Value: &Process_Node_Map_Output_Ref{
								Ref: &Process_Node_Map_Output_Reference{NodeKey: "invalid"},
							},
						},
					},
				}},
			}),
		}, err: "node \"invalid\" not found"},
		{w: &Process{
			Hash: hash.Int(1),
			Name: "inputs-with-valid-ref",
			Nodes: append(nodes, &Process_Node{
				Key: "mapping",
				Type: &Process_Node_Map_{&Process_Node_Map{
					Outputs: map[string]*Process_Node_Map_Output{
						"key": {
							Value: &Process_Node_Map_Output_Ref{
								Ref: &Process_Node_Map_Output_Reference{NodeKey: "nodeKey1"},
							},
						},
					},
				}},
			}),
			Edges: append(edges,
				&Process_Edge{Src: "nodeKey1", Dst: "mapping"},
				&Process_Edge{Src: "mapping", Dst: "nodeKey2"},
			),
		}, valid: true},
		{w: &Process{
			Hash: hash.Int(1),
			Name: "filter-gte-invalid-value-type",
			Nodes: append(nodes, &Process_Node{
				Key: "filter",
				Type: &Process_Node_Filter_{&Process_Node_Filter{
					Conditions: []Process_Node_Filter_Condition{
						{
							Key:       "foo",
							Predicate: Process_Node_Filter_Condition_GT,
							Value: &types.Value{
								Kind: &types.Value_StringValue{StringValue: "bar"},
							},
						},
					},
				}},
			}),
		}, err: "filter with condition GT, GTE, LT or LTE only works with value of type Number"},
	}
	for _, test := range tests {
		err := test.w.Validate()
		if test.valid {
			assert.Nil(t, err, test.w.Name)
		} else {
			assert.Contains(t, test.w.Validate().Error(), test.err, test.w.Name)
		}
	}
}
