package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/result"
	"github.com/stretchr/testify/require"
)

func testExecuteTask(t *testing.T) {
	var (
		executionStream pb.Execution_StreamClient
		resultStream    pb.Result_StreamClient
		err             error
		executorHash    = testRunnerHash
	)

	t.Run("create stream execution", func(t *testing.T) {
		t.Run("with nil filter", func(t *testing.T) {
			_, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{})
			require.NoError(t, err)
		})

		t.Run("normal", func(t *testing.T) {
			executionStream, err = client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
				Filter: &pb.StreamExecutionRequest_Filter{
					ExecutorHash: executorHash,
				},
			})
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(executionStream)
		})
	})

	t.Run("create stream result", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			resultStream, err = client.ResultClient.Stream(context.Background(), &pb.StreamResultRequest{})
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(resultStream)
		})
	})

	t.Run("simple execution", func(t *testing.T) {
		var (
			execHash  hash.Hash
			taskKey   = "task1"
			eventHash = hash.Int(1)
			inputs    = &types.Struct{
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
			execHash = resp.Hash
		})
		var exec *execution.Execution
		t.Run("in progress", func(t *testing.T) {
			exec, err = executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, execHash, exec.Hash)
			require.Equal(t, taskKey, exec.TaskKey)
			require.Equal(t, eventHash, exec.EventHash)
			require.Equal(t, executorHash, exec.ExecutorHash)
			require.True(t, inputs.Equal(exec.Inputs))
		})
		t.Run("get request", func(t *testing.T) {
			exec, err := client.ExecutionClient.Get(context.Background(), &pb.GetExecutionRequest{Hash: execHash})
			require.NoError(t, err)
			require.True(t, exec.Equal(exec))
		})
		var res *result.Result
		t.Run("completed", func(t *testing.T) {
			res, err = resultStream.Recv()
			require.NoError(t, err)
			require.True(t, res.RequestHash.Equal(execHash))
			require.Equal(t, "test", res.GetOutputs().Fields["msg"].GetStringValue())
			require.NotEmpty(t, res.GetOutputs().Fields["timestamp"].GetNumberValue())
		})
		t.Run("get result", func(t *testing.T) {
			exec, err := client.ResultClient.Get(context.Background(), &pb.GetResultRequest{Hash: res.Hash})
			require.NoError(t, err)
			require.True(t, exec.Equal(res))
		})
	})

	t.Run("complex execution", func(t *testing.T) {
		var (
			execHash  hash.Hash
			taskKey   = "task_complex"
			eventHash = hash.Int(2)
			inputs    = &types.Struct{
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
			execHash = resp.Hash
		})
		var exec *execution.Execution
		t.Run("in progress", func(t *testing.T) {
			exec, err = executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, execHash, exec.Hash)
			require.Equal(t, taskKey, exec.TaskKey)
			require.Equal(t, eventHash, exec.EventHash)
			require.Equal(t, executorHash, exec.ExecutorHash)
			require.True(t, inputs.Equal(exec.Inputs))
		})
		t.Run("get request", func(t *testing.T) {
			exec, err := client.ExecutionClient.Get(context.Background(), &pb.GetExecutionRequest{Hash: execHash})
			require.NoError(t, err)
			require.True(t, exec.Equal(exec))
		})
		var res *result.Result
		t.Run("completed", func(t *testing.T) {
			res, err = resultStream.Recv()
			require.NoError(t, err)
			require.True(t, res.RequestHash.Equal(execHash))
			require.Equal(t, "complex", res.GetOutputs().Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
			require.Len(t, res.GetOutputs().Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values, 3)
			require.Equal(t, "first", res.GetOutputs().Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[0].GetStringValue())
			require.Equal(t, "second", res.GetOutputs().Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[1].GetStringValue())
			require.Equal(t, "third", res.GetOutputs().Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[2].GetStringValue())
			require.NotEmpty(t, res.GetOutputs().Fields["msg"].GetStructValue().Fields["timestamp"].GetNumberValue())
		})
		t.Run("get result", func(t *testing.T) {
			exec, err := client.ResultClient.Get(context.Background(), &pb.GetResultRequest{Hash: res.Hash})
			require.NoError(t, err)
			require.True(t, exec.Equal(res))
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("list request", func(t *testing.T) {
			resp, err := client.ExecutionClient.List(context.Background(), &pb.ListExecutionRequest{})
			require.NoError(t, err)
			require.Len(t, resp.Executions, 2)
		})
		t.Run("list result", func(t *testing.T) {
			resp, err := client.ResultClient.List(context.Background(), &pb.ListResultRequest{})
			require.NoError(t, err)
			require.Len(t, resp.Results, 2)
		})
	})
}
