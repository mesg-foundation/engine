package workflow

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
	w.visitChildren(node, visited)
	return len(visited) == len(w.Nodes)
}

// walk through all the children of a node and populate a map of visited children.
func (w *Workflow) visitChildren(node string, visited map[string]bool) {
	visited[node] = true
	for _, edge := range w.ChildrenIDs(node) {
		w.visitChildren(edge, visited)
	}
}
