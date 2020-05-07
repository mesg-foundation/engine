package main

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	processmodule "github.com/mesg-foundation/engine/x/process"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func testOrchestratorLogs(runnerHash, instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			processHash     hash.Hash
			err             error
			logsAll         orchestrator.Orchestrator_LogsClient
			logsExact       orchestrator.Orchestrator_LogsClient
			logsWrong       orchestrator.Orchestrator_LogsClient
			cancelLogsWrong context.CancelFunc
		)

		t.Run("create process", func(t *testing.T) {
			msg := processmodule.MsgCreate{
				Owner: cliAddress,
				Name:  "logs",
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
		t.Run("init logs stream #1 - no filter", func(t *testing.T) {
			req := orchestrator.OrchestratorLogsRequest{}
			logsAll, err = client.OrchestratorClient.Logs(context.Background(), &req, grpc.PerRPCCredentials(&signCred{req}))
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(logsAll)
		})
		t.Run("init logs stream #2 - exact filter", func(t *testing.T) {
			req := orchestrator.OrchestratorLogsRequest{
				ProcessHashes: []hash.Hash{processHash},
			}
			logsExact, err = client.OrchestratorClient.Logs(context.Background(), &req, grpc.PerRPCCredentials(&signCred{req}))
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(logsExact)
		})
		t.Run("init logs stream #3 - wrong filter", func(t *testing.T) {
			h, err := hash.Random()
			require.NoError(t, err)
			req := orchestrator.OrchestratorLogsRequest{
				ProcessHashes: []hash.Hash{h},
			}
			var ctx context.Context
			ctx, cancelLogsWrong = context.WithCancel(context.Background())
			logsWrong, err = client.OrchestratorClient.Logs(ctx, &req, grpc.PerRPCCredentials(&signCred{req}))
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(logsWrong)
		})
		defer cancelLogsWrong()
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
		t.Run("exact filter", func(t *testing.T) {
			t.Run("check process is triggered", func(t *testing.T) {
				log, err := logsExact.Recv()
				require.NoError(t, err)
				require.True(t, processHash.Equal(log.ProcessHash))
				require.Equal(t, "n0", log.NodeKey)
				require.Equal(t, process.NodeTypeEvent, log.NodeType)
				require.False(t, log.EventHash.IsZero())
			})
			t.Run("check process creates execution", func(t *testing.T) {
				log, err := logsExact.Recv()
				require.NoError(t, err)
				require.True(t, processHash.Equal(log.ProcessHash))
				require.Equal(t, "n1", log.NodeKey)
				require.Equal(t, process.NodeTypeTask, log.NodeType)
				require.False(t, log.CreatedExecHash.IsZero())
			})
		})
		t.Run("no filter", func(t *testing.T) {
			t.Run("check process is triggered", func(t *testing.T) {
				log, err := logsAll.Recv()
				require.NoError(t, err)
				require.True(t, processHash.Equal(log.ProcessHash))
				require.Equal(t, "n0", log.NodeKey)
				require.Equal(t, process.NodeTypeEvent, log.NodeType)
				require.False(t, log.EventHash.IsZero())
			})
			t.Run("check process creates execution", func(t *testing.T) {
				log, err := logsAll.Recv()
				require.NoError(t, err)
				require.True(t, processHash.Equal(log.ProcessHash))
				require.Equal(t, "n1", log.NodeKey)
				require.Equal(t, process.NodeTypeTask, log.NodeType)
				require.False(t, log.CreatedExecHash.IsZero())
			})
		})
		t.Run("wrong filter - wait 1 sec", func(t *testing.T) {
			time.AfterFunc(1*time.Second, cancelLogsWrong)
			_, err := logsWrong.Recv()
			require.Error(t, err)
			require.Contains(t, err.Error(), "context canceled")
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
