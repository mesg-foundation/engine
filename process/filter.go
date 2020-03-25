package process

import (
	"errors"
	"strings"

	"github.com/mesg-foundation/engine/protobuf/types"
)

// Match returns true the current filter matches the given data
func (f *Process_Node_Filter_Condition) Match(data *types.Value) (bool, error) {
	switch f.Predicate {
	case Process_Node_Filter_Condition_EQ:
		return data.Equal(f.Value), nil
	case Process_Node_Filter_Condition_GT,
		Process_Node_Filter_Condition_GTE,
		Process_Node_Filter_Condition_LT,
		Process_Node_Filter_Condition_LTE:
		n1, ok1 := data.Kind.(*types.Value_NumberValue)
		n2, ok2 := f.Value.Kind.(*types.Value_NumberValue)
		if !ok1 || !ok2 {
			return false, errors.New("predicates GT, GTE, LT, and LTE are only compatible with type Number")
		}
		switch f.Predicate {
		case Process_Node_Filter_Condition_GT:
			return n1.NumberValue > n2.NumberValue, nil
		case Process_Node_Filter_Condition_GTE:
			return n1.NumberValue >= n2.NumberValue, nil
		case Process_Node_Filter_Condition_LT:
			return n1.NumberValue < n2.NumberValue, nil
		case Process_Node_Filter_Condition_LTE:
			return n1.NumberValue <= n2.NumberValue, nil
		}
	case Process_Node_Filter_Condition_CONTAINS:
		switch dataTyped := data.Kind.(type) {
		case *types.Value_ListValue:
			for _, value := range dataTyped.ListValue.Values {
				if value.Equal(f.Value) {
					return true, nil
				}
			}
			return false, nil
		case *types.Value_StringValue:
			filter, ok := f.Value.Kind.(*types.Value_StringValue)
			if !ok {
				return false, errors.New("predicates CONTAINS on data of type String is only compatible with value of type String")
			}
			return strings.Contains(dataTyped.StringValue, filter.StringValue), nil
		default:
			return false, errors.New("predicate CONTAINS is only compatible on data of type List or String")
		}
	}
	return false, errors.New("predicates type is unknown")
}
