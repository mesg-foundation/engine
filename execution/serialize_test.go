package execution

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var instanceHash = sdk.AccAddress(crypto.AddressHash([]byte("2")))
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

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "5:cosmos163e4uw3xtctwacplt9cchx6aqvqecp7cfc5u00;6:taskKey;7:1:a:3:b;;b:2:3.14159265359;;c:4:true;;d:6:1:1:3:hello;;;;;f:5:1:a:3:hello;;;;;;;", exec.HashSerialize())
	require.Equal(t, "cosmos1s8v8tm8tmxfpan56vspn5uflsammry0upf5xed", sdk.AccAddress(crypto.AddressHash([]byte(exec.HashSerialize()))).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec.HashSerialize()
	}
}
