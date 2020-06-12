package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestOrchestratorRefGrandParentTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orch, ep, store, _, testInstanceHash, _, _, execChan := newTestOrchestrator(ctx, t)
	defer orch.Stop()

	var (
		testProcessHash hash.Hash
		err             error
	)
	t.Run("create process", func(t *testing.T) {
		testProcessHash, err = store.CreateProcess(
			context.Background(),
			"nested-map",
			[]*process.Process_Node{
				{
					Key: "n0",
					Type: &process.Process_Node_Event_{
						Event: &process.Process_Node_Event{
							InstanceHash: testInstanceHash,
							EventKey:     "event_trigger",
						},
					},
				},
				{
					Key: "n1",
					Type: &process.Process_Node_Task_{
						Task: &process.Process_Node_Task{
							InstanceHash: testInstanceHash,
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
							InstanceHash: testInstanceHash,
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
							InstanceHash: testInstanceHash,
							TaskKey:      "task1",
						},
					},
				},
			},
			[]*process.Process_Edge{
				{Src: "n0", Dst: "n1"},
				{Src: "n1", Dst: "n2"},
				{Src: "n2", Dst: "n3"},
				{Src: "n3", Dst: "n4"},
				{Src: "n4", Dst: "n5"},
			},
		)
		require.NoError(t, err)
	})
	t.Run("trigger process", func(t *testing.T) {
		_, err := ep.Publish(
			testInstanceHash,
			"event_trigger",
			&types.Struct{
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
		)
		require.NoError(t, err)
	})
	t.Run("check first task", func(t *testing.T) {
		var exec *execution.Execution
		t.Run("check created execution", func(t *testing.T) {
			exec = <-execChan
			require.Equal(t, "n1", exec.NodeKey)
			require.Equal(t, "task1", exec.TaskKey)
			require.True(t, testProcessHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "foo_event", exec.Inputs.Fields["msg"].GetStringValue())
		})
		t.Run("update exec", func(t *testing.T) {
			err := store.UpdateExecution(
				context.Background(),
				exec.Hash,
				time.Now().UnixNano(),
				time.Now().UnixNano(),
				&types.Struct{
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
				"",
			)
			require.NoError(t, err)
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec := <-execChan
			require.Equal(t, "n1", exec.NodeKey)
			require.Equal(t, "task1", exec.TaskKey)
			require.True(t, testProcessHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "foo_event", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
	})
	t.Run("check second task", func(t *testing.T) {
		var exec *execution.Execution
		t.Run("check created execution", func(t *testing.T) {
			exec = <-execChan
			require.Equal(t, "n3", exec.NodeKey)
			require.Equal(t, "task1", exec.TaskKey)
			require.True(t, testProcessHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "itsAConstant", exec.Inputs.Fields["msg"].GetStringValue())
		})
		t.Run("update exec", func(t *testing.T) {
			err := store.UpdateExecution(
				context.Background(),
				exec.Hash,
				time.Now().UnixNano(),
				time.Now().UnixNano(),
				&types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StringValue{
								StringValue: "itsAConstant",
							},
						},
						"timestamp": {
							Kind: &types.Value_NumberValue{
								NumberValue: float64(time.Now().Unix()),
							},
						},
					},
				},
				"",
			)
			require.NoError(t, err)
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec := <-execChan
			require.Equal(t, "n3", exec.NodeKey)
			require.Equal(t, "task1", exec.TaskKey)
			require.True(t, testProcessHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "itsAConstant", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
	})
	t.Run("check third task", func(t *testing.T) {
		var exec *execution.Execution
		t.Run("check created execution", func(t *testing.T) {
			exec = <-execChan
			require.Equal(t, "n5", exec.NodeKey)
			require.Equal(t, "task1", exec.TaskKey)
			require.True(t, testProcessHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "foo_event", exec.Inputs.Fields["msg"].GetStringValue())
		})
	})
}
