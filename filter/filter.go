package filter

// Predicate is the type of conditions that can be applied in a filter of a process trigger
type Predicate uint

// List of possible conditions for process's filter
const (
	EQ Predicate = iota + 1
)

// Condition is the condition to apply in a filter
type Condition struct {
	Key       string    `hash:"name:1" validate:"required,printascii"`
	Predicate Predicate `hash:"name:2" validate:"required"`
	Value     string    `hash:"name:3"`
}

// Filter contains a list of conditions to apply
type Filter struct {
	Conditions []Condition `hash:"name:1" validate:"dive,required"`
}

// Match returns true if the data match the current list of filters
func (f Filter) Match(data map[string]interface{}) bool {
	for _, condition := range f.Conditions {
		if !condition.Match(data) {
			return false
		}
	}

	return true
}

// Match returns true the current filter matches the given data
func (f Condition) Match(inputs map[string]interface{}) bool {
	switch f.Predicate {
	case EQ:
		return inputs[f.Key] == f.Value
	default:
		return false
	}
}
