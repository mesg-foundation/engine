package main

import (
	"context"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	executionmodule "github.com/mesg-foundation/engine/x/execution"
	"github.com/mesg-foundation/engine/x/ownership"
	"github.com/stretchr/testify/require"
)

func testExecution(t *testing.T) {
	var (
		streamInProgress pb.Execution_StreamClient
		streamCompleted  pb.Execution_StreamClient
		err              error
		executorHash     = testRunnerHash
		executorAddress  = testRunnerAddress
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

	t.Run("simple execution with price and withdraw", func(t *testing.T) {
		var (
			executionHash hash.Hash
			execAddress   sdk.AccAddress
			exec          *execution.Execution
			taskKey       = "task1"
			eventHash     = hash.Int(1)
			price         = sdk.NewCoins(sdk.NewInt64Coin("atto", 50000))
			inputs        = &types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StringValue{
							StringValue: "test",
						},
					},
				},
			}
			expectedCoinsForExecutor = sdk.NewCoins(sdk.NewInt64Coin("atto", 40000)) // 80%
			expectedCoinsForService  = sdk.NewCoins(sdk.NewInt64Coin("atto", 5000))  // 10%
			expectedCoinsForEmitter  = sdk.NewCoins(sdk.NewInt64Coin("atto", 5000))  // 10%
			executorBalance          sdk.Coins
			serviceBalance           sdk.Coins
		)

		lcdGet(t, "bank/balances/"+executorAddress.String(), &executorBalance)
		lcdGet(t, "bank/balances/"+testServiceAddress.String(), &serviceBalance)

		t.Run("create", func(t *testing.T) {
			msg := executionmodule.MsgCreate{
				Signer:       engineAddress,
				TaskKey:      taskKey,
				EventHash:    eventHash,
				ExecutorHash: executorHash,
				Inputs:       inputs,
				Price:        price.String(),
			}
			res, err := cclient.BuildAndBroadcastMsg(msg)
			require.NoError(t, err)
			executionHash = res.Data
		})
		t.Run("get execution address", func(t *testing.T) {
			var exec *execution.Execution
			lcdGet(t, "execution/get/"+executionHash.String(), &exec)
			require.Equal(t, exec.Hash, executionHash)
			execAddress = exec.Address
		})
		t.Run("execution balance before completed", func(t *testing.T) {
			coins := sdk.Coins{}
			lcdGet(t, "bank/balances/"+execAddress.String(), &coins)
			require.True(t, coins.IsEqual(price), price, coins)
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
			var exec *execution.Execution
			lcdGet(t, "execution/get/"+executionHash.String(), &exec)
			require.True(t, exec.Equal(exec))
		})
		t.Run("executor + emitter balance", func(t *testing.T) {
			var coins sdk.Coins
			lcdGet(t, "bank/balances/"+executorAddress.String(), &coins)
			expectedCoins := expectedCoinsForExecutor.Add(expectedCoinsForEmitter...).Add(executorBalance...)
			require.True(t, expectedCoins.IsEqual(coins), expectedCoins, coins)
		})
		t.Run("service balance", func(t *testing.T) {
			var coins sdk.Coins
			lcdGet(t, "bank/balances/"+testServiceAddress.String(), &coins)
			expectedCoins := expectedCoinsForService.Add(serviceBalance...)
			require.True(t, expectedCoins.IsEqual(coins), expectedCoins, coins)
		})
		t.Run("execution balance", func(t *testing.T) {
			var coins sdk.Coins
			lcdGet(t, "bank/balances/"+execAddress.String(), &coins)
			require.True(t, coins.IsZero(), coins)
		})
		t.Run("withdraw from service", func(t *testing.T) {
			msg := ownership.MsgWithdraw{
				Owner:        engineAddress,
				Amount:       expectedCoinsForService.String(),
				ResourceHash: testServiceHash,
			}
			_, err = cclient.BuildAndBroadcastMsg(msg)
			require.NoError(t, err)

			// check balance
			var coins sdk.Coins
			param := bank.NewQueryBalanceParams(testServiceAddress)
			require.NoError(t, cclient.QueryJSON("custom/bank/balances", param, &coins))
			require.True(t, serviceBalance.IsEqual(coins), serviceBalance, coins)
		})
		t.Run("withdraw from runner", func(t *testing.T) {
			msg := ownership.MsgWithdraw{
				Owner:        engineAddress,
				Amount:       expectedCoinsForExecutor.Add(expectedCoinsForEmitter...).String(),
				ResourceHash: testRunnerHash,
			}
			_, err = cclient.BuildAndBroadcastMsg(msg)
			require.NoError(t, err)

			// check balance
			var coins sdk.Coins
			param := bank.NewQueryBalanceParams(testRunnerAddress)
			require.NoError(t, cclient.QueryJSON("custom/bank/balances", param, &coins))
			require.True(t, executorBalance.IsEqual(coins), executorBalance, coins)
		})
	})

	t.Run("complex execution", func(t *testing.T) {
		var (
			executionHash hash.Hash
			exec          *execution.Execution
			taskKey       = "task_complex"
			eventHash     = hash.Int(2)
			price         = "10000atto"
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
			msg := executionmodule.MsgCreate{
				Signer:       engineAddress,
				TaskKey:      taskKey,
				EventHash:    eventHash,
				ExecutorHash: executorHash,
				Inputs:       inputs,
				Price:        price,
			}
			res, err := cclient.BuildAndBroadcastMsg(msg)
			require.NoError(t, err)
			executionHash = res.Data
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
			var exec *execution.Execution
			lcdGet(t, "execution/get/"+executionHash.String(), &exec)
			require.True(t, exec.Equal(exec))
		})
	})

	t.Run("list", func(t *testing.T) {
		execs := make([]*execution.Execution, 0)
		lcdGet(t, "execution/list", &execs)
		require.Len(t, execs, 2)
	})

	t.Run("many executions in parallel", func(t *testing.T) {
		var (
			n          = 10
			executions = make([]hash.Hash, 0)
			taskKey    = "task1"
			price      = "10000atto"
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
					msg := executionmodule.MsgCreate{
						Signer:       engineAddress,
						TaskKey:      taskKey,
						EventHash:    hash,
						ExecutorHash: executorHash,
						Inputs:       inputs,
						Price:        price,
					}
					res, err := cclient.BuildAndBroadcastMsg(msg)
					require.NoError(t, err)
					mutex.Lock()
					defer mutex.Unlock()
					require.NotContains(t, executions, res.Data)
					executions = append(executions, res.Data)
				}()
			}
			wg.Wait()
			require.Len(t, executions, n)
		})
		t.Run("check in progress", func(t *testing.T) {
			execs := make([]hash.Hash, 0)
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
			execs := make([]hash.Hash, 0)
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
