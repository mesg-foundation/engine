package execution

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

var instanceHash, _ = hash.Decode("5M8pQvBCPYwzwxe2bZUbV2g8bSgpsotp441xYvVBNMhd")
var exec = New(nil, instanceHash, nil, nil, "", "taskKey", "", &types.Struct{
	Fields: map[string]*types.Value{
		"a": {
			Kind: &types.Value_StringValue{
				StringValue: "b",
			},
		},
		"b": {
			Kind: &types.Value_NumberValue{
				NumberValue: 3.14159265359,
			},
		},
		"c": {
			Kind: &types.Value_BoolValue{
				BoolValue: true,
			},
		},
		"d": {
			Kind: &types.Value_ListValue{
				ListValue: &types.ListValue{
					Values: []*types.Value{
						{
							Kind: &types.Value_NullValue{
								NullValue: types.NullValue_NULL_VALUE,
							},
						},
						{
							Kind: &types.Value_StringValue{
								StringValue: "hello",
							},
						},
					},
				},
			},
		},
		"e": {
			Kind: &types.Value_NullValue{
				NullValue: types.NullValue_NULL_VALUE,
			},
		},
		"f": {
			Kind: &types.Value_StructValue{
				StructValue: &types.Struct{
					Fields: map[string]*types.Value{
						"a": {
							Kind: &types.Value_StringValue{
								StringValue: "hello",
							},
						},
					},
				},
			},
		},
	},
}, nil, nil)

func TestSerialize(t *testing.T) {
	require.Equal(t, "5:5M8pQvBCPYwzwxe2bZUbV2g8bSgpsotp441xYvVBNMhd;6:taskKey;7:1:a:3:b;;b:2:3.14159265359;;c:4:true;;d:6:1:1:3:hello;;;;;f:5:1:a:3:hello;;;;;;;", exec.HashSerialize())
	require.Equal(t, "CNT7drUzuRuv59bbTXhoD2AUzs8vrQgi7XqfAxeHvxGf", hash.Dump(exec).String())
}

func BenchmarkSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec.HashSerialize()
	}
}
