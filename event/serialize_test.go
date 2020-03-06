package event

import (
	"testing"

	"github.com/mesg-foundation/engine/cosmos/address"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var data = &Event{
	Key:          "eventKey",
	InstanceHash: address.InstAddress(crypto.AddressHash([]byte("10"))),
	Data: &types.Struct{
		Fields: map[string]*types.Value{
			"foo": {
				Kind: &types.Value_StringValue{StringValue: "bar"},
			},
		},
	},
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:cosmos1ffzdc9fkggz2srlgp6grj32uc9sg9qvzu4jym8;3:eventKey;4:1:foo:3:bar;;;;", data.HashSerialize())
	require.Equal(t, "cosmos186j7ddz0046caf63h566ez5nxd23svd44ga94h", address.EventAddress(crypto.AddressHash([]byte(data.HashSerialize()))).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
