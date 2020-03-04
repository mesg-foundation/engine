package event

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

var data = &Event{
	Key:          "eventKey",
	InstanceHash: hash.Int(10),
	Data: &types.Struct{
		Fields: map[string]*types.Value{
			"foo": &types.Value{
				Kind: &types.Value_StringValue{StringValue: "bar"},
			},
		},
	},
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:g35TxFqwMx95vCk63fTxGTHb6ei4W24qg5t2x6xD3cT;3:eventKey;4:1:foo:3:bar;;;;", data.HashSerialize())
	require.Equal(t, "CGQ1DWeSsf13BDovLNb9zHXTVMUPNcTWY4DaSAFNN88T", hash.Dump(data).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
