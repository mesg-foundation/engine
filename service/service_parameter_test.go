package service

import (
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	var tests = []struct {
		name      string
		params    []*Service_Parameter
		data      *types.Struct
		expecterr bool
	}{
		{
			name: "no parameters and no data",
		},
		{
			name: "one optional parameter without data",
			params: []*Service_Parameter{
				{
					Optional: true,
				},
			},
		},
		{
			name: "simple types (string,number,boolean)",
			params: []*Service_Parameter{
				{
					Key:  "string",
					Type: "String",
				},
				{
					Key:  "number",
					Type: "Number",
				},
				{
					Key:  "boolean",
					Type: "Boolean",
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"string": {
						Kind: &types.Value_StringValue{},
					},
					"number": {
						Kind: &types.Value_NumberValue{},
					},
					"boolean": {
						Kind: &types.Value_BoolValue{},
					},
				},
			},
		},
		{
			name: "any type",
			params: []*Service_Parameter{
				{
					Key:  "key",
					Type: "Any",
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"key": {
						Kind: &types.Value_ListValue{},
					},
				},
			},
		},
		{
			name: "array type",
			params: []*Service_Parameter{
				{
					Key:      "key",
					Type:     "Number",
					Repeated: true,
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"key": {
						Kind: &types.Value_ListValue{
							ListValue: &types.ListValue{
								Values: []*types.Value{
									{
										Kind: &types.Value_NumberValue{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "object type",
			params: []*Service_Parameter{
				{
					Key:  "key",
					Type: "Object",
					Object: []*Service_Parameter{
						{
							Key:  "string",
							Type: "String",
						},
					},
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"key": {
						Kind: &types.Value_StructValue{
							StructValue: &types.Struct{
								Fields: map[string]*types.Value{
									"string": {
										Kind: &types.Value_StringValue{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "required parameter without data",
			params: []*Service_Parameter{
				{
					Optional: false,
				},
			},
			expecterr: true,
		},
		{
			name: "invalid parameter type",
			params: []*Service_Parameter{
				{
					Key:  "key",
					Type: "-",
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"key": {
						Kind: &types.Value_NullValue{},
					},
				},
			},
			expecterr: true,
		},
		{
			name: "invalid string type",
			params: []*Service_Parameter{
				{
					Key:  "key",
					Type: "String",
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"key": {
						Kind: &types.Value_NullValue{},
					},
				},
			},
			expecterr: true,
		},
		{
			name: "invalid number type",
			params: []*Service_Parameter{
				{
					Key:  "key",
					Type: "Number",
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"key": {
						Kind: &types.Value_NullValue{},
					},
				},
			},
			expecterr: true,
		},
		{
			name: "invalid boolean type",
			params: []*Service_Parameter{
				{
					Key:  "key",
					Type: "Boolean",
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"key": {
						Kind: &types.Value_NullValue{},
					},
				},
			},
			expecterr: true,
		},
		{
			name: "invalid list type",
			params: []*Service_Parameter{
				{
					Key:      "key",
					Repeated: true,
					Type:     "String",
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"key": {
						Kind: &types.Value_NullValue{},
					},
				},
			},
			expecterr: true,
		},
		{
			name: "invalid object type",
			params: []*Service_Parameter{
				{
					Key:  "key",
					Type: "Object",
				},
			},
			data: &types.Struct{
				Fields: map[string]*types.Value{
					"key": {
						Kind: &types.Value_NullValue{},
					},
				},
			},
			expecterr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.expecterr {
				assert.Error(t, validateServiceParameters(tt.params, tt.data))
			} else {
				assert.NoError(t, validateServiceParameters(tt.params, tt.data))
			}
		})
	}
}
