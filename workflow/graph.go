package workflow

import "fmt"

// Graph is a graph structure
type Graph struct {
	Nodes []Node `hash:"name:1" validate:"dive,required"`
	Edges []Edge `hash:"name:2" validate:"dive,required"`
}

// Node type
type Node interface {
	ID() string
}

// Edge connects two nodes together based on their ID
type Edge struct {
	Src string `hash:"name:1" validate:"required"`
	Dst string `hash:"name:2" validate:"required"`
}

// ChildrenIDs returns the list of node IDs with a dependency to the current node
func (g Graph) ChildrenIDs(nodeID string) []string {
	nodeIDs := make([]string, 0)
	for _, edge := range g.EdgesFrom(nodeID) {
		nodeIDs = append(nodeIDs, edge.Dst)
	}
	return nodeIDs
}

// ParentIDs returns the list of node IDs with the current node as child
func (g Graph) ParentIDs(nodeID string) []string {
	nodeIDs := make([]string, 0)
	for _, edge := range g.Edges {
		if edge.Dst == nodeID {
			nodeIDs = append(nodeIDs, edge.Src)
		}
	}
	return nodeIDs
}

// FindNodes returns a list of nodes matching a specific filter
func (g Graph) FindNodes(filter func(n Node) bool) []Node {
	nodes := make([]Node, 0)
	for _, node := range g.Nodes {
		if filter(node) {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// FindNode return a specific node in a graph identifies by its ID. Returns an error if there is no match or multiple matches
func (g Graph) FindNode(id string) (Node, error) {
	nodes := g.FindNodes(func(n Node) bool {
		return n.ID() == id
	})
	if len(nodes) == 0 {
		return nil, fmt.Errorf("node %q not found", id)
	}
	if len(nodes) > 1 {
		return nil, fmt.Errorf("multiple nodes with the id %q", id)
	}
	return nodes[0], nil
}

// EdgesFrom return all the edges that has a common source
func (g Graph) EdgesFrom(src string) []Edge {
	edges := make([]Edge, 0)
	for _, edge := range g.Edges {
		if edge.Src == src {
			edges = append(edges, edge)
		}
	}
	return edges
}

// A null graph is a graph that contains no nodes
func (g Graph) hasNodes() bool {
	return len(g.Nodes) > 0
}

// An acyclic graph is a graph that doesn't contain a cycle. If you walk through the graph you will go maximum one time on each node.
func (g Graph) isAcyclic() bool {
	visited := make(map[string]bool)
	recursive := make(map[string]bool)
	for _, node := range g.Nodes {
		if g.hasCycle(node.ID(), visited, recursive) {
			return false
		}
	}
	return true
}

// Check if the descendant of a node are creating any cycle. https://algorithms.tutorialhorizon.com/graph-detect-cycle-in-a-directed-graph/
func (g Graph) hasCycle(node string, visited map[string]bool, recursive map[string]bool) bool {
	visited[node] = true
	recursive[node] = true
	for _, child := range g.ChildrenIDs(node) {
		if !visited[child] && g.hasCycle(child, visited, recursive) {
			return true
		}
		if recursive[child] {
			return true
		}
	}
	recursive[node] = false
	return false
}

// A connected graph is a graph where all the nodes are connected with each other through edges.
// Warning: this function will have a stack overflow if the graph is not acyclic.
func (g Graph) isConnected() bool {
	root := g.getRoot(g.Nodes[0].ID())
	visited := make(map[string]bool)
	g.dfs(root, func(node string) {
		visited[node] = true
	})
	return len(visited) == len(g.Nodes)
}

// walk through all the children of a node and populate a map of visited children.
func (g Graph) dfs(node string, fn func(node string)) {
	fn(node)
	for _, n := range g.ChildrenIDs(node) {
		g.dfs(n, fn)
	}
}

// getRoot get the root of the tree graph
func (g Graph) getRoot(node string) string {
	parents := g.ParentIDs(node)
	if len(parents) == 0 {
		return node
	}
	if len(parents) > 1 {
		panic("multiple parents is not supported")
	}
	return g.getRoot(parents[0])
}

// Return the maximum number of parent found in the graph.
func (g Graph) maximumParents() int {
	max := 0
	for _, node := range g.Nodes {
		if l := len(g.ParentIDs(node.ID())); max < l {
			max = l
		}
	}
	return max
}

func (g Graph) shouldBeDirectedTree() error {
	if !g.hasNodes() {
		return fmt.Errorf("workflow needs to have at least one node")
	}
	if !g.isAcyclic() {
		return fmt.Errorf("workflow should not contain any cycles")
	}
	if !g.isConnected() {
		return fmt.Errorf("workflow should be a connected graph")
	}
	if g.maximumParents() > 1 {
		return fmt.Errorf("workflow should contain nodes with one parent maximum")
	}
	return nil
}

func (g Graph) validate() error {
	for _, edge := range g.Edges {
		if _, err := g.FindNode(edge.Src); err != nil {
			return err
		}
		if _, err := g.FindNode(edge.Dst); err != nil {
			return err
		}
	}
	return nil
}
