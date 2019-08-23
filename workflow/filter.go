package workflow

// // Predicate is the type of conditions that can be applied in a filter of a workflow trigger
// type Predicate uint

// // List of possible conditions for workflow's filter
// const (
// 	EQ Predicate = iota + 1
// )

// // TriggerFilters is a list of filters to apply
// type TriggerFilters []*TriggerFilter

// // TriggerFilter is the filter definition that can be applied to a workflow trigger
// type TriggerFilter struct {
// 	Key       string      `hash:"name:1" validate:"required,printascii"`
// 	Predicate Predicate   `hash:"name:2" validate:"required"`
// 	Value     interface{} `hash:"name:3"`
// }

// // Match returns true if the data match the current list of filters
// func (f TriggerFilters) Match(data map[string]interface{}) bool {
// 	filters := []*TriggerFilter(f)
// 	for _, filter := range filters {
// 		if !filter.Match(data) {
// 			return false
// 		}
// 	}

// 	return true
// }

// // Match returns true the current filter matches the given data
// func (f *TriggerFilter) Match(inputs map[string]interface{}) bool {
// 	switch f.Predicate {
// 	case EQ:
// 		return inputs[f.Key] == f.Value
// 	default:
// 		return false
// 	}
// }
