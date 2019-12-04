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
	"github.com/stretchr/testify/require"
)

func testOrchestratorEventMapTaskMapTask(executionStream pb.Execution_StreamClient, instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var processHash hash.Hash
		t.Skip("this test doesn't work as map cannot access the trigger event")
		t.Run("create process", func(t *testing.T) {
			respProc, err := client.ProcessClient.Create(context.Background(), &pb.CreateProcessRequest{
				Key: "result-map-task-map-task-process",
				Nodes: []*process.Process_Node{
					{
						Type: &process.Process_Node_Event_{
							Event: &process.Process_Node_Event{
								Key:          "n0",
								InstanceHash: instanceHash,
								EventKey:     "test_event",
							},
						},
					},
					{
						Type: &process.Process_Node_Map_{
							Map: &process.Process_Node_Map{
								Key: "n1",
								Outputs: []*process.Process_Node_Map_Output{
									{
										Key: "msg",
										Value: &process.Process_Node_Map_Output_Constant{
											Constant: &types.Value{Kind: &types.Value_StringValue{StringValue: "itsAConstant"}},
										},
									},
								},
							},
						},
					},
					{
						Type: &process.Process_Node_Task_{
							Task: &process.Process_Node_Task{
								Key:          "n2",
								InstanceHash: instanceHash,
								TaskKey:      "task1",
							},
						},
					},
					{
						Type: &process.Process_Node_Map_{
							Map: &process.Process_Node_Map{
								Key: "n3",
								Outputs: []*process.Process_Node_Map_Output{
									{
										Key: "msg",
										Value: &process.Process_Node_Map_Output_Ref{
											Ref: &process.Process_Node_Map_Output_Reference{
												NodeKey: "n0",
												Key:     "msg",
											},
										},
									},
								},
							},
						},
					},
					{
						Type: &process.Process_Node_Task_{
							Task: &process.Process_Node_Task{
								Key:          "n4",
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
				},
			})
			require.NoError(t, err)
			processHash = respProc.Hash
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
				exec, err := executionStream.Recv()
				require.NoError(t, err)
				require.Equal(t, "n2", exec.StepID)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_InProgress, exec.Status)
				require.Equal(t, "itsAConstant", exec.Inputs.Fields["msg"].GetStringValue())
			})
			t.Run("check completed execution", func(t *testing.T) {
				exec, err := executionStream.Recv()
				require.NoError(t, err)
				require.Equal(t, "n2", exec.StepID)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_Completed, exec.Status)
				require.Equal(t, "itsAConstant", exec.Outputs.Fields["msg"].GetStringValue())
				require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
			})
		})
		t.Run("check second task", func(t *testing.T) {
			t.Run("check in progress execution", func(t *testing.T) {
				exec, err := executionStream.Recv()
				require.NoError(t, err)
				require.Equal(t, "n4", exec.StepID)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_InProgress, exec.Status)
				require.Equal(t, "foo_event", exec.Inputs.Fields["msg"].GetStringValue())
			})
			t.Run("check completed execution", func(t *testing.T) {
				exec, err := executionStream.Recv()
				require.NoError(t, err)
				require.Equal(t, "n4", exec.StepID)
				require.Equal(t, "task1", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_Completed, exec.Status)
				require.Equal(t, "foo_event", exec.Outputs.Fields["msg"].GetStringValue())
				require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
			})
		})
		t.Run("delete process", func(t *testing.T) {
			_, err := client.ProcessClient.Delete(context.Background(), &pb.DeleteProcessRequest{Hash: processHash})
			require.NoError(t, err)
		})
	}
}
