package main

import (
	"testing"
)

func testEvent(t *testing.T) {
	t.SkipNow()
	// TODO: cannot create an event so this test is not functional
	// stream, err := client.EventClient.Stream(context.Background(), &orchestrator.EventStreamRequest{
	// 	Filter: &orchestrator.EventStreamRequest_Filter{},
	// }, grpc.PerRPCCredentials(&signCred{req}))
	// require.NoError(t, err)
	// acknowledgement.WaitForStreamToBeReady(stream)

	// t.Run("simple event", func(t *testing.T) {
	// 	var (
	// 		eventHash hash.Hash
	// 		data      = &types.Struct{
	// 			Fields: map[string]*types.Value{
	// 				"msg": {
	// 					Kind: &types.Value_StringValue{
	// 						StringValue: "foo",
	// 					},
	// 				},
	// 				"timestamp": {
	// 					Kind: &types.Value_NumberValue{
	// 						NumberValue: float64(time.Now().Unix()),
	// 					},
	// 				},
	// 			},
	// 		}
	// 	)
	// 	t.Run("create", func(t *testing.T) {
	// 		resp, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
	// 			InstanceHash: testInstanceHash,
	// 			Key:          "event_trigger",
	// 			Data:         data,
	// 		})
	// 		require.NoError(t, err)
	// 		eventHash = resp.Hash
	// 	})
	// 	t.Run("receive", func(t *testing.T) {
	// 		event, err := stream.Recv()
	// 		require.NoError(t, err)
	// 		require.Equal(t, eventHash, event.Hash)
	// 		require.Equal(t, testInstanceHash, event.InstanceHash)
	// 		require.Equal(t, "event_trigger", event.Key)
	// 		require.True(t, data.Equal(event.Data))
	// 	})
	// })

	// t.Run("complex event", func(t *testing.T) {
	// 	var (
	// 		eventHash hash.Hash
	// 		data      = &types.Struct{
	// 			Fields: map[string]*types.Value{
	// 				"msg": {
	// 					Kind: &types.Value_StructValue{
	// 						StructValue: &types.Struct{
	// 							Fields: map[string]*types.Value{
	// 								"msg": {
	// 									Kind: &types.Value_StringValue{
	// 										StringValue: "complex",
	// 									},
	// 								},
	// 								"timestamp": {
	// 									Kind: &types.Value_NumberValue{
	// 										NumberValue: float64(time.Now().Unix()),
	// 									},
	// 								},
	// 								"array": {
	// 									Kind: &types.Value_ListValue{
	// 										ListValue: &types.ListValue{Values: []*types.Value{
	// 											{Kind: &types.Value_StringValue{StringValue: "first"}},
	// 											{Kind: &types.Value_StringValue{StringValue: "second"}},
	// 											{Kind: &types.Value_StringValue{StringValue: "third"}},
	// 										}},
	// 									},
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		}
	// 	)
	// 	t.Run("create", func(t *testing.T) {
	// 		resp, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
	// 			InstanceHash: testInstanceHash,
	// 			Key:          "event_complex_trigger",
	// 			Data:         data,
	// 		})
	// 		require.NoError(t, err)
	// 		eventHash = resp.Hash
	// 	})
	// 	t.Run("receive", func(t *testing.T) {
	// 		event, err := stream.Recv()
	// 		require.NoError(t, err)
	// 		require.Equal(t, eventHash, event.Hash)
	// 		require.Equal(t, testInstanceHash, event.InstanceHash)
	// 		require.Equal(t, "event_complex_trigger", event.Key)
	// 		require.True(t, data.Equal(event.Data))
	// 	})
	// })
}
