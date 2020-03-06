package main

import (
	"context"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos/address"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
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
			executionHash address.ExecAddress
			exec          *execution.Execution
			taskKey       = "task1"
			eventHash     = address.EventAddress(crypto.AddressHash([]byte("1")))
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
			require.True(t, exec.Hash.Equals(exec.Hash))
		})
	})

	t.Run("complex execution", func(t *testing.T) {
		var (
			executionHash address.ExecAddress
			exec          *execution.Execution
			taskKey       = "task_complex"
			eventHash     = address.EventAddress(crypto.AddressHash([]byte("2")))
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
			require.True(t, exec.Hash.Equals(exec.Hash))
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("grpc", func(t *testing.T) {
			resp, err := client.ExecutionClient.List(context.Background(), &pb.ListExecutionRequest{})
			require.NoError(t, err)
			require.Len(t, resp.Executions, 2)
		})
		t.Run("lcd", func(t *testing.T) {
			execs := make([]*execution.Execution, 0)
			lcdGet(t, "execution/list", &execs)
			require.Len(t, execs, 2)
		})
	})

	t.Run("execution with price", func(t *testing.T) {
		var (
			inputs = &types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StringValue{
							StringValue: "test",
						},
					},
				},
			}
		)
		resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
			TaskKey:      "task1",
			EventHash:    address.EventAddress(crypto.AddressHash([]byte("1"))),
			ExecutorHash: executorHash,
			Price:        "50000atto",
			Inputs:       inputs,
		})
		require.NoError(t, err)

		// check balance of execution before completed
		t.Run("execution balance before completed", func(t *testing.T) {
			coins := sdk.Coins{}
			lcdGet(t, "bank/balances/"+resp.Hash.String(), &coins)
			require.True(t, coins.AmountOf("atto").Equal(sdk.NewInt(50000)))
		})

		var executorBalance sdk.Coins
		var serviceBalance sdk.Coins
		lcdGet(t, "bank/balances/"+executorHash.String(), &executorBalance)
		lcdGet(t, "bank/balances/"+testServiceHash.String(), &serviceBalance)

		_, err = streamInProgress.Recv()
		require.NoError(t, err)

		exec, err := streamCompleted.Recv()
		require.NoError(t, err)
		require.Equal(t, resp.Hash, exec.Hash)

		// check balance of executor
		t.Run("executor balance", func(t *testing.T) {
			coins := sdk.Coins{}
			lcdGet(t, "bank/balances/"+executorHash.String(), &coins)
			require.True(t, coins.AmountOf("atto").Equal(sdk.NewInt(45000).Add(executorBalance.AmountOf("atto"))))
		})
		// check balance of service
		t.Run("service balance", func(t *testing.T) {
			coins := sdk.Coins{}
			lcdGet(t, "bank/balances/"+testServiceHash.String(), &coins)
			require.True(t, coins.AmountOf("atto").Equal(sdk.NewInt(5000).Add(serviceBalance.AmountOf("atto"))))
		})
		// check balance of execution
		t.Run("execution balance", func(t *testing.T) {
			coins := sdk.Coins{}
			lcdGet(t, "bank/balances/"+resp.Hash.String(), &coins)
			require.True(t, coins.AmountOf("atto").Equal(sdk.NewInt(0)))
		})
	})

	t.Run("many executions in parallel", func(t *testing.T) {
		var (
			n          = 10
			executions = make([]address.ExecAddress, 0)
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
			var mutex sync.Mutex
			wg.Add(n)
			for i := 1; i <= n; i++ {
				go func() {
					defer wg.Done()
					hash, err := hash.Random()
					require.Nil(t, err)
					resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
						TaskKey:      taskKey,
						EventHash:    address.EventAddress(crypto.AddressHash(hash)),
						ExecutorHash: executorHash,
						Inputs:       inputs,
					})
					require.NoError(t, err)
					mutex.Lock()
					defer mutex.Unlock()
					require.NotContains(t, executions, resp.Hash)
					executions = append(executions, resp.Hash)
				}()
			}
			wg.Wait()
			require.Len(t, executions, n)
		})
		t.Run("check in progress", func(t *testing.T) {
			execs := make([]address.ExecAddress, 0)
			for i := 1; i <= n; i++ {
				exec, err := streamInProgress.Recv()
				require.NoError(t, err)
				require.Contains(t, executions, exec.Hash)
				require.NotContains(t, execs, exec.Hash)
				execs = append(execs, exec.Hash)
			}
			require.Len(t, execs, n)
		})
		t.Run("check completed", func(t *testing.T) {
			execs := make([]address.ExecAddress, 0)
			for i := 1; i <= n; i++ {
				exec, err := streamCompleted.Recv()
				require.NoError(t, err)
				require.Contains(t, executions, exec.Hash)
				require.NotContains(t, execs, exec.Hash)
				execs = append(execs, exec.Hash)
			}
			require.Len(t, execs, n)
		})
	})
}
