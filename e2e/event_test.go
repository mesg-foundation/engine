package main

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/hash"
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

	t.Run("simple event", func(t *testing.T) {
		var (
			eventHash hash.Hash
			data      = &types.Struct{
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
			}
		)
		t.Run("create", func(t *testing.T) {
			resp, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
				InstanceHash: testInstanceHash,
				Key:          "test_event",
				Data:         data,
			})
			require.NoError(t, err)
			eventHash = resp.Hash
		})
		t.Run("receive", func(t *testing.T) {
			event, err := stream.Recv()
			require.NoError(t, err)
			require.Equal(t, eventHash, event.Hash)
			require.Equal(t, testInstanceHash, event.InstanceHash)
			require.Equal(t, "test_event", event.Key)
			require.True(t, data.Equal(event.Data))
		})
	})

	t.Run("complex event", func(t *testing.T) {
		var (
			eventHash hash.Hash
			data      = &types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StructValue{
							StructValue: &types.Struct{
								Fields: map[string]*types.Value{
									"msg": {
										Kind: &types.Value_StringValue{
											StringValue: "complex",
										},
									},
									"timestamp": {
										Kind: &types.Value_NumberValue{
											NumberValue: float64(time.Now().Unix()),
										},
									},
									"array": {
										Kind: &types.Value_ListValue{
											ListValue: &types.ListValue{Values: []*types.Value{
												&types.Value{Kind: &types.Value_StringValue{StringValue: "first"}},
												&types.Value{Kind: &types.Value_StringValue{StringValue: "second"}},
												&types.Value{Kind: &types.Value_StringValue{StringValue: "third"}},
											}},
										},
									},
								},
							},
						},
					},
				},
			}
		)
		t.Run("create", func(t *testing.T) {
			resp, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
				InstanceHash: testInstanceHash,
				Key:          "test_event_complex",
				Data:         data,
			})
			require.NoError(t, err)
			eventHash = resp.Hash
		})
		t.Run("receive", func(t *testing.T) {
			event, err := stream.Recv()
			require.NoError(t, err)
			require.Equal(t, eventHash, event.Hash)
			require.Equal(t, testInstanceHash, event.InstanceHash)
			require.Equal(t, "test_event_complex", event.Key)
			require.True(t, data.Equal(event.Data))
		})
	})
}
