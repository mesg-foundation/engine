package main

import (
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	executionmodule "github.com/mesg-foundation/engine/x/execution"
	"github.com/mesg-foundation/engine/x/ownership"
	"github.com/stretchr/testify/require"
)

func testExecution(t *testing.T) {
	var (
		executorHash    = testRunnerHash
		executorAddress = testRunnerAddress
	)

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
			executionHash = lcdBroadcastMsg(t, msg)
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
			execInProgress := pollExecution(t, executionHash, execution.Status_InProgress)
			require.Equal(t, executionHash, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, eventHash, execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("completed", func(t *testing.T) {
			exec = pollExecution(t, executionHash, execution.Status_Completed)
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
			var execR *execution.Execution
			lcdGet(t, "execution/get/"+executionHash.String(), &execR)
			require.True(t, exec.Equal(execR), exec, execR)
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
			lcdBroadcastMsg(t, msg)

			// check balance
			var coins sdk.Coins
			lcdGet(t, "bank/balances/"+testServiceAddress.String(), &coins)
			require.True(t, serviceBalance.IsEqual(coins), serviceBalance, coins)
		})
		t.Run("withdraw from runner", func(t *testing.T) {
			msg := ownership.MsgWithdraw{
				Owner:        engineAddress,
				Amount:       expectedCoinsForExecutor.Add(expectedCoinsForEmitter...).String(),
				ResourceHash: testRunnerHash,
			}
			lcdBroadcastMsg(t, msg)

			// check balance
			var coins sdk.Coins
			lcdGet(t, "bank/balances/"+testRunnerAddress.String(), &coins)
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
			executionHash = lcdBroadcastMsg(t, msg)
		})
		t.Run("in progress", func(t *testing.T) {
			execInProgress := pollExecution(t, executionHash, execution.Status_InProgress)
			require.Equal(t, executionHash, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, eventHash, execInProgress.EventHash)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("completed", func(t *testing.T) {
			exec = pollExecution(t, executionHash, execution.Status_Completed)
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
			var execR *execution.Execution
			lcdGet(t, "execution/get/"+executionHash.String(), &execR)
			require.True(t, exec.Equal(execR))
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
			msgs := make([]sdk.Msg, 0)
			for i := 0; i < n; i++ {
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
				msgs = append(msgs, msg)
			}
			execsHash := lcdBroadcastMsgs(t, msgs)
			// split hash
			hashSize := hash.DefaultHash().Size()
			for i := 0; i < n; i++ {
				execHash := execsHash[hashSize*i : hashSize*(i+1)]
				require.NotContains(t, executions, execHash)
				executions = append(executions, execHash)
			}
			require.Len(t, executions, n)
		})
		t.Run("check in progress", func(t *testing.T) {
			execs := make([]hash.Hash, 0)
			wg := sync.WaitGroup{}
			var mutex sync.Mutex
			wg.Add(n)
			for i := 0; i < n; i++ {
				go func(i int) {
					defer wg.Done()
					exec := pollExecution(t, executions[i], execution.Status_InProgress)
					mutex.Lock()
					defer mutex.Unlock()
					require.Contains(t, executions, exec.Hash)
					require.NotContains(t, execs, exec.Hash)
					execs = append(execs, exec.Hash)
				}(i)
			}
			wg.Wait()
			require.Len(t, execs, n)
		})
		t.Run("check completed", func(t *testing.T) {
			execs := make([]hash.Hash, 0)
			wg := sync.WaitGroup{}
			var mutex sync.Mutex
			wg.Add(n)
			for i := 0; i < n; i++ {
				go func(i int) {
					defer wg.Done()
					exec := pollExecution(t, executions[i], execution.Status_Completed)
					mutex.Lock()
					defer mutex.Unlock()
					require.Contains(t, executions, exec.Hash)
					require.NotContains(t, execs, exec.Hash)
					execs = append(execs, exec.Hash)
				}(i)
			}
			wg.Wait()
			require.Len(t, execs, n)
		})
	})
}
