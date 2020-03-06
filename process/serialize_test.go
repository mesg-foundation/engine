package process

import (
	"testing"

	"github.com/mesg-foundation/engine/cosmos/address"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var data = &Process{
	Name: "name",
	Nodes: []*Process_Node{
		{
			Key: "nodeKey1",
			Type: &Process_Node_Event_{
				&Process_Node_Event{
					InstanceHash: address.InstAddress(crypto.AddressHash([]byte("5"))),
					EventKey:     "eventKey",
				},
			},
		}, {
			Key: "nodeKey2",
			Type: &Process_Node_Task_{&Process_Node_Task{
				InstanceHash: address.InstAddress(crypto.AddressHash([]byte("2"))),
				TaskKey:      "-",
			}},
		},
	},
	Edges: []*Process_Edge{
		{Src: "nodeKey1", Dst: "nodeKey2"},
	},
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:name;4:0:1:nodeKey1;3:2:cosmos1auk3yl0r0w2zh2ksv9z72jcvvxdp7g3jseahyw;3:eventKey;;;1:1:nodeKey2;4:2:cosmos163e4uw3xtctwacplt9cchx6aqvqecp7cfc5u00;3:-;;;;5:0:1:nodeKey1;2:nodeKey2;;;", data.HashSerialize())
	require.Equal(t, "cosmos1n9awwtu5972rrjcgxp2lysk7ca2kfxf89np8wk", address.ProcAddress(crypto.AddressHash([]byte(data.HashSerialize()))).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
