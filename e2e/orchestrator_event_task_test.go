package main

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	processmodule "github.com/mesg-foundation/engine/x/process"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func testOrchestratorEventTask(runnerHash hash.Hash, instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			processHash hash.Hash
			err         error
			execHash    hash.Hash
			logs        orchestrator.Orchestrator_LogsClient
		)

		t.Run("create process", func(t *testing.T) {
			msg := processmodule.MsgCreate{
				Owner: cliAddress,
				Name:  "event-task-process",
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
			require.True(t, processHash.Equal(log.Data.ProcessHash))
			require.Equal(t, "n0", log.Data.NodeKey)
			require.Equal(t, process.NodeTypeEvent, log.Data.NodeType)
			require.False(t, log.Data.EventHash.IsZero())
		})
		t.Run("check process creates execution", func(t *testing.T) {
			log, err := logs.Recv()
			require.NoError(t, err)
			require.True(t, processHash.Equal(log.Data.ProcessHash))
			require.Equal(t, "n1", log.Data.NodeKey)
			require.Equal(t, process.NodeTypeTask, log.Data.NodeType)
			require.False(t, log.Data.ExecutionHash.IsZero())
			execHash = log.Data.ExecutionHash
		})
		t.Run("check in progress execution", func(t *testing.T) {
			exec, err := pollExecutionOfProcess(processHash, execution.Status_InProgress, "n1")
			require.NoError(t, err)
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.True(t, execHash.Equal(exec.Hash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "foo_1", exec.Inputs.Fields["msg"].GetStringValue())
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec, err := pollExecutionOfProcess(processHash, execution.Status_Completed, "n1")
			require.NoError(t, err)
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.True(t, execHash.Equal(exec.Hash))
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "foo_1", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
		t.Run("delete process", func(t *testing.T) {
			_, err := lcd.BroadcastMsg(processmodule.MsgDelete{
				Owner: cliAddress,
				Hash:  processHash,
			})
			require.NoError(t, err)
		})
	}
}
