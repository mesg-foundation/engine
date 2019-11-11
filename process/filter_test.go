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
		data   []*types.Value
		match  bool
	}{
		{
			name: "not matching filter",
			filter: Process_Node_Filter{
				Conditions: []Process_Node_Filter_Condition{
					{
						Index:     0,
						Predicate: Process_Node_Filter_Condition_EQ,
						Value:     "xx",
					},
				},
			},
			data: []*types.Value{
				{
					Kind: &types.Value_StringValue{
						StringValue: "bar",
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
						Index:     0,
						Predicate: Process_Node_Filter_Condition_EQ,
						Value:     "bar",
					},
					{
						Index:     1,
						Predicate: Process_Node_Filter_Condition_EQ,
						Value:     "yyy",
					},
				},
			},
			data: []*types.Value{
				{
					Kind: &types.Value_StringValue{
						StringValue: "bar",
					},
				},
				{
					Kind: &types.Value_StringValue{
						StringValue: "yyy",
					},
				},
				{
					Kind: &types.Value_StringValue{
						StringValue: "bbb",
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
						Index:     0,
						Predicate: Process_Node_Filter_Condition_EQ,
						Value:     "bar",
					},
					{
						Index:     1,
						Predicate: Process_Node_Filter_Condition_EQ,
						Value:     "aaa",
					},
				},
			},
			data: []*types.Value{
				{
					Kind: &types.Value_StringValue{
						StringValue: "bar",
					},
				},
				{
					Kind: &types.Value_StringValue{
						StringValue: "yyy",
					},
				},
				{
					Kind: &types.Value_StringValue{
						StringValue: "bbb",
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
