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

func TestOrchestratorNestedData(t *testing.T) {
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
			"nested-data",
			[]*process.Process_Node{
				{
					Key: "n0",
					Type: &process.Process_Node_Event_{
						Event: &process.Process_Node_Event{
							InstanceHash: testInstanceHash,
							EventKey:     "event_complex_trigger",
						},
					},
				},
				{
					Key: "n1",
					Type: &process.Process_Node_Task_{
						Task: &process.Process_Node_Task{
							InstanceHash: testInstanceHash,
							TaskKey:      "task_complex",
						},
					},
				},
			},
			[]*process.Process_Edge{
				{Src: "n0", Dst: "n1"},
			},
		)
		require.NoError(t, err)
	})
	t.Run("trigger process", func(t *testing.T) {
		_, err := ep.Publish(
			testInstanceHash,
			"event_complex_trigger",
			&types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StructValue{
							StructValue: &types.Struct{
								Fields: map[string]*types.Value{
									"msg": {
										Kind: &types.Value_StringValue{
											StringValue: "complex",
										},
									},
									"timestamp": {
										Kind: &types.Value_NumberValue{
											NumberValue: float64(time.Now().Unix()),
										},
									},
									"array": {
										Kind: &types.Value_ListValue{
											ListValue: &types.ListValue{Values: []*types.Value{
												{Kind: &types.Value_StringValue{StringValue: "first"}},
												{Kind: &types.Value_StringValue{StringValue: "second"}},
												{Kind: &types.Value_StringValue{StringValue: "third"}},
											}},
										},
									},
								},
							},
						},
					},
				},
			},
		)
		require.NoError(t, err)
	})
	t.Run("check created execution", func(t *testing.T) {
		exec := <-execChan
		require.Equal(t, "task_complex", exec.TaskKey)
		require.Equal(t, "n1", exec.NodeKey)
		require.True(t, testProcessHash.Equal(exec.ProcessHash))
		require.Equal(t, execution.Status_InProgress, exec.Status)
		require.Equal(t, "complex", exec.Inputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
		require.Len(t, exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values, 3)
		require.Equal(t, "first", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[0].GetStringValue())
		require.Equal(t, "second", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[1].GetStringValue())
		require.Equal(t, "third", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[2].GetStringValue())
	})
}
