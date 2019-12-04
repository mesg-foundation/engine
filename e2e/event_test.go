package main

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func testEvent(t *testing.T) {
	stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{
		Filter: &pb.StreamEventRequest_Filter{},
	})
	require.NoError(t, err)
	acknowledgement.WaitForStreamToBeReady(stream)

	resp, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
		InstanceHash: testInstanceHash,
		Key:          "test_event",
		Data: &types.Struct{
			Fields: map[string]*types.Value{
				"msg": {
					Kind: &types.Value_StringValue{
						StringValue: "foo",
					},
				},
				"timestamp": {
					Kind: &types.Value_NumberValue{
						NumberValue: float64(time.Now().Unix()),
					},
				},
			},
		},
	})
	require.NoError(t, err)

	event, err := stream.Recv()
	require.NoError(t, err)

	require.Equal(t, resp.Hash, event.Hash)
	require.Equal(t, testInstanceHash, event.InstanceHash)
	require.Equal(t, "test_event", event.Key)
	require.Equal(t, "foo", event.Data.Fields["msg"].GetStringValue())
	require.NotEmpty(t, event.Data.Fields["timestamp"].GetNumberValue())
}
