package process

import (
	"testing"

	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	var tests = []struct {
		name   string
		filter Process_Node_Filter
		data   *types.Struct
		match  bool
	}{
		{
			name: "not matching filter",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_EQ,
						Value: &types.Value{
							Kind: &types.Value_StringValue{
								StringValue: "xx",
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_StringValue{
							StringValue: "bar",
						},
					},
				},
			},
			match: false,
		},
		{
			name: "matching multiple conditions",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_EQ,
						Value: &types.Value{
							Kind: &types.Value_StringValue{
								StringValue: "bar",
							},
						},
					},
					{
						Key:       "xxx",
						Predicate: Process_Node_Filter_Condition_EQ,
						Value: &types.Value{
							Kind: &types.Value_StringValue{
								StringValue: "yyy",
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_StringValue{
							StringValue: "bar",
						},
					},
					"xxx": {
						Kind: &types.Value_StringValue{
							StringValue: "yyy",
						},
					},
					"aaa": {
						Kind: &types.Value_StringValue{
							StringValue: "bbb",
						},
					},
				},
			},
			match: true,
		},
		{
			name: "non matching multiple conditions",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_EQ,
						Value: &types.Value{
							Kind: &types.Value_StringValue{
								StringValue: "bar",
							},
						},
					},
					{
						Key:       "xxx",
						Predicate: Process_Node_Filter_Condition_EQ,
						Value: &types.Value{
							Kind: &types.Value_StringValue{
								StringValue: "aaa",
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_StringValue{
							StringValue: "bar",
						},
					},
					"xxx": {
						Kind: &types.Value_StringValue{
							StringValue: "yyy",
						},
					},
					"aaa": {
						Kind: &types.Value_StringValue{
							StringValue: "bbb",
						},
					},
				},
			},
			match: false,
		},
		{
			name: "matching GT",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_GT,
						Value: &types.Value{
							Kind: &types.Value_NumberValue{
								NumberValue: 1,
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_NumberValue{
							NumberValue: 2,
						},
					},
				},
			},
			match: true,
		},
		{
			name: "non matching GT",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_GT,
						Value: &types.Value{
							Kind: &types.Value_NumberValue{
								NumberValue: 2,
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_NumberValue{
							NumberValue: 1,
						},
					},
				},
			},
			match: false,
		},
		{
			name: "GT wrong type",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_LTE,
						Value: &types.Value{
							Kind: &types.Value_NumberValue{
								NumberValue: 1,
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_StringValue{
							StringValue: "foo",
						},
					},
				},
			},
			match: false,
		},
		{
			name: "matching GTE",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_GTE,
						Value: &types.Value{
							Kind: &types.Value_NumberValue{
								NumberValue: 1,
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_NumberValue{
							NumberValue: 1,
						},
					},
				},
			},
			match: true,
		},
		{
			name: "non matching GTE",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_GTE,
						Value: &types.Value{
							Kind: &types.Value_NumberValue{
								NumberValue: 2,
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_NumberValue{
							NumberValue: 1,
						},
					},
				},
			},
			match: false,
		},
		{
			name: "matching LT",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_LT,
						Value: &types.Value{
							Kind: &types.Value_NumberValue{
								NumberValue: 2,
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_NumberValue{
							NumberValue: 1,
						},
					},
				},
			},
			match: true,
		},
		{
			name: "non matching LT",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_LT,
						Value: &types.Value{
							Kind: &types.Value_NumberValue{
								NumberValue: 1,
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_NumberValue{
							NumberValue: 2,
						},
					},
				},
			},
			match: false,
		},
		{
			name: "matching LTE",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_LTE,
						Value: &types.Value{
							Kind: &types.Value_NumberValue{
								NumberValue: 1,
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_NumberValue{
							NumberValue: 1,
						},
					},
				},
			},
			match: true,
		},
		{
			name: "non matching LTE",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_LTE,
						Value: &types.Value{
							Kind: &types.Value_NumberValue{
								NumberValue: 1,
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_NumberValue{
							NumberValue: 2,
						},
					},
				},
			},
			match: false,
		},
		{
			name: "matching contains",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_CONTAINS,
						Value: &types.Value{
							Kind: &types.Value_StringValue{
								StringValue: "foo",
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
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
				},
			},
			match: true,
		},
		{
			name: "non matching contains",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_CONTAINS,
						Value: &types.Value{
							Kind: &types.Value_StringValue{
								StringValue: "foo",
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
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
				},
			},
			match: false,
		},
		{
			name: "wrong type contains",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Key:       "foo",
						Predicate: Process_Node_Filter_Condition_CONTAINS,
						Value: &types.Value{
							Kind: &types.Value_StringValue{
								StringValue: "foo",
							},
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"foo": {
						Kind: &types.Value_StringValue{
							StringValue: "foo",
						},
					},
				},
			},
			match: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.match, tt.filter.Match(tt.data))
		})
	}
}
