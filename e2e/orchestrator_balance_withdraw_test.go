package main

import (
	"context"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	"github.com/mesg-foundation/engine/x/ownership"
	processmodule "github.com/mesg-foundation/engine/x/process"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func testOrchestratorProcessBalanceWithdraw(runnerHash, instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			processHash hash.Hash
			procAddress sdk.AccAddress
			err         error
			execHash    hash.Hash
			logs        orchestrator.Orchestrator_LogsClient
		)

		t.Run("create process", func(t *testing.T) {
			msg := processmodule.MsgCreate{
				Owner: cliAddress,
				Name:  "balance-withdraw-process",
				Nodes: []*process.Process_Node{
					{
						Key: "n0",
						Type: &process.Process_Node_Event_{
							Event: &process.Process_Node_Event{
								InstanceHash: instanceHash,
								EventKey:     "event_trigger",
							},
						},
					},
					{
						Key: "n1",
						Type: &process.Process_Node_Task_{
							Task: &process.Process_Node_Task{
								InstanceHash: instanceHash,
								TaskKey:      "task1",
							},
						},
					},
				},
				Edges: []*process.Process_Edge{
					{Src: "n0", Dst: "n1"},
				},
			}
			processHash, err = lcd.BroadcastMsg(msg)
			require.NoError(t, err)
		})
		t.Run("get process address", func(t *testing.T) {
			var proc *process.Process
			require.NoError(t, lcd.Get("process/get/"+processHash.String(), &proc))
			require.Equal(t, proc.Hash, processHash)
			procAddress = proc.Address
		})
		t.Run("check coins on process", func(t *testing.T) {
			var coins sdk.Coins
			require.NoError(t, lcd.Get("bank/balances/"+procAddress.String(), &coins))
			require.True(t, coins.IsEqual(processInitialBalance), coins)
		})
		t.Run("init logs stream", func(t *testing.T) {
			req := orchestrator.OrchestratorLogsRequest{
				ProcessHashes: []hash.Hash{processHash},
			}
			logs, err = client.OrchestratorClient.Logs(context.Background(), &req, grpc.PerRPCCredentials(&signCred{req}))
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(logs)
		})
		t.Run("trigger process", func(t *testing.T) {
			req := orchestrator.ExecutionCreateRequest{
				Price:        "10000atto",
				TaskKey:      "task_trigger",
				ExecutorHash: runnerHash,
				Inputs: &types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StringValue{
								StringValue: "foo_1",
							},
						},
						"timestamp": {
							Kind: &types.Value_NumberValue{
								NumberValue: float64(time.Now().Unix()),
							},
						},
					},
				},
			}
			_, err := client.ExecutionClient.Create(context.Background(), &req, grpc.PerRPCCredentials(&signCred{req}))
			require.NoError(t, err)
		})
		t.Run("check process is triggered", func(t *testing.T) {
			log, err := logs.Recv()
			require.NoError(t, err)
			require.True(t, processHash.Equal(log.ProcessHash))
			require.Equal(t, "n0", log.NodeKey)
			require.Equal(t, process.NodeTypeEvent, log.NodeType)
			require.False(t, log.EventHash.IsZero())
		})
		t.Run("check process creates execution", func(t *testing.T) {
			log, err := logs.Recv()
			require.NoError(t, err)
			require.True(t, processHash.Equal(log.ProcessHash))
			require.Equal(t, "n1", log.NodeKey)
			require.Equal(t, process.NodeTypeTask, log.NodeType)
			require.False(t, log.CreatedExecHash.IsZero())
			execHash = log.CreatedExecHash
		})
		t.Run("check in progress execution", func(t *testing.T) {
			exec, err := pollExecutionOfProcess(processHash, execution.Status_InProgress, "n1")
			require.NoError(t, err)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.True(t, execHash.Equal(exec.Hash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.Equal(t, "foo_1", exec.Inputs.Fields["msg"].GetStringValue())
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec, err := pollExecutionOfProcess(processHash, execution.Status_Completed, "n1")
			require.NoError(t, err)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.True(t, execHash.Equal(exec.Hash))
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.Equal(t, "foo_1", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
		t.Run("check coins on process after 1 execution", func(t *testing.T) {
			var coins sdk.Coins
			require.NoError(t, lcd.Get("bank/balances/"+procAddress.String(), &coins))
			require.True(t, coins.IsEqual(processInitialBalance.Sub(executionPrice)), coins)
		})
		t.Run("withdraw from process", func(t *testing.T) {
			coins := executionPrice
			msg := ownership.MsgWithdraw{
				Owner:        cliAddress,
				Amount:       coins.String(),
				ResourceHash: processHash,
			}
			_, err := lcd.BroadcastMsg(msg)
			require.NoError(t, err)

			require.NoError(t, lcd.Get("bank/balances/"+procAddress.String(), &coins))
			require.True(t, coins.IsEqual(processInitialBalance.Sub(executionPrice).Sub(executionPrice)), coins)
		})
		t.Run("delete process", func(t *testing.T) {
			_, err := lcd.BroadcastMsg(processmodule.MsgDelete{
				Owner: cliAddress,
				Hash:  processHash,
			})
			require.NoError(t, err)
		})
		t.Run("check coins on process after deletion", func(t *testing.T) {
			var coins sdk.Coins
			require.NoError(t, lcd.Get("bank/balances/"+procAddress.String(), &coins))
			require.True(t, coins.IsZero(), coins)
		})
	}
}
