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

func TestOrchestratorResultTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orch, _, store, _, testInstanceHash, testRunnerHash, _, execChan := newTestOrchestrator(ctx, t)
	defer orch.Stop()

	var (
		testProcessHash hash.Hash
		triggerExecHash hash.Hash
		err             error
	)
	t.Run("create process", func(t *testing.T) {
		testProcessHash, err = store.CreateProcess(
			context.Background(),
			"result-task-process",
			[]*process.Process_Node{
				{
					Key: "n0",
					Type: &process.Process_Node_Result_{
						Result: &process.Process_Node_Result{
							InstanceHash: testInstanceHash,
							TaskKey:      "task1",
						},
					},
				},
				{
					Key: "n1",
					Type: &process.Process_Node_Task_{
						Task: &process.Process_Node_Task{
							InstanceHash: testInstanceHash,
							TaskKey:      "task2",
						},
					},
				},
			}, []*process.Process_Edge{
				{Src: "n0", Dst: "n1"},
			},
		)
		require.NoError(t, err)
	})
	t.Run("trigger process", func(t *testing.T) {
		eventHash, err := hash.Random()
		require.NoError(t, err)
		triggerExecHash, err = store.CreateExecution(
			context.Background(),
			"task1",
			&types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StringValue{
							StringValue: "foo_2",
						},
					},
					"timestamp": {
						Kind: &types.Value_NumberValue{
							NumberValue: float64(time.Now().Unix()),
						},
					},
				},
			},
			nil,
			nil,
			eventHash,
			testProcessHash,
			"",
			testRunnerHash,
		)
		require.NoError(t, err)
	})
	t.Run("check trigger process execution", func(t *testing.T) {
		t.Run("in progress", func(t *testing.T) {
			exec := <-execChan
			require.NoError(t, err)
			require.Equal(t, triggerExecHash, exec.Hash)
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "", exec.NodeKey)
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.True(t, exec.Inputs.Equal(&types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StringValue{
							StringValue: "foo_2",
						},
					},
				},
			}))
		})
		t.Run("update exec", func(t *testing.T) {
			err := store.UpdateExecution(
				context.Background(),
				triggerExecHash,
				time.Now().UnixNano(),
				time.Now().UnixNano(),
				&types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StringValue{
								StringValue: "foo_2",
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
		t.Run("completed", func(t *testing.T) {
			exec := <-execChan
			require.NoError(t, err)
			require.Equal(t, triggerExecHash, exec.Hash)
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "", exec.NodeKey)
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "foo_2", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
	})
	t.Run("check created execution", func(t *testing.T) {
		exec := <-execChan
		require.NoError(t, err)
		require.Equal(t, "task2", exec.TaskKey)
		require.Equal(t, "n1", exec.NodeKey)
		require.Equal(t, testProcessHash, exec.ProcessHash)
		require.Equal(t, execution.Status_InProgress, exec.Status)
		require.Equal(t, "foo_2", exec.Inputs.Fields["msg"].GetStringValue())
	})
}
