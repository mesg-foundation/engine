package process

// TypeString returns the type of the node in string
func (node *Process_Node) TypeString() string {
	switch node.Type.(type) {
	case *Process_Node_Result_:
		return "result"
	case *Process_Node_Event_:
		return "event"
	case *Process_Node_Task_:
		return "task"
	case *Process_Node_Map_:
		return "map"
	case *Process_Node_Filter_:
		return "filter"
	default:
		return "unknown"
	}
}
