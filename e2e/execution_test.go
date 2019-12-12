package main

import (
	"context"
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
		stream       pb.Execution_StreamClient
		err          error
		executorHash = testRunnerHash
	)

	t.Run("create stream nil filter", func(t *testing.T) {
		_, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{})
		require.NoError(t, err)
	})

	t.Run("create stream", func(t *testing.T) {
		stream, err = client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
			Filter: &pb.StreamExecutionRequest_Filter{
				ExecutorHash: executorHash,
			},
		})
		require.NoError(t, err)
		acknowledgement.WaitForStreamToBeReady(stream)
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
			execInProgress, err := stream.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, eventHash, execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("completed", func(t *testing.T) {
			exec, err = stream.Recv()
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
			execInProgress, err := stream.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash1, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, hash.Int(2), execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("second in progress", func(t *testing.T) {
			execInProgress, err := stream.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash2, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, hash.Int(3), execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("first completed", func(t *testing.T) {
			exec, err := stream.Recv()
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
			exec, err := stream.Recv()
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
			execInProgress, err := stream.Recv()
			require.NoError(t, err)
			require.Equal(t, executionHash, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, eventHash, execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("completed", func(t *testing.T) {
			exec, err = stream.Recv()
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
		require.Len(t, resp.Executions, 4)
	})
}
