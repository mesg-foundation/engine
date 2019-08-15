package workflow

import "fmt"

// ChildrenIDs returns the list of node IDs with a dependency to the current node
func (w Workflow) ChildrenIDs(nodeKey string) []string {
	nodeKeys := make([]string, 0)
	for _, edge := range w.Edges {
		if edge.Src == nodeKey {
			nodeKeys = append(nodeKeys, edge.Dst)
		}
	}
	return nodeKeys
}

// ParentIDs returns the list of node IDs with the current node as child
func (w Workflow) ParentIDs(nodeKey string) []string {
	nodeKeys := make([]string, 0)
	for _, edge := range w.Edges {
		if edge.Dst == nodeKey {
			nodeKeys = append(nodeKeys, edge.Src)
		}
	}
	return nodeKeys
}

// FindNode returns the node matching the key in parameter or an error if not found
func (w Workflow) FindNode(key string) (Node, error) {
	for _, node := range w.Nodes {
		if node.Key == key {
			return node, nil
		}
	}
	return Node{}, fmt.Errorf("node %q not found", key)
}

// A null graph is a graph that contains no nodes
func (w *Workflow) hasNodes() bool {
	return len(w.Nodes) > 0
}

// An acyclic graph is a graph that doesn't contain a cycle. If you walk through the graph you will go maximum one time on each node.
func (w *Workflow) isAcyclic() bool {
	visited := make(map[string]bool)
	recursive := make(map[string]bool)
	for _, node := range w.Nodes {
		if w.hasCycle(node.Key, visited, recursive) {
			return false
		}
	}
	return true
}

// Check if the descendant of a node are creating any cycle. https://algorithms.tutorialhorizon.com/graph-detect-cycle-in-a-directed-graph/
func (w *Workflow) hasCycle(node string, visited map[string]bool, recursive map[string]bool) bool {
	visited[node] = true
	recursive[node] = true
	for _, child := range w.ChildrenIDs(node) {
		if !visited[child] && w.hasCycle(child, visited, recursive) {
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
func (w *Workflow) isConnected() bool {
	root := w.getRoot(w.Nodes[0].Key)
	visited := make(map[string]bool)
	w.dfs(root, func(node string) {
		visited[node] = true
	})
	return len(visited) == len(w.Nodes)
}

// walk through all the children of a node and populate a map of visited children.
func (w *Workflow) dfs(node string, fn func(node string)) {
	fn(node)
	for _, n := range w.ChildrenIDs(node) {
		w.dfs(n, fn)
	}
}

func (w *Workflow) getRoot(node string) string {
	parents := w.ParentIDs(node)
	if len(parents) == 0 {
		return node
	}
	if len(parents) > 1 {
		panic("multiple parents is not supported")
	}
	return w.getRoot(parents[0])
}

// Return the maximum number of parent found in the graph.
func (w *Workflow) maximumParents() int {
	max := 0
	for _, node := range w.Nodes {
		if max < len(w.ParentIDs(node.Key)) {
			max = len(w.ParentIDs(node.Key))
		}
	}
	return max
}

func (w *Workflow) shouldBeDirectedTree() error {
	if !w.hasNodes() {
		return fmt.Errorf("workflow needs to have at least one node")
	}
	if !w.isAcyclic() {
		return fmt.Errorf("workflow should not contain any cycles")
	}
	if !w.isConnected() {
		return fmt.Errorf("workflow should be a connected graph")
	}
	if w.maximumParents() > 1 {
		return fmt.Errorf("workflow should contain nodes with one parent maximum")
	}
	return nil
}
