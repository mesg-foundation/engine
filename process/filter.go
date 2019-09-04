package process

// Match returns true if the data match the current list of filters
func (f Process_Node_Filter) Match(data map[string]interface{}) bool {
	for _, condition := range f.Conditions {
		if !condition.Match(data) {
			return false
		}
	}

	return true
}

// Match returns true the current filter matches the given data
func (f Process_Node_Filter_Condition) Match(inputs map[string]interface{}) bool {
	return f.Predicate == Process_Node_Filter_Condition_EQ &&
		inputs[f.Key] == f.Value
}
