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
				Kind: &types.Value_StringValue{
					StringValue: "foo",
				},
			},
			match: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.match, tt.condition.Match(tt.data))
		})
	}
}
