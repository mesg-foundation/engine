package main

import (
	"context"
	"math"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func testExecution(t *testing.T) {
	var (
		executorHash    = testRunnerHash
		executorAddress = testRunnerAddress
		err             error
	)

	t.Run("simple execution with price", func(t *testing.T) {
		var (
			executionHash hash.Hash
			exec          *execution.Execution
			taskKey       = "task1"
			inputs        = &types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StringValue{
							StringValue: "test",
						},
					},
				},
			}
			executorBalance sdk.Int
		)

		require.NoError(t, lcd.Get("credit/get/"+executorAddress.String(), &executorBalance))
		require.Equal(t, sdk.NewInt(0), executorBalance)

		t.Run("create", func(t *testing.T) {
			req := orchestrator.ExecutionCreateRequest{
				TaskKey:      taskKey,
				ExecutorHash: executorHash,
				Inputs:       inputs,
			}
			resp, err := client.ExecutionClient.Create(context.Background(), &req, grpc.PerRPCCredentials(&signCred{req}))
			require.NoError(t, err)
			executionHash = resp.Hash
		})
		t.Run("get execution address", func(t *testing.T) {
			var exec *execution.Execution
			require.NoError(t, lcd.Get("execution/get/"+executionHash.String(), &exec))
			require.Equal(t, exec.Hash, executionHash)
		})
		t.Run("in progress", func(t *testing.T) {
			execInProgress, err := pollExecution(executionHash, execution.Status_InProgress)
			require.NoError(t, err)
			require.Equal(t, executionHash, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("completed", func(t *testing.T) {
			exec, err = pollExecution(executionHash, execution.Status_Completed)
			require.NoError(t, err)
			require.Equal(t, executionHash, exec.Hash)
			require.Equal(t, taskKey, exec.TaskKey)
			require.Equal(t, executorHash, exec.ExecutorHash)
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.True(t, inputs.Equal(exec.Inputs))
			require.Equal(t, "test", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
		t.Run("get", func(t *testing.T) {
			var execR *execution.Execution
			require.NoError(t, lcd.Get("execution/get/"+executionHash.String(), &execR))
			require.True(t, exec.Equal(execR), exec, execR)
		})
		t.Run("executor balance", func(t *testing.T) {
			var credits sdk.Int
			var execR *execution.Execution
			require.NoError(t, lcd.Get("execution/get/"+executionHash.String(), &execR))
			require.NoError(t, lcd.Get("credit/get/"+executorAddress.String(), &credits))

			duration := sdk.NewInt(int64(math.Ceil(float64((exec.Stop - exec.Start) / 1e9))))
			require.True(t, duration.GTE(sdk.ZeroInt()))
			datasize := sdk.NewInt(int64(execR.Inputs.XXX_Size() + execR.Outputs.XXX_Size()))
			perCall, ok1 := sdk.NewIntFromString(task1Price.PerCall)
			require.True(t, ok1)
			perSec, ok2 := sdk.NewIntFromString(task1Price.PerSec)
			require.True(t, ok2)
			perKB, ok3 := sdk.NewIntFromString(task1Price.PerKB)
			require.True(t, ok3)
			expected := perCall.Add(duration.Mul(perSec)).Add(datasize.Mul(perKB))
			require.Equal(t, expected, credits)
		})
	})

	t.Run("complex execution", func(t *testing.T) {
		var (
			executionHash hash.Hash
			exec          *execution.Execution
			taskKey       = "task_complex"
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
			req := orchestrator.ExecutionCreateRequest{
				TaskKey:      taskKey,
				ExecutorHash: executorHash,
				Inputs:       inputs,
			}
			resp, err := client.ExecutionClient.Create(context.Background(), &req, grpc.PerRPCCredentials(&signCred{req}))
			require.NoError(t, err)
			executionHash = resp.Hash
		})
		t.Run("in progress", func(t *testing.T) {
			execInProgress, err := pollExecution(executionHash, execution.Status_InProgress)
			require.NoError(t, err)
			require.Equal(t, executionHash, execInProgress.Hash)
			require.Equal(t, taskKey, execInProgress.TaskKey)
			require.Equal(t, executorHash, execInProgress.ExecutorHash)
			require.Equal(t, execution.Status_InProgress, execInProgress.Status)
			require.True(t, inputs.Equal(execInProgress.Inputs))
		})
		t.Run("completed", func(t *testing.T) {
			exec, err = pollExecution(executionHash, execution.Status_Completed)
			require.NoError(t, err)
			require.Equal(t, executionHash, exec.Hash)
			require.Equal(t, taskKey, exec.TaskKey)
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
			require.NoError(t, lcd.Get("execution/get/"+executionHash.String(), &execR))
			require.True(t, exec.Equal(execR))
		})
	})

	t.Run("list", func(t *testing.T) {
		execs := make([]*execution.Execution, 0)
		require.NoError(t, lcd.Get("execution/list", &execs))
		require.Len(t, execs, 2)
	})

	t.Run("many executions in parallel", func(t *testing.T) {
		var (
			n          = 10
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
			var mutex sync.Mutex
			wg.Add(n)
			for i := 0; i < n; i++ {
				go func() {
					defer wg.Done()
					req := orchestrator.ExecutionCreateRequest{
						TaskKey:      taskKey,
						ExecutorHash: executorHash,
						Inputs:       inputs,
					}
					resp, err := client.ExecutionClient.Create(context.Background(), &req, grpc.PerRPCCredentials(&signCred{req}))
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
			execs := make([]hash.Hash, 0)
			wg := sync.WaitGroup{}
			var mutex sync.Mutex
			wg.Add(n)
			for i := 0; i < n; i++ {
				go func(i int) {
					defer wg.Done()
					exec, err := pollExecution(executions[i], execution.Status_InProgress)
					require.NoError(t, err)
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
					exec, err := pollExecution(executions[i], execution.Status_Completed)
					require.NoError(t, err)
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
