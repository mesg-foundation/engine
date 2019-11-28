package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
)

var testProcessHash hash.Hash

func testProcess(t *testing.T) {
	req := &pb.CreateProcessRequest{
		Name: "test-process",
		Nodes: []*process.Process_Node{
			{
				Key: "n0",
				Type: &process.Process_Node_Event_{
					Event: &process.Process_Node_Event{
						InstanceHash: testInstanceHash,
						EventKey:     "test_service_ready",
					},
				},
			},
			{
				Key: "n1",
				Type: &process.Process_Node_Task_{
					Task: &process.Process_Node_Task{
						InstanceHash: testInstanceHash,
						TaskKey:      "test_service_ready",
					},
				},
			},
		},
		Edges: []*process.Process_Edge{
			{
				Src: "n0",
				Dst: "n1",
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		resp, err := client.ProcessClient.Create(context.Background(), req)
		require.NoError(t, err)
		testProcessHash = resp.Hash
	})
	t.Run("get", func(t *testing.T) {
		p, err := client.ProcessClient.Get(context.Background(), &pb.GetProcessRequest{Hash: testProcessHash})
		require.NoError(t, err)
		require.True(t, p.Equal(&process.Process{
			Hash:  p.Hash,
			Name:  req.Name,
			Nodes: req.Nodes,
			Edges: req.Edges,
		}))
	})

	t.Run("list", func(t *testing.T) {
		ps, err := client.ProcessClient.List(context.Background(), &pb.ListProcessRequest{})
		require.NoError(t, err)
		require.Len(t, ps.Processes, 1)
	})

	t.Run("delete", func(t *testing.T) {
		_, err := client.ProcessClient.Delete(context.Background(), &pb.DeleteProcessRequest{Hash: testProcessHash})
		require.NoError(t, err)
	})
}
