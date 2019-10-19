package process

import "github.com/mesg-foundation/engine/protobuf/types"

// Match returns true if the data match the current list of filters
func (f Process_Node_Filter) Match(data *types.Struct) bool {
	for _, condition := range f.Conditions {
		if !condition.Match(data) {
			return false
		}
	}

	return true
}

// Match returns true the current filter matches the given data
func (f Process_Node_Filter_Condition) Match(data *types.Struct) bool {
	return f.Predicate == Process_Node_Filter_Condition_PREDICATE_EQ &&
		data.Fields[f.Key].GetStringValue() == f.Value
}
