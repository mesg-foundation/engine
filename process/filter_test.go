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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.match, tt.filter.Match(tt.data))
		})
	}
}
