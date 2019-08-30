package process

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/assert"
)

func TestValidateProcess(t *testing.T) {

	trigger := Result{
		Key:          "trigger:result",
		InstanceHash: hash.Int(2),
		TaskKey:      "-",
	}

	nodes := []Node{
		trigger,
		Task{
			Key:          "nodeKey1",
			InstanceHash: hash.Int(2),
			TaskKey:      "-",
		},
		Task{
			Key:          "nodeKey2",
			InstanceHash: hash.Int(3),
			TaskKey:      "-",
		},
	}

	edges := []Edge{
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
			Graph: Graph{
				Nodes: []Node{Result{InstanceHash: hash.Int(1)}},
			},
			Hash: hash.Int(1),
			Key:  "missing-key",
		}, err: "Error:Field validation for 'TaskKey' failed on the 'required' tag"},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "edge-src-missing-node",
			Graph: Graph{
				Nodes: nodes,
				Edges: append(edges,
					Edge{Src: "-", Dst: "nodeKey2"},
				),
			},
		}, err: "node \"-\" not found"},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "edge-dst-missing-node",
			Graph: Graph{
				Nodes: nodes,
				Edges: append(edges,
					Edge{Src: "nodeKey1", Dst: "-"},
				),
			},
		}, err: "node \"-\" not found"},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "cyclic-graph",
			Graph: Graph{
				Nodes: nodes,
				Edges: append(edges,
					Edge{Src: "nodeKey1", Dst: "nodeKey2"},
					Edge{Src: "nodeKey2", Dst: "nodeKey1"},
				),
			},
		}, err: "process should not contain any cycles"},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "non-connected-graph",
			Graph: Graph{
				Nodes: append(nodes, Task{
					Key:          "nodeKey3",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}, Task{
					Key:          "nodeKey4",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}),
				Edges: append(edges,
					Edge{Src: "nodeKey1", Dst: "nodeKey2"},
					Edge{Src: "nodeKey3", Dst: "nodeKey4"},
				),
			},
		}, err: "process should be a connected graph"},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "multiple-parent-graph",
			Graph: Graph{
				Nodes: append(nodes, Task{
					Key:          "nodeKey3",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}, Task{
					Key:          "nodeKey4",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}),
				Edges: append(edges,
					Edge{Src: "nodeKey1", Dst: "nodeKey2"},
					Edge{Src: "nodeKey1", Dst: "nodeKey3"},
					Edge{Src: "nodeKey2", Dst: "nodeKey4"},
					Edge{Src: "nodeKey3", Dst: "nodeKey4"},
				),
			},
		}, err: "process should contain nodes with one parent maximum"},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "multiple-parent-graph",
			Graph: Graph{
				Nodes: append(nodes, Task{
					Key:          "nodeKey3",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}, Task{
					Key:          "nodeKey4",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}, Task{
					Key:          "nodeKey5",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}, Task{
					Key:          "nodeKey6",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}, Task{
					Key:          "nodeKey7",
					InstanceHash: hash.Int(2),
					TaskKey:      "-",
				}),
				Edges: append(edges,
					Edge{Src: "nodeKey1", Dst: "nodeKey2"},
					Edge{Src: "nodeKey2", Dst: "nodeKey3"},
					Edge{Src: "nodeKey2", Dst: "nodeKey4"},
					Edge{Src: "nodeKey3", Dst: "nodeKey5"},
					Edge{Src: "nodeKey4", Dst: "nodeKey6"},
					Edge{Src: "nodeKey4", Dst: "nodeKey7"},
				),
			},
		}, valid: true},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "inputs-with-invalid-node",
			Graph: Graph{
				Nodes: append(nodes, Map{
					Key: "mapping",
					Outputs: []Output{
						{Key: "-", Ref: &OutputReference{Key: "-", NodeKey: "invalid"}},
					},
				}),
			},
		}, err: "node \"invalid\" not found"},
		{w: &Process{
			Hash: hash.Int(1),
			Key:  "inputs-with-valid-ref",
			Graph: Graph{
				Nodes: append(nodes, Map{
					Key: "mapping",
					Outputs: []Output{
						{Key: "-", Ref: &OutputReference{Key: "-", NodeKey: "nodeKey1"}},
					},
				}),
				Edges: append(edges,
					Edge{Src: "nodeKey1", Dst: "mapping"},
					Edge{Src: "mapping", Dst: "nodeKey2"},
				),
			},
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
