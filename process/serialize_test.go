package process

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
)

var data = &Process{
	Name: "name",
	Nodes: []*Process_Node{
		&Process_Node{
			Key: "nodeKey1",
			Type: &Process_Node_Event_{
				&Process_Node_Event{
					InstanceHash: hash.Int(5),
					EventKey:     "eventKey",
				},
			},
		}, &Process_Node{
			Key: "nodeKey2",
			Type: &Process_Node_Task_{&Process_Node_Task{
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}},
		},
	},
	Edges: []*Process_Edge{
		&Process_Edge{Src: "nodeKey1", Dst: "nodeKey2"},
	},
}

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "2:name;4:0:1:nodeKey1;3:2:LX3EUdRUBUa3TbsYXLEUdj9J3prXkWXvLYSWyYyc2Jj;3:eventKey;;;1:1:nodeKey2;4:2:8opHzTAnfzRpPEx21XtnrVTX28YQuCpAjcn1PczScKh;3:-;;;;5:0:1:nodeKey1;2:nodeKey2;;;", data.HashSerialize())
	require.Equal(t, "HUQ6EKW3fQPDhpCZu65x5ETQxtSCDKLJk9Qw9PKkLvzg", hash.Dump(data).String())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data.HashSerialize()
	}
}
