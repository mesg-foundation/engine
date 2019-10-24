package main

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestExecution(t *testing.T) {
	stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{
		Filter: &pb.StreamEventRequest_Filter{},
	})
	require.NoError(t, err)
	go func() {
		fmt.Println("RECIVE EVENT ")
		event, err := stream.Recv()
		fmt.Println("RECIVE EVENT ")
		require.NoError(t, err)
		require.Equal(t, "ping_ok", event.Key)
	}()

	t.Run("create", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), passmd)
		_, err := client.ExecutionClient.Create(ctx, &pb.CreateExecutionRequest{
			InstanceHash: testInstanceHash,
			TaskKey:      "ping",
			Inputs: &types.Struct{
				Fields: map[string]*types.Value{
					"msg": &types.Value{
						Kind: &types.Value_StringValue{
							StringValue: "test",
						},
					},
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("get", func(t *testing.T) {
	})
	t.Run("stream", func(t *testing.T) {
	})
	t.Run("update", func(t *testing.T) {
	})
}
