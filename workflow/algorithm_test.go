package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func defaultGraph() *Workflow {
	return &Workflow{
		Nodes: []Node{
			{Key: "nodeKey1"},
			{Key: "nodeKey2"},
			{Key: "nodeKey3"},
			{Key: "nodeKey4"},
			{Key: "nodeKey5"},
			{Key: "nodeKey6"},
			{Key: "nodeKey7"}},
		Edges: []Edge{
			{Src: "nodeKey1", Dst: "nodeKey2"},
			{Src: "nodeKey2", Dst: "nodeKey3"},
			{Src: "nodeKey2", Dst: "nodeKey4"},
			{Src: "nodeKey3", Dst: "nodeKey5"},
			{Src: "nodeKey4", Dst: "nodeKey6"},
			{Src: "nodeKey4", Dst: "nodeKey7"},
		},
	}
}

func TestChildrenIDs(t *testing.T) {
	var tests = []struct {
		graph    *Workflow
		node     string
		children []string
	}{
		{graph: defaultGraph(), node: "nodeKey1", children: []string{"nodeKey2"}},
		{graph: defaultGraph(), node: "nodeKey2", children: []string{"nodeKey3", "nodeKey4"}},
		{graph: defaultGraph(), node: "nodeKey3", children: []string{"nodeKey5"}},
		{graph: defaultGraph(), node: "nodeKey4", children: []string{"nodeKey6", "nodeKey7"}},
		{graph: defaultGraph(), node: "nodeKey5", children: []string{}},
		{graph: defaultGraph(), node: "nodeKey6", children: []string{}},
		{graph: defaultGraph(), node: "nodeKey7", children: []string{}},
	}
	for _, test := range tests {
		assert.Equal(t, test.graph.ChildrenIDs(test.node), test.children)
	}
}

func TestParentIDs(t *testing.T) {
	var tests = []struct {
		graph   *Workflow
		node    string
		parents []string
	}{
		{graph: defaultGraph(), node: "nodeKey1", parents: []string{}},
		{graph: defaultGraph(), node: "nodeKey2", parents: []string{"nodeKey1"}},
		{graph: defaultGraph(), node: "nodeKey3", parents: []string{"nodeKey2"}},
		{graph: defaultGraph(), node: "nodeKey4", parents: []string{"nodeKey2"}},
		{graph: defaultGraph(), node: "nodeKey5", parents: []string{"nodeKey3"}},
		{graph: defaultGraph(), node: "nodeKey6", parents: []string{"nodeKey4"}},
		{graph: defaultGraph(), node: "nodeKey7", parents: []string{"nodeKey4"}},
	}
	for _, test := range tests {
		assert.Equal(t, test.graph.ParentIDs(test.node), test.parents)
	}
}

func TestFindNode(t *testing.T) {
	var tests = []struct {
		graph   *Workflow
		node    string
		present bool
	}{
		{graph: defaultGraph(), node: "nodeKey1", present: true},
		{graph: defaultGraph(), node: "nodeKey2", present: true},
		{graph: defaultGraph(), node: "nodeKey3", present: true},
		{graph: defaultGraph(), node: "nodeKey4", present: true},
		{graph: defaultGraph(), node: "nodeKey5", present: true},
		{graph: defaultGraph(), node: "nodeKey6", present: true},
		{graph: defaultGraph(), node: "nodeKey7", present: true},
		{graph: defaultGraph(), node: "nodeKey8", present: false},
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
		graph *Workflow
		null  bool
	}{
		{graph: defaultGraph(), null: false},
		{graph: &Workflow{}, null: true},
	}
	for _, test := range tests {
		assert.Equal(t, test.graph.hasNodes(), test.null)
	}
}

func TestIsAcyclic(t *testing.T) {
	var tests = []struct {
		graph   *Workflow
		acyclic bool
	}{
		{graph: defaultGraph(), acyclic: true},
		{graph: &Workflow{
			Nodes: []Node{
				{Key: "nodeKey1"},
				{Key: "nodeKey2"},
			},
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey1"},
			},
		}, acyclic: false},
		{graph: &Workflow{
			Nodes: []Node{
				{Key: "nodeKey1"},
				{Key: "nodeKey2"},
				{Key: "nodeKey3"},
			},
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey3"},
				{Src: "nodeKey3", Dst: "nodeKey1"},
			},
		}, acyclic: false},
		{graph: &Workflow{
			Nodes: []Node{
				{Key: "nodeKey1"},
				{Key: "nodeKey2"},
				{Key: "nodeKey3"},
			},
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey3"},
				{Src: "nodeKey3", Dst: "nodeKey2"},
			},
		}, acyclic: false},
		{graph: &Workflow{
			Nodes: []Node{
				{Key: "nodeKey1"},
				{Key: "nodeKey2"},
				{Key: "nodeKey3"},
				{Key: "nodeKey4"},
			},
			Edges: []Edge{
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
		graph     *Workflow
		node      string
		connected bool
	}{
		{graph: defaultGraph(), connected: true},
		{graph: &Workflow{
			Nodes: []Node{
				{Key: "nodeKey1"},
				{Key: "nodeKey2"},
				{Key: "nodeKey3"},
				{Key: "nodeKey4"},
			},
			Edges: []Edge{
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
		graph    *Workflow
		node     string
		children []string
	}{
		{graph: defaultGraph(), node: "nodeKey1", children: []string{"nodeKey2", "nodeKey3", "nodeKey4", "nodeKey5", "nodeKey6", "nodeKey7"}},
		{graph: defaultGraph(), node: "nodeKey2", children: []string{"nodeKey3", "nodeKey4", "nodeKey5", "nodeKey6", "nodeKey7"}},
		{graph: defaultGraph(), node: "nodeKey3", children: []string{"nodeKey5"}},
		{graph: defaultGraph(), node: "nodeKey4", children: []string{"nodeKey6", "nodeKey7"}},
		{graph: defaultGraph(), node: "nodeKey5", children: []string{}},
		{graph: defaultGraph(), node: "nodeKey6", children: []string{}},
		{graph: defaultGraph(), node: "nodeKe7", children: []string{}},
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
		graph *Workflow
		node  string
		root  string
	}{
		{graph: defaultGraph(), node: "nodeKey1", root: "nodeKey1"},
		{graph: defaultGraph(), node: "nodeKey5", root: "nodeKey1"},
		{graph: defaultGraph(), node: "nodeKey6", root: "nodeKey1"},
		{graph: defaultGraph(), node: "nodeKey4", root: "nodeKey1"},
	}
	for _, test := range tests {
		assert.Equal(t, test.graph.getRoot(test.node), test.root)
	}
}

func TestIsMonoParental(t *testing.T) {
	var tests = []struct {
		graph      *Workflow
		monoParent bool
	}{
		{graph: defaultGraph(), monoParent: true},
		{graph: &Workflow{
			Nodes: []Node{
				{Key: "nodeKey1"},
				{Key: "nodeKey2"},
				{Key: "nodeKey3"},
			},
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey3"},
				{Src: "nodeKey2", Dst: "nodeKey3"},
			},
		}, monoParent: false},
		{graph: &Workflow{
			Nodes: []Node{
				{Key: "nodeKey1"},
				{Key: "nodeKey2"},
				{Key: "nodeKey3"},
				{Key: "nodeKey4"},
			},
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey1", Dst: "nodeKey3"},
				{Src: "nodeKey2", Dst: "nodeKey4"},
				{Src: "nodeKey3", Dst: "nodeKey4"},
			},
		}, monoParent: false},
	}
	for _, test := range tests {
		assert.Equal(t, test.monoParent, test.graph.isMonoParental())
	}
}
