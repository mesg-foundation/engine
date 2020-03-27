package main

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	processmodule "github.com/mesg-foundation/engine/x/process"
	"github.com/stretchr/testify/require"
)

func testOrchestratorRefGrandParentTask(instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			processHash hash.Hash
		)
		t.Run("create process", func(t *testing.T) {
			processHash = lcdBroadcastMsg(t, processmodule.MsgCreate{
				Owner: engineAddress,
				Name:  "ref-grand-parent-task",
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
					{
						Key: "n2",
						Type: &process.Process_Node_Map_{
							Map: &process.Process_Node_Map{
								Outputs: map[string]*process.Process_Node_Map_Output{
									"msg": {
										Value: &process.Process_Node_Map_Output_StringConst{
											StringConst: "itsAConstant",
										},
									},
								},
							},
						},
					},
					{
						Key: "n3",
						Type: &process.Process_Node_Task_{
							Task: &process.Process_Node_Task{
								InstanceHash: instanceHash,
								TaskKey:      "task1",
							},
						},
					},
					{
						Key: "n4",
						Type: &process.Process_Node_Map_{
							Map: &process.Process_Node_Map{
								Outputs: map[string]*process.Process_Node_Map_Output{
									"msg": {
										Value: &process.Process_Node_Map_Output_Ref{
											Ref: &process.Process_Node_Reference{
												NodeKey: "n1",
												Path: &process.Process_Node_Reference_Path{
													Selector: &process.Process_Node_Reference_Path_Key{
														Key: "msg",
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Key: "n5",
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
					{Src: "n2", Dst: "n3"},
					{Src: "n3", Dst: "n4"},
					{Src: "n4", Dst: "n5"},
				},
			})
		})
		t.Run("trigger process", func(t *testing.T) {
			_, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
				InstanceHash: instanceHash,
				Key:          "test_event",
				Data: &types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StringValue{
								StringValue: "foo_event",
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
		t.Run("check first task", func(t *testing.T) {
			t.Run("check in progress execution", func(t *testing.T) {
				exec := pollExecutionOfProcess(t, processHash, execution.Status_InProgress, "n1")
				require.Equal(t, "n1", exec.NodeKey)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_InProgress, exec.Status)
				require.Equal(t, "foo_event", exec.Inputs.Fields["msg"].GetStringValue())
			})
			t.Run("check completed execution", func(t *testing.T) {
				exec := pollExecutionOfProcess(t, processHash, execution.Status_Completed, "n1")
				require.Equal(t, "n1", exec.NodeKey)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_Completed, exec.Status)
				require.Equal(t, "foo_event", exec.Outputs.Fields["msg"].GetStringValue())
				require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
			})
		})
		t.Run("check second task", func(t *testing.T) {
			t.Run("check in progress execution", func(t *testing.T) {
				exec := pollExecutionOfProcess(t, processHash, execution.Status_InProgress, "n3")
				require.Equal(t, "n3", exec.NodeKey)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_InProgress, exec.Status)
				require.Equal(t, "itsAConstant", exec.Inputs.Fields["msg"].GetStringValue())
			})
			t.Run("check completed execution", func(t *testing.T) {
				exec := pollExecutionOfProcess(t, processHash, execution.Status_Completed, "n3")
				require.Equal(t, "n3", exec.NodeKey)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_Completed, exec.Status)
				require.Equal(t, "itsAConstant", exec.Outputs.Fields["msg"].GetStringValue())
				require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
			})
		})
		t.Run("check third task", func(t *testing.T) {
			t.Run("check in progress execution", func(t *testing.T) {
				exec := pollExecutionOfProcess(t, processHash, execution.Status_InProgress, "n5")
				require.Equal(t, "n5", exec.NodeKey)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_InProgress, exec.Status)
				require.Equal(t, "foo_event", exec.Inputs.Fields["msg"].GetStringValue())
			})
			t.Run("check completed execution", func(t *testing.T) {
				exec := pollExecutionOfProcess(t, processHash, execution.Status_Completed, "n5")
				require.Equal(t, "n5", exec.NodeKey)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_Completed, exec.Status)
				require.Equal(t, "foo_event", exec.Outputs.Fields["msg"].GetStringValue())
				require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
			})
		})
		t.Run("delete process", func(t *testing.T) {
			lcdBroadcastMsg(t, processmodule.MsgDelete{
				Owner: engineAddress,
				Hash:  processHash,
			})
		})
	}
}
