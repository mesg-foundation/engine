package process

// NodeType are the string constant corresponding to each type of node.
const (
	NodeTypeUnknown = "unknown"
	NodeTypeResult  = "result"
	NodeTypeEvent   = "event"
	NodeTypeTask    = "task"
	NodeTypeMap     = "map"
	NodeTypeFilter  = "filter"
)

// TypeString returns the type of the node in string
func (node *Process_Node) TypeString() string {
	switch node.Type.(type) {
	case *Process_Node_Result_:
		return NodeTypeResult
	case *Process_Node_Event_:
		return NodeTypeEvent
	case *Process_Node_Task_:
		return NodeTypeTask
	case *Process_Node_Map_:
		return NodeTypeMap
	case *Process_Node_Filter_:
		return NodeTypeFilter
	default:
		return NodeTypeUnknown
	}
}
