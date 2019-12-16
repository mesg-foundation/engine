package main

import (
	"context"
	"sync"
	"testing"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func testExecution(t *testing.T) {
	var (
		streamInProgress pb.Execution_StreamClient
		streamCompleted  pb.Execution_StreamClient
		err              error
		executorHash     = testRunnerHash
	)

	t.Run("create stream nil filter", func(t *testing.T) {
		_, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{})
		require.NoError(t, err)
	})

	t.Run("create stream", func(t *testing.T) {
		streamInProgress, err = client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
			Filter: &pb.StreamExecutionRequest_Filter{
				ExecutorHash: executorHash,
				Statuses:     []execution.Status{execution.Status_InProgress},
			},
		})
		require.NoError(t, err)
		streamCompleted, err = client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
			Filter: &pb.StreamExecutionRequest_Filter{
				ExecutorHash: executorHash,
				Statuses:     []execution.Status{execution.Status_Completed},
			},
		})
		require.NoError(t, err)
		acknowledgement.WaitForStreamToBeReady(streamInProgress)
		acknowledgement.WaitForStreamToBeReady(streamCompleted)
	})

	t.Run("simple execution", func(t *testing.T) {
		var (
			executionHash hash.Hash
			exec          *execution.Execution
			taskKey       = "task1"
			eventHash     = hash.Int(1)
			inputs        = &types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StringValue{
							StringValue: "test",
						},
					},
				},
			}
		)
		t.Run("create", func(t *testing.T) {
			resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
				TaskKey:      taskKey,
				EventHash:    eventHash,
				ExecutorHash: executorHash,
				Inputs:       inputs,
			})
			require.NoError(t, err)
			executionHash = resp.Hash
		})
		t.Run("in progress", func(t *testing.T) {
			execInProgress, err := streamInProgress.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, eventHash, execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("completed", func(t *testing.T) {
			exec, err = streamCompleted.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash, exec.Hash)
			require.Equal(t, taskKey, exec.TaskKey)
			require.Equal(t, eventHash, exec.EventHash)
			require.Equal(t, executorHash, exec.ExecutorHash)
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.True(t, inputs.Equal(exec.Inputs))
			require.Equal(t, "test", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
		t.Run("get", func(t *testing.T) {
			exec, err := client.ExecutionClient.Get(context.Background(), &pb.GetExecutionRequest{Hash: executionHash})
			require.NoError(t, err)
			require.True(t, exec.Equal(exec))
		})
	})

	t.Run("double execution", func(t *testing.T) {
		var (
			executionHash1 hash.Hash
			executionHash2 hash.Hash
			taskKey        = "task1"
			inputs         = &types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StringValue{
							StringValue: "test",
						},
					},
				},
			}
		)
		t.Run("create first", func(t *testing.T) {
			resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
				TaskKey:      taskKey,
				EventHash:    hash.Int(2),
				ExecutorHash: executorHash,
				Inputs:       inputs,
			})
			require.NoError(t, err)
			executionHash1 = resp.Hash
		})
		t.Run("create second", func(t *testing.T) {
			resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
				TaskKey:      taskKey,
				EventHash:    hash.Int(3),
				ExecutorHash: executorHash,
				Inputs:       inputs,
			})
			require.NoError(t, err)
			executionHash2 = resp.Hash
		})
		t.Run("first in progress", func(t *testing.T) {
			execInProgress, err := streamInProgress.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash1, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, hash.Int(2), execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("second in progress", func(t *testing.T) {
			execInProgress, err := streamInProgress.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash2, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, hash.Int(3), execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("first completed", func(t *testing.T) {
			exec, err := streamCompleted.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash1, exec.Hash)
			require.Equal(t, taskKey, exec.TaskKey)
			require.Equal(t, hash.Int(2), exec.EventHash)
			require.Equal(t, executorHash, exec.ExecutorHash)
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.True(t, inputs.Equal(exec.Inputs))
			require.Equal(t, "test", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
		t.Run("second completed", func(t *testing.T) {
			exec, err := streamCompleted.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash2, exec.Hash)
			require.Equal(t, taskKey, exec.TaskKey)
			require.Equal(t, hash.Int(3), exec.EventHash)
			require.Equal(t, executorHash, exec.ExecutorHash)
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.True(t, inputs.Equal(exec.Inputs))
			require.Equal(t, "test", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
	})

	t.Run("double execution in parallel", func(t *testing.T) {
		var (
			executions = make([]hash.Hash, 0)
			taskKey    = "task1"
			inputs     = &types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StringValue{
							StringValue: "test",
						},
					},
				},
			}
		)
		t.Run("create executions", func(t *testing.T) {
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
					TaskKey:      taskKey,
					EventHash:    hash.Int(4),
					ExecutorHash: executorHash,
					Inputs:       inputs,
				})
				require.NoError(t, err)
				executions = append(executions, resp.Hash)
			}()
			go func() {
				defer wg.Done()
				resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
					TaskKey:      taskKey,
					EventHash:    hash.Int(5),
					ExecutorHash: executorHash,
					Inputs:       inputs,
				})
				require.NoError(t, err)
				executions = append(executions, resp.Hash)
			}()
			wg.Wait()
			require.Len(t, executions, 2)
			require.False(t, executions[0].Equal(executions[1]))
		})
		t.Run("check in progress", func(t *testing.T) {
			execInProgress1, err := streamInProgress.Recv()
			require.NoError(t, err)
			execInProgress2, err := streamInProgress.Recv()
			require.NoError(t, err)
			require.False(t, execInProgress1.Hash.Equal(execInProgress2.Hash))
			require.Contains(t, executions, execInProgress1.Hash)
			require.Contains(t, executions, execInProgress2.Hash)
		})
		t.Run("check completed", func(t *testing.T) {
			exec1, err := streamCompleted.Recv()
			require.NoError(t, err)
			exec2, err := streamCompleted.Recv()
			require.False(t, exec1.Hash.Equal(exec2.Hash))
			require.Contains(t, executions, exec1.Hash)
			require.Contains(t, executions, exec2.Hash)
		})
	})

	t.Run("complex execution", func(t *testing.T) {
		var (
			executionHash hash.Hash
			exec          *execution.Execution
			taskKey       = "task_complex"
			eventHash     = hash.Int(2)
			inputs        = &types.Struct{
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
									"array": {
										Kind: &types.Value_ListValue{
											ListValue: &types.ListValue{Values: []*types.Value{
												{Kind: &types.Value_StringValue{StringValue: "first"}},
												{Kind: &types.Value_StringValue{StringValue: "second"}},
												{Kind: &types.Value_StringValue{StringValue: "third"}},
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
			resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
				TaskKey:      taskKey,
				EventHash:    eventHash,
				ExecutorHash: executorHash,
				Inputs:       inputs,
			})
			require.NoError(t, err)
			executionHash = resp.Hash
		})
		t.Run("in progress", func(t *testing.T) {
			execInProgress, err := streamInProgress.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, eventHash, execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("completed", func(t *testing.T) {
			exec, err = streamCompleted.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash, exec.Hash)
			require.Equal(t, taskKey, exec.TaskKey)
			require.Equal(t, eventHash, exec.EventHash)
			require.Equal(t, executorHash, exec.ExecutorHash)
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.True(t, inputs.Equal(exec.Inputs))
			require.Equal(t, "complex", exec.Outputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
			require.Len(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values, 3)
			require.Equal(t, "first", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[0].GetStringValue())
			require.Equal(t, "second", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[1].GetStringValue())
			require.Equal(t, "third", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[2].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["timestamp"].GetNumberValue())
		})
		t.Run("get", func(t *testing.T) {
			exec, err := client.ExecutionClient.Get(context.Background(), &pb.GetExecutionRequest{Hash: executionHash})
			require.NoError(t, err)
			require.True(t, exec.Equal(exec))
		})
	})

	t.Run("list", func(t *testing.T) {
		resp, err := client.ExecutionClient.List(context.Background(), &pb.ListExecutionRequest{})
		require.NoError(t, err)
		require.Len(t, resp.Executions, 6)
	})
}
