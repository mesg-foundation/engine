package main

import (
	"context"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/x/ownership"
	processmodule "github.com/mesg-foundation/engine/x/process"
	"github.com/stretchr/testify/require"
)

func testOrchestratorProcessBalanceWithdraw(instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			processHash hash.Hash
			procAddress sdk.AccAddress
		)

		t.Run("create process", func(t *testing.T) {
			processHash = lcdBroadcastMsg(t, processmodule.MsgCreate{
				Owner: engineAddress,
				Name:  "balance-withdraw-process",
				Nodes: []*process.Process_Node{
					{
						Key: "n0",
						Type: &process.Process_Node_Event_{
							Event: &process.Process_Node_Event{
								InstanceHash: instanceHash,
								EventKey:     "test_event",
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
			})
		})
		t.Run("get process address", func(t *testing.T) {
			var proc *process.Process
			lcdGet(t, "process/get/"+processHash.String(), &proc)
			require.Equal(t, proc.Hash, processHash)
			procAddress = proc.Address
		})
		t.Run("check coins on process", func(t *testing.T) {
			var coins sdk.Coins
			lcdGet(t, "bank/balances/"+procAddress.String(), &coins)
			require.True(t, coins.IsEqual(processInitialBalance), coins)
		})
		t.Run("trigger process", func(t *testing.T) {
			_, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
				InstanceHash: instanceHash,
				Key:          "test_event",
				Data: &types.Struct{
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
			})
			require.NoError(t, err)
		})
		t.Run("check in progress execution", func(t *testing.T) {
			exec := pollExecutionOfProcess(t, processHash, execution.Status_InProgress, "n1")
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.Equal(t, "foo_1", exec.Inputs.Fields["msg"].GetStringValue())
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec := pollExecutionOfProcess(t, processHash, execution.Status_Completed, "n1")
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.Equal(t, "foo_1", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
		t.Run("check coins on process after 1 execution", func(t *testing.T) {
			var coins sdk.Coins
			lcdGet(t, "bank/balances/"+procAddress.String(), &coins)
			require.True(t, coins.IsEqual(processInitialBalance.Sub(minExecutionPrice)), coins)
		})
		t.Run("withdraw from process", func(t *testing.T) {
			coins := minExecutionPrice
			msg := ownership.MsgWithdraw{
				Owner:        engineAddress,
				Amount:       coins.String(),
				ResourceHash: processHash,
			}
			lcdBroadcastMsg(t, msg)

			lcdGet(t, "bank/balances/"+procAddress.String(), &coins)
			require.True(t, coins.IsEqual(processInitialBalance.Sub(minExecutionPrice).Sub(minExecutionPrice)), coins)
		})
		t.Run("delete process", func(t *testing.T) {
			lcdBroadcastMsg(t, processmodule.MsgDelete{
				Owner: engineAddress,
				Hash:  processHash,
			})
		})
		t.Run("check coins on process after deletion", func(t *testing.T) {
			var coins sdk.Coins
			lcdGet(t, "bank/balances/"+procAddress.String(), &coins)
			require.True(t, coins.IsZero(), coins)
		})
	}
}
