package process

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/assert"
)

func TestValidateProcess(t *testing.T) {
	trigger := &Process_Node{
		Type: &Process_Node_Result_{
			Result: &Process_Node_Result{
				Key:          "trigger:result",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			},
		},
	}

	nodes := []*Process_Node{
		trigger,
		{
			Type: &Process_Node_Task_{&Process_Node_Task{
				Key:          "nodeKey1",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}},
		},
		{
			Type: &Process_Node_Task_{&Process_Node_Task{
				Key:          "nodeKey2",
				InstanceHash: hash.Int(3),
				TaskKey:      "-",
			}},
		},
	}

	edges := []*Process_Edge{
		{Src: trigger.ID(), Dst: "nodeKey1"},
	}

	var tests = []struct {
		w     *Process
		valid bool
		err   string
	}{
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "invalid-struct",
		}, err: "should contain exactly one trigger"},
		{w: &Process{
			Nodes: []*Process_Node{
				{
					Type: &Process_Node_Result_{&Process_Node_Result{InstanceHash: hash.Int(1)}},
				},
			},
			Hash: hash.Int(1),
			Key:  "missing-key",
		}, err: "Error:Field validation for 'TaskKey' failed on the 'required' tag"},
		{w: &Process{
			Hash:  hash.Int(1),
			Key:   "edge-src-missing-node",
			Nodes: nodes,
			Edges: append(edges, &Process_Edge{Src: "-", Dst: "nodeKey2"}),
		}, err: "node \"-\" not found"},
		{w: &Process{
			Hash:  hash.Int(1),
			Key:   "edge-dst-missing-node",
			Nodes: nodes,
			Edges: append(edges, &Process_Edge{Src: "nodeKey1", Dst: "-"}),
		}, err: "node \"-\" not found"},
		{w: &Process{
			Hash:  hash.Int(1),
			Key:   "cyclic-graph",
			Nodes: nodes,
			Edges: append(edges,
				&Process_Edge{Src: "nodeKey1", Dst: "nodeKey2"},
				&Process_Edge{Src: "nodeKey2", Dst: "nodeKey1"},
			),
		}, err: "process should not contain any cycles"},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "non-connected-graph",
			Nodes: append(nodes, &Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{
					Key:          "nodeKey3",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{
					Key:          "nodeKey4",
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
			Key:  "multiple-parent-graph",
			Nodes: append(nodes, &Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{
					Key:          "nodeKey3",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{
					Key:          "nodeKey4",
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
			Key:  "multiple-parent-graph",
			Nodes: append(nodes, &Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{
					Key:          "nodeKey3",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{
					Key:          "nodeKey4",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{
					Key:          "nodeKey5",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{
					Key:          "nodeKey6",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}},
			}, &Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{
					Key:          "nodeKey7",
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
			Key:  "inputs-with-invalid-node",
			Nodes: append(nodes, &Process_Node{
				Type: &Process_Node_Map_{&Process_Node_Map{
					Key: "mapping",
					Outputs: []*Process_Node_Map_Output{
						{
							Value: &Process_Node_Map_Output_Ref{
								Ref: &Process_Node_Map_Output_Reference{OutputIndex: 0, NodeKey: "invalid"},
							},
						},
					},
				}},
			}),
		}, err: "node \"invalid\" not found"},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "inputs-with-valid-ref",
			Nodes: append(nodes, &Process_Node{
				Type: &Process_Node_Map_{&Process_Node_Map{
					Key: "mapping",
					Outputs: []*Process_Node_Map_Output{
						{
							Value: &Process_Node_Map_Output_Ref{
								Ref: &Process_Node_Map_Output_Reference{OutputIndex: 0, NodeKey: "nodeKey1"},
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
	}
	for _, test := range tests {
		err := test.w.Validate()
		if test.valid {
			assert.Nil(t, err, test.w.Key)
		} else {
			assert.Contains(t, test.w.Validate().Error(), test.err, test.w.Key)
		}
	}
}
