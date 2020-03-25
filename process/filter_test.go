package process

import (
	"testing"

	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	var tests = []struct {
		name      string
		condition *Process_Node_Filter_Condition
		data      *types.Value
		match     bool
		err       string
	}{
		{
			name: "not matching filter",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_EQ,
				Value: &types.Value{
					Kind: &types.Value_StringValue{
						StringValue: "xx",
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_StringValue{
					StringValue: "bar",
				},
			},
			match: false,
		},
		{
			name: "matching GT",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_GT,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 1,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_NumberValue{
					NumberValue: 2,
				},
			},
			match: true,
		},
		{
			name: "non matching GT",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_GT,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 2,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_NumberValue{
					NumberValue: 1,
				},
			},
			match: false,
		},
		{
			name: "GT wrong type",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_LTE,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 1,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_StringValue{
					StringValue: "foo",
				},
			},
			match: false,
			err:   "predicates GT, GTE, LT, and LTE are only compatible with type Number",
		},
		{
			name: "matching GTE",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_GTE,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 1,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_NumberValue{
					NumberValue: 1,
				},
			},
			match: true,
		},
		{
			name: "non matching GTE",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_GTE,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 2,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_NumberValue{
					NumberValue: 1,
				},
			},
			match: false,
		},
		{
			name: "matching LT",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_LT,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 2,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_NumberValue{
					NumberValue: 1,
				},
			},
			match: true,
		},
		{
			name: "non matching LT",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_LT,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 1,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_NumberValue{
					NumberValue: 2,
				},
			},
			match: false,
		},
		{
			name: "matching LTE",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_LTE,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 1,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_NumberValue{
					NumberValue: 1,
				},
			},
			match: true,
		},
		{
			name: "non matching LTE",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_LTE,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 1,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_NumberValue{
					NumberValue: 2,
				},
			},
			match: false,
		},
		{
			name: "matching contains",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_CONTAINS,
				Value: &types.Value{
					Kind: &types.Value_StringValue{
						StringValue: "foo",
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_ListValue{
					ListValue: &types.ListValue{
						Values: []*types.Value{
							{
								Kind: &types.Value_StringValue{
									StringValue: "foo",
								},
							},
						},
					},
				},
			},
			match: true,
		},
		{
			name: "non matching contains",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_CONTAINS,
				Value: &types.Value{
					Kind: &types.Value_StringValue{
						StringValue: "foo",
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_ListValue{
					ListValue: &types.ListValue{
						Values: []*types.Value{
							{
								Kind: &types.Value_StringValue{
									StringValue: "bar",
								},
							},
						},
					},
				},
			},
			match: false,
		},
		{
			name: "wrong type contains",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_CONTAINS,
				Value: &types.Value{
					Kind: &types.Value_StringValue{
						StringValue: "foo",
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_NumberValue{
					NumberValue: 10,
				},
			},
			match: false,
			err:   "predicate CONTAINS is only compatible on data of type List or String",
		},
		{
			name: "wrong type contains string",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_CONTAINS,
				Value: &types.Value{
					Kind: &types.Value_NumberValue{
						NumberValue: 10,
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_StringValue{
					StringValue: "foo",
				},
			},
			match: false,
			err:   "predicates CONTAINS on data of type String is only compatible with value of type String",
		},
		{
			name: "string contain",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_CONTAINS,
				Value: &types.Value{
					Kind: &types.Value_StringValue{
						StringValue: "world",
					},
				},
			},
			data: &types.Value{
				Kind: &types.Value_StringValue{
					StringValue: "hello world",
				},
			},
			match: true,
		},
		{
			name: "unknown predicate",
			condition: &Process_Node_Filter_Condition{
				Predicate: Process_Node_Filter_Condition_Unknown,
			},
			match: false,
			err:   "predicates type is unknown",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			match, err := tt.condition.Match(tt.data)
			require.Equal(t, tt.match, match)
			if len(tt.err) > 0 || err != nil {
				require.EqualError(t, err, tt.err)
			}
		})
	}
}
