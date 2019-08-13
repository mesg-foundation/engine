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
func (w *Workflow) isNullGraph() bool {
	if w.Nodes == nil {
		return true
	}
	return len(w.Nodes) == 0
}

// An acyclic graph is a graph that doesn't contain a cycle. If you walk through the graph you will go maximum one time on each node.
func (w *Workflow) isAcyclic(node string, discovered map[string]bool) bool {
	if discovered[node] {
		return false
	}
	discovered[node] = true
	for _, child := range w.ChildrenIDs(node) {
		copyMap := make(map[string]bool)
		for k, v := range discovered {
			copyMap[k] = v
		}
		if !w.isAcyclic(child, copyMap) {
			return false
		}
	}
	return true
}

// A connected graph is a graph where all the nodes are connected with each other through edges.
// Warning: this function will have a stack overflow if the graph is not acyclic.
func (w *Workflow) isConnected(node string) bool {
	visited := make(map[string]bool)
	root := w.getRoot(node)
	w.visitChildren(root, visited)
	return len(visited) == len(w.Nodes)
}

// walk through all the children of a node and populate a map of visited children.
func (w *Workflow) visitChildren(node string, visited map[string]bool) {
	visited[node] = true
	for _, n := range w.ChildrenIDs(node) {
		w.visitChildren(n, visited)
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

func (w *Workflow) isMonoParental() bool {
	for _, node := range w.Nodes {
		if len(w.ParentIDs(node.Key)) > 1 {
			return false
		}
	}
	return true
}
