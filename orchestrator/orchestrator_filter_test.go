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

func TestOrchestratorFilterTask(t *testing.T) {
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
			"filter",
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
							InstanceHash: testInstanceHash,
							TaskKey:      "task1",
						},
					},
				},
			},
			[]*process.Process_Edge{
				{Src: "n0", Dst: "n1"},
				{Src: "n1", Dst: "n2"},
			},
		)
		require.NoError(t, err)
	})
	t.Run("pass filter", func(t *testing.T) {
		t.Run("trigger process", func(t *testing.T) {
			_, err := ep.Publish(
				testInstanceHash,
				"event_trigger",
				&types.Struct{
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
			)
			require.NoError(t, err)
		})
		t.Run("check created execution", func(t *testing.T) {
			exec := <-execChan
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n2", exec.NodeKey)
			require.True(t, testProcessHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "shouldMatch", exec.Inputs.Fields["msg"].GetStringValue())
		})
	})
	t.Run("stop at filter", func(t *testing.T) {
		t.Run("trigger process", func(t *testing.T) {
			_, err := ep.Publish(
				testInstanceHash,
				"event_trigger",
				&types.Struct{
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
			)
			require.NoError(t, err)
		})
		t.Run("wait timeout to check execution is not created", func(t *testing.T) {
			select {
			case <-time.After(time.Second):
			case <-execChan:
				require.FailNow(t, "should timeout")
			}
		})
	})
}
