package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/service"
	"github.com/stretchr/testify/require"
)

var testServiceHash hash.Hash

func benchmarkExecution(b *testing.B) {
	b.StopTimer()
	// create service
	if testServiceHash.IsZero() {
		resp, err := client.ServiceClient.Create(context.Background(), newTestCreateServiceRequest())
		require.NoError(b, err)
		testServiceHash = resp.Hash
	}

	// start runner
	stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{
		Filter: &pb.StreamEventRequest_Filter{
			Key: "test_service_ready",
		},
	})
	require.NoError(b, err)
	acknowledgement.WaitForStreamToBeReady(stream)

	respR, err := client.RunnerClient.Create(context.Background(), &pb.CreateRunnerRequest{
		ServiceHash: testServiceHash,
		Env:         []string{"BAR=3", "REQUIRED=4"},
	})
	require.NoError(b, err)
	testRunnerHash := respR.Hash

	// wait for service to be ready
	_, err = stream.Recv()
	require.NoError(b, err)

	// init stuff
	executorHash := testRunnerHash
	streamInProgress, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		Filter: &pb.StreamExecutionRequest_Filter{
			ExecutorHash: executorHash,
			Statuses:     []execution.Status{execution.Status_InProgress},
		},
	})
	require.NoError(b, err)
	streamCompleted, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		Filter: &pb.StreamExecutionRequest_Filter{
			ExecutorHash: executorHash,
			Statuses:     []execution.Status{execution.Status_Completed},
		},
	})
	require.NoError(b, err)
	acknowledgement.WaitForStreamToBeReady(streamInProgress)
	acknowledgement.WaitForStreamToBeReady(streamCompleted)

	// run fake service
	// go func() {
	// 	for {
	// 		exec, err := streamInProgress.Recv()
	// 		if err != nil {
	// 			require.NoError(b, err)
	// 		}
	// 		go func(exec *execution.Execution) {
	// 			req := &pb.UpdateExecutionRequest{
	// 				Hash: exec.Hash,
	// 				Result: &pb.UpdateExecutionRequest_Outputs{
	// 					Outputs: &types.Struct{
	// 						Fields: map[string]*types.Value{
	// 							"msg": {
	// 								Kind: &types.Value_StringValue{
	// 									StringValue: exec.Inputs.Fields["msg"].GetStringValue(),
	// 								},
	// 							},
	// 							"timestamp": {
	// 								Kind: &types.Value_NumberValue{
	// 									NumberValue: float64(time.Now().Unix()),
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			}
	// 			if _, err := client.ExecutionClient.Update(context.Background(), req); err != nil {
	// 				require.NoError(b, err)
	// 			}
	// 		}(exec)
	// 	}
	// }()

	fmt.Printf("start benchmark for %d execution\n", b.N)
	b.StartTimer()

	// create execution
	var (
		taskKey = "task1"
		inputs  = &types.Struct{
			Fields: map[string]*types.Value{
				"msg": {
					Kind: &types.Value_StringValue{
						StringValue: "test",
					},
				},
			},
		}
	)
	for i := 0; i < b.N; i++ {
		go func(i int) {
			hash, err := hash.Random()
			require.Nil(b, err)
			fmt.Printf("creating execution #%d\n", i)
			_, err = client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
				TaskKey:      taskKey,
				EventHash:    hash,
				ExecutorHash: executorHash,
				Inputs:       inputs,
			})
			require.NoError(b, err)
			fmt.Printf("execution #%d created\n", i)
		}(i)
	}

	// wait for completion
	for i := 0; i < b.N; i++ {
		_, err = streamCompleted.Recv()
		require.NoError(b, err)
		fmt.Printf("received execution #%d\n", i)
	}

	b.StopTimer()
	// stop runner
	_, err = client.RunnerClient.Delete(context.Background(), &pb.DeleteRunnerRequest{Hash: testRunnerHash})
	require.NoError(b, err)
}

func newTestCreateServiceRequest() *pb.CreateServiceRequest {
	return &pb.CreateServiceRequest{
		Sid:  "test-service",
		Name: "test-service",
		Configuration: service.Service_Configuration{
			Env: []string{"FOO=1", "BAR=2", "REQUIRED"},
		},
		Tasks: []*service.Service_Task{
			{
				Key: "task1",
				Inputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
				},
				Outputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
					{
						Key:  "timestamp",
						Type: "Number",
					},
				},
			},
			{
				Key: "task2",
				Inputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
				},
				Outputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
					{
						Key:  "timestamp",
						Type: "Number",
					},
				},
			},
			{
				Key: "task_complex",
				Inputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "Object",
						Object: []*service.Service_Parameter{
							{
								Key:  "msg",
								Type: "String",
							},
							{
								Key:      "array",
								Type:     "String",
								Repeated: true,
								Optional: true,
							},
						},
					},
				},
				Outputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "Object",
						Object: []*service.Service_Parameter{
							{
								Key:  "msg",
								Type: "String",
							},
							{
								Key:  "timestamp",
								Type: "Number",
							},
							{
								Key:      "array",
								Type:     "String",
								Repeated: true,
								Optional: true,
							},
						},
					},
				},
			},
		},
		Events: []*service.Service_Event{
			{
				Key: "test_service_ready",
			},
			{
				Key: "test_event",
				Data: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
					{
						Key:  "timestamp",
						Type: "Number",
					},
				},
			},
			{
				Key: "test_event_complex",
				Data: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "Object",
						Object: []*service.Service_Parameter{
							{
								Key:  "msg",
								Type: "String",
							},
							{
								Key:  "timestamp",
								Type: "Number",
							},
							{
								Key:      "array",
								Type:     "String",
								Repeated: true,
							},
						},
					},
				},
			},
			{
				Key: "event_after_task",
				Data: []*service.Service_Parameter{
					{
						Key:  "task_key",
						Type: "String",
					},
					{
						Key:  "timestamp",
						Type: "Number",
					},
				},
			},
		},
		Source: "QmWHKNvJ4wT83TLHPLMjTRBjJYvcwyr8oqTtqxaJZVXPbQ",
	}
}
