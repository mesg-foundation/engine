package main

import (
	"context"
	"fmt"
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

func testOrchestratorFilter(runnerHash, instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			processHash hash.Hash
			err         error
			execHash    hash.Hash
			logs        orchestrator.Orchestrator_LogsClient
			cancelLogs  context.CancelFunc
		)

		t.Run("create process", func(t *testing.T) {
			msg := processmodule.MsgCreate{
				Owner: cliAddress,
				Name:  "filter",
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
						Type: &process.Process_Node_Filter_{
							Filter: &process.Process_Node_Filter{
								Conditions: []process.Process_Node_Filter_Condition{
									{
										Ref: &process.Process_Node_Reference{
											NodeKey: "n0",
											Path: &process.Process_Node_Reference_Path{
												Selector: &process.Process_Node_Reference_Path_Key{
													Key: "msg",
												},
											},
										},
										Predicate: process.Process_Node_Filter_Condition_EQ,
										Value: &types.Value{
											Kind: &types.Value_StringValue{
												StringValue: "shouldMatch",
											},
										},
									},
									{
										Ref: &process.Process_Node_Reference{
											NodeKey: "n0",
											Path: &process.Process_Node_Reference_Path{
												Selector: &process.Process_Node_Reference_Path_Key{
													Key: "timestamp",
												},
											},
										},
										Predicate: process.Process_Node_Filter_Condition_GT,
										Value: &types.Value{
											Kind: &types.Value_NumberValue{
												NumberValue: 10,
											},
										},
									},
								},
							},
						},
					},
					{
						Key: "n2",
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
					{Src: "n1", Dst: "n2"},
				},
			}
			processHash, err = lcd.BroadcastMsg(msg)
			require.NoError(t, err)
		})
		t.Run("init logs stream", func(t *testing.T) {
			req := orchestrator.OrchestratorLogsRequest{
				ProcessHashes: []hash.Hash{processHash},
			}
			var ctx context.Context
			ctx, cancelLogs = context.WithCancel(context.Background())
			logs, err = client.OrchestratorClient.Logs(ctx, &req, grpc.PerRPCCredentials(&signCred{req}))
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(logs)
		})
		t.Run("pass filter", func(t *testing.T) {
			t.Run("trigger process", func(t *testing.T) {
				req := orchestrator.ExecutionCreateRequest{
					Price:        "10000atto",
					TaskKey:      "task_trigger",
					ExecutorHash: runnerHash,
					Inputs: &types.Struct{
						Fields: map[string]*types.Value{
							"msg": {
								Kind: &types.Value_StringValue{
									StringValue: "shouldMatch",
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
			t.Run("check process executes filter", func(t *testing.T) {
				log, err := logs.Recv()
				require.NoError(t, err)
				require.True(t, processHash.Equal(log.Data.ProcessHash))
				require.Equal(t, "n1", log.Data.NodeKey)
				require.Equal(t, process.NodeTypeFilter, log.Data.NodeType)
			})
			t.Run("check process creates execution", func(t *testing.T) {
				log, err := logs.Recv()
				require.NoError(t, err)
				require.True(t, processHash.Equal(log.Data.ProcessHash))
				require.Equal(t, "n2", log.Data.NodeKey)
				require.Equal(t, process.NodeTypeTask, log.Data.NodeType)
				require.False(t, log.Data.ExecutionHash.IsZero())
				execHash = log.Data.ExecutionHash
			})
			t.Run("check in progress execution", func(t *testing.T) {
				exec, err := pollExecutionOfProcess(processHash, execution.Status_InProgress, "n2")
				require.NoError(t, err)
				require.Equal(t, "task1", exec.TaskKey)
				require.Equal(t, "n2", exec.NodeKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.True(t, execHash.Equal(exec.Hash))
				require.Equal(t, execution.Status_InProgress, exec.Status)
				require.Equal(t, "shouldMatch", exec.Inputs.Fields["msg"].GetStringValue())
			})
			t.Run("check completed execution", func(t *testing.T) {
				exec, err := pollExecutionOfProcess(processHash, execution.Status_Completed, "n2")
				require.NoError(t, err)
				require.Equal(t, "task1", exec.TaskKey)
				require.Equal(t, "n2", exec.NodeKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.True(t, execHash.Equal(exec.Hash))
				require.Equal(t, execution.Status_Completed, exec.Status)
				require.Equal(t, "shouldMatch", exec.Outputs.Fields["msg"].GetStringValue())
				require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
			})
		})
		t.Run("stop at filter", func(t *testing.T) {
			t.Run("trigger process", func(t *testing.T) {
				req := orchestrator.ExecutionCreateRequest{
					Price:        "10000atto",
					TaskKey:      "task_trigger",
					ExecutorHash: runnerHash,
					Inputs: &types.Struct{
						Fields: map[string]*types.Value{
							"msg": {
								Kind: &types.Value_StringValue{
									StringValue: "shouldNotMatch",
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
			t.Run("check process executes filter", func(t *testing.T) {
				log, err := logs.Recv()
				require.NoError(t, err)
				require.True(t, processHash.Equal(log.Data.ProcessHash))
				require.Equal(t, "n1", log.Data.NodeKey)
				require.Equal(t, process.NodeTypeFilter, log.Data.NodeType)
			})
			t.Run("check process no execution is created - wait 1sec", func(t *testing.T) {
				time.AfterFunc(1*time.Second, cancelLogs)
				_, err := logs.Recv()
				require.Error(t, err)
				require.Contains(t, err.Error(), "context canceled")
			})
			t.Run("wait timeout to check execution is not created", func(t *testing.T) {
				_, err := pollExecutionOfProcess(processHash, execution.Status_InProgress, "n2")
				require.EqualError(t, err, fmt.Sprintf("pollExecutionOfProcess timeout with process hash %q and status %q and nodeKey %q", processHash, execution.Status_InProgress, "n2"))
			})
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
