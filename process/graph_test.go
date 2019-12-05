package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func defaultProcess() *Process {
	return &Process{
		Nodes: []*Process_Node{
			{Key: "nodeKey1", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			{Key: "nodeKey2", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			{Key: "nodeKey3", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			{Key: "nodeKey4", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			{Key: "nodeKey5", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			{Key: "nodeKey6", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			{Key: "nodeKey7", Type: &Process_Node_Task_{&Process_Node_Task{}}}},
		Edges: []*Process_Edge{
			{Src: "nodeKey1", Dst: "nodeKey2"},
			{Src: "nodeKey2", Dst: "nodeKey3"},
			{Src: "nodeKey2", Dst: "nodeKey4"},
			{Src: "nodeKey3", Dst: "nodeKey5"},
			{Src: "nodeKey4", Dst: "nodeKey6"},
			{Src: "nodeKey4", Dst: "nodeKey7"},
		},
	}
}

func TestChildrenKeys(t *testing.T) {
	var tests = []struct {
		graph    *Process
		node     string
		children []string
	}{
		{graph: defaultProcess(), node: "nodeKey1", children: []string{"nodeKey2"}},
		{graph: defaultProcess(), node: "nodeKey2", children: []string{"nodeKey3", "nodeKey4"}},
		{graph: defaultProcess(), node: "nodeKey3", children: []string{"nodeKey5"}},
		{graph: defaultProcess(), node: "nodeKey4", children: []string{"nodeKey6", "nodeKey7"}},
		{graph: defaultProcess(), node: "nodeKey5", children: []string{}},
		{graph: defaultProcess(), node: "nodeKey6", children: []string{}},
		{graph: defaultProcess(), node: "nodeKey7", children: []string{}},
	}
	for _, test := range tests {
		assert.Equal(t, test.graph.ChildrenKeys(test.node), test.children)
	}
}

func TestParentKeys(t *testing.T) {
	var tests = []struct {
		graph   *Process
		node    string
		parents []string
	}{
		{graph: defaultProcess(), node: "nodeKey1", parents: []string{}},
		{graph: defaultProcess(), node: "nodeKey2", parents: []string{"nodeKey1"}},
		{graph: defaultProcess(), node: "nodeKey3", parents: []string{"nodeKey2"}},
		{graph: defaultProcess(), node: "nodeKey4", parents: []string{"nodeKey2"}},
		{graph: defaultProcess(), node: "nodeKey5", parents: []string{"nodeKey3"}},
		{graph: defaultProcess(), node: "nodeKey6", parents: []string{"nodeKey4"}},
		{graph: defaultProcess(), node: "nodeKey7", parents: []string{"nodeKey4"}},
	}
	for _, test := range tests {
		assert.Equal(t, test.graph.ParentKeys(test.node), test.parents)
	}
}

func TestFindNode(t *testing.T) {
	var tests = []struct {
		graph   *Process
		node    string
		present bool
	}{
		{graph: defaultProcess(), node: "nodeKey1", present: true},
		{graph: defaultProcess(), node: "nodeKey2", present: true},
		{graph: defaultProcess(), node: "nodeKey3", present: true},
		{graph: defaultProcess(), node: "nodeKey4", present: true},
		{graph: defaultProcess(), node: "nodeKey5", present: true},
		{graph: defaultProcess(), node: "nodeKey6", present: true},
		{graph: defaultProcess(), node: "nodeKey7", present: true},
		{graph: defaultProcess(), node: "nodeKey8", present: false},
	}
	for _, test := range tests {
		node, err := test.graph.FindNode(test.node)
		if test.present {
			assert.NoError(t, err)
			assert.Equal(t, node.Key, test.node)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestHasNodes(t *testing.T) {
	var tests = []struct {
		graph    *Process
		hasNodes bool
	}{
		{graph: defaultProcess(), hasNodes: true},
		{graph: &Process{}, hasNodes: false},
	}
	for _, test := range tests {
		assert.Equal(t, test.graph.hasNodes(), test.hasNodes)
	}
}

func TestIsAcyclic(t *testing.T) {
	var tests = []struct {
		graph   *Process
		acyclic bool
	}{
		{graph: defaultProcess(), acyclic: true},
		{graph: &Process{
			Nodes: []*Process_Node{
				{Key: "nodeKey1", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey2", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			},
			Edges: []*Process_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey1"},
			},
		}, acyclic: false},
		{graph: &Process{
			Nodes: []*Process_Node{
				{Key: "nodeKey1", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey2", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey3", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			},
			Edges: []*Process_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey3"},
				{Src: "nodeKey3", Dst: "nodeKey1"},
			},
		}, acyclic: false},
		{graph: &Process{
			Nodes: []*Process_Node{
				{Key: "nodeKey1", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey2", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey3", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			},
			Edges: []*Process_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey3"},
				{Src: "nodeKey3", Dst: "nodeKey2"},
			},
		}, acyclic: false},
		{graph: &Process{
			Nodes: []*Process_Node{
				{Key: "nodeKey1", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey2", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey3", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey4", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			},
			Edges: []*Process_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey1", Dst: "nodeKey3"},
				{Src: "nodeKey2", Dst: "nodeKey4"},
				{Src: "nodeKey3", Dst: "nodeKey4"},
			},
		}, acyclic: true},
	}
	for _, test := range tests {
		assert.Equal(t, test.graph.isAcyclic(), test.acyclic)
	}
}

func TestIsConnected(t *testing.T) {
	var tests = []struct {
		graph     *Process
		node      string
		connected bool
	}{
		{graph: defaultProcess(), connected: true},
		{graph: &Process{
			Nodes: []*Process_Node{
				{Key: "nodeKey1", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey2", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey3", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey4", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			},
			Edges: []*Process_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey3", Dst: "nodeKey4"},
			},
		}, connected: false},
	}
	for i, test := range tests {
		assert.Equal(t, test.graph.isConnected(), test.connected, i)
	}
}

func TestVisitChildren(t *testing.T) {
	var tests = []struct {
		graph    *Process
		node     string
		children []string
	}{
		{graph: defaultProcess(), node: "nodeKey1", children: []string{"nodeKey2", "nodeKey3", "nodeKey4", "nodeKey5", "nodeKey6", "nodeKey7"}},
		{graph: defaultProcess(), node: "nodeKey2", children: []string{"nodeKey3", "nodeKey4", "nodeKey5", "nodeKey6", "nodeKey7"}},
		{graph: defaultProcess(), node: "nodeKey3", children: []string{"nodeKey5"}},
		{graph: defaultProcess(), node: "nodeKey4", children: []string{"nodeKey6", "nodeKey7"}},
		{graph: defaultProcess(), node: "nodeKey5", children: []string{}},
		{graph: defaultProcess(), node: "nodeKey6", children: []string{}},
		{graph: defaultProcess(), node: "nodeKe7", children: []string{}},
	}
	for _, test := range tests {
		visit := make(map[string]bool)
		test.graph.dfs(test.node, func(node string) {
			visit[node] = true
		})
		for _, child := range test.children {
			assert.True(t, visit[child])
		}
	}
}

func TestGetRoot(t *testing.T) {
	var tests = []struct {
		graph *Process
		node  string
		root  string
	}{
		{graph: defaultProcess(), node: "nodeKey1", root: "nodeKey1"},
		{graph: defaultProcess(), node: "nodeKey5", root: "nodeKey1"},
		{graph: defaultProcess(), node: "nodeKey6", root: "nodeKey1"},
		{graph: defaultProcess(), node: "nodeKey4", root: "nodeKey1"},
	}
	for _, test := range tests {
		assert.Equal(t, test.graph.getRoot(test.node), test.root)
	}
}

func TestIsMonoParental(t *testing.T) {
	var tests = []struct {
		graph *Process
		max   int
	}{
		{graph: defaultProcess(), max: 1},
		{graph: &Process{
			Nodes: []*Process_Node{
				{Key: "nodeKey1", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey2", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey3", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			},
			Edges: []*Process_Edge{
				{Src: "nodeKey1", Dst: "nodeKey3"},
				{Src: "nodeKey2", Dst: "nodeKey3"},
			},
		}, max: 2},
		{graph: &Process{
			Nodes: []*Process_Node{
				{Key: "nodeKey1", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey2", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey3", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey4", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			},
			Edges: []*Process_Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey1", Dst: "nodeKey3"},
				{Src: "nodeKey2", Dst: "nodeKey4"},
				{Src: "nodeKey3", Dst: "nodeKey4"},
			},
		}, max: 2},
		{graph: &Process{
			Nodes: []*Process_Node{
				{Key: "nodeKey1", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey2", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey3", Type: &Process_Node_Task_{&Process_Node_Task{}}},
				{Key: "nodeKey4", Type: &Process_Node_Task_{&Process_Node_Task{}}},
			},
			Edges: []*Process_Edge{
				{Src: "nodeKey1", Dst: "nodeKey4"},
				{Src: "nodeKey2", Dst: "nodeKey4"},
				{Src: "nodeKey3", Dst: "nodeKey4"},
			},
		}, max: 3},
	}
	for _, test := range tests {
		assert.Equal(t, test.max, test.graph.maximumParents())
	}
}
