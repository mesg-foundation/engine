package workflow

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	var tests = []struct {
		filters []*Workflow_Trigger_Filter
		data    map[string]interface{}
		match   bool
	}{
		{ // not matching filter
			filters: []*Workflow_Trigger_Filter{
				{Key: "foo", Predicate: Workflow_Trigger_Filter_EQ, Value: "xx"},
			},
			data:  map[string]interface{}{"foo": "bar"},
			match: false,
		},
		{ // matching multiple filters
			filters: []*Workflow_Trigger_Filter{
				{Key: "foo", Predicate: Workflow_Trigger_Filter_EQ, Value: "bar"},
				{Key: "xxx", Predicate: Workflow_Trigger_Filter_EQ, Value: "yyy"},
			},
			data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			},
			match: true,
		},
		{ // non matching multiple filters
			filters: []*Workflow_Trigger_Filter{
				{Key: "foo", Predicate: Workflow_Trigger_Filter_EQ, Value: "bar"},
				{Key: "xxx", Predicate: Workflow_Trigger_Filter_EQ, Value: "aaa"},
			},
			data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			},
			match: false,
		},
	}
	for i, test := range tests {
		match := true
		for _, f := range test.filters {
			if !f.Match(test.data) {
				match = false
				break
			}
		}
		assert.Equal(t, test.match, match, i)
	}
}

func TestValidateWorkflow(t *testing.T) {
	trigger := Workflow_Trigger{
		InstanceHash: hash.Int(2),
		Key: &Workflow_Trigger_TaskKey{
			TaskKey: "-",
		},
		NodeKey: "nodeKey1",
	}

	nodes := []*Workflow_Node{
		{
			Key:          "nodeKey1",
			InstanceHash: hash.Int(2),
			TaskKey:      "-",
		},
		{
			Key:          "nodeKey2",
			InstanceHash: hash.Int(3),
			TaskKey:      "-",
		},
	}

	var tests = []struct {
		w     *Workflow
		valid bool
		err   string
	}{
		{w: &Workflow{
			Hash: hash.Int(1),
			Key:  "invalid-struct",
		}, err: "Error:Field validation"},
		{w: &Workflow{
			Trigger: Workflow_Trigger{InstanceHash: hash.Int(1), NodeKey: "-"},
			Hash:    hash.Int(1),
			Key:     "missing-key",
		}, err: "Key: 'Workflow.Trigger.TaskKey' Error:Field validation for 'TaskKey' failed on the 'required_without' tag"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "trigger-missing-node",
			Trigger: trigger,
		}, err: "node \"nodeKey1\" not found"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "edge-src-missing-node",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []*Workflow_Edge{
				{Src: "-", Dst: "nodeKey2"},
			},
		}, err: "node \"-\" not found"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "edge-dst-missing-node",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []*Workflow_Edge{
				{Src: "nodeKey1", Dst: "-"},
			},
		}, err: "node \"-\" not found"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "cyclic-graph",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []*Workflow_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey1"},
			},
		}, err: "workflow should not contain any cycles"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "non-connected-graph",
			Trigger: trigger,
			Nodes: append(nodes, &Workflow_Node{
				Key:          "nodeKey3",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, &Workflow_Node{
				Key:          "nodeKey4",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}),
			Edges: []*Workflow_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey3", Dst: "nodeKey4"},
			},
		}, err: "workflow should be a connected graph"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "multiple-parent-graph",
			Trigger: trigger,
			Nodes: append(nodes, &Workflow_Node{
				Key:          "nodeKey3",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, &Workflow_Node{
				Key:          "nodeKey4",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}),
			Edges: []*Workflow_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey1", Dst: "nodeKey3"},
				{Src: "nodeKey2", Dst: "nodeKey4"},
				{Src: "nodeKey3", Dst: "nodeKey4"},
			},
		}, err: "workflow should contain nodes with one parent maximum"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "multiple-parent-graph",
			Trigger: trigger,
			Nodes: append(nodes, &Workflow_Node{
				Key:          "nodeKey3",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, &Workflow_Node{
				Key:          "nodeKey4",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, &Workflow_Node{
				Key:          "nodeKey5",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, &Workflow_Node{
				Key:          "nodeKey6",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, &Workflow_Node{
				Key:          "nodeKey7",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}),
			Edges: []*Workflow_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey3"},
				{Src: "nodeKey2", Dst: "nodeKey4"},
				{Src: "nodeKey3", Dst: "nodeKey5"},
				{Src: "nodeKey4", Dst: "nodeKey6"},
				{Src: "nodeKey4", Dst: "nodeKey7"},
			},
		}, valid: true},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "inputs-with-invalid-node",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []*Workflow_Edge{
				{
					Src: "nodeKey1",
					Dst: "nodeKey2",
					Inputs: []*Workflow_Edge_Input{
						{
							Key: "-",
							Value: &Workflow_Edge_Input_Ref{
								Ref: &Workflow_Edge_Input_Reference{
									Key:     "-",
									NodeKey: "invalid"},
							},
						},
					}},
			},
		}, err: "node \"invalid\" not found"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "inputs-with-valid-ref",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []*Workflow_Edge{
				{
					Src: "nodeKey1",
					Dst: "nodeKey2",
					Inputs: []*Workflow_Edge_Input{
						{
							Key: "-",
							Value: &Workflow_Edge_Input_Ref{
								Ref: &Workflow_Edge_Input_Reference{
									Key:     "-",
									NodeKey: "nodeKey1"},
							},
						},
					}},
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
