package workflowvm

import (
	"testing"

	"github.com/mesg-foundation/core/workflow"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	vm := New()

	vm.Add(workflow.Workflow{
		Hash: "workflowA",
		Trigger: workflow.Trigger{
			InstanceHash: "instanceA",
			EventKey:     "eventA",
		},
		Tasks: []workflow.Task{
			{
				InstanceHash: "instanceB", // once this executed, it should both trig taskC & taskD.
				Key:          "taskB",
			},
			{
				InstanceHash: "instanceC",
				Key:          "taskC",
			},
		},
	})

	vm.Add(workflow.Workflow{
		Hash: "workflowB",
		Trigger: workflow.Trigger{
			InstanceHash: "instanceB",
			EventKey:     "executionFinished",
			Filter: workflow.Filter{
				TaskKey: "taskB",
			},
		},
		Tasks: []workflow.Task{
			{
				InstanceHash: "instanceD",
				Key:          "taskD",
			},
		},
	})

	go vm.Process(Event{
		InstanceHash: "instanceA",
		Key:          "eventA",
		Data:         map[string]interface{}{"dataA": 1},
	})
	req := <-vm.ExecuctionRequests
	require.Equal(t, map[string]interface{}{"dataA": 1}, req.Inputs)
	require.Equal(t, "instanceB", req.InstanceHash)
	require.Empty(t, req.ParentHash)
	require.Equal(t, "taskB", req.TaskKey)
	require.NotEmpty(t, req.Secret)

	go vm.Process(Event{
		InstanceHash: "instanceB",
		Key:          "executionFinished",
		TaskKey:      "taskB",
		Data:         map[string]interface{}{"dataB": 2},
		Secret:       req.Secret,
		ParentHash:   []byte{0xB},
	})
	req2 := <-vm.ExecuctionRequests
	require.Equal(t, map[string]interface{}{"dataB": 2}, req2.Inputs)
	require.Equal(t, "instanceC", req2.InstanceHash)
	require.Equal(t, []byte{0xB}, req2.ParentHash)
	require.Equal(t, "taskC", req2.TaskKey)
	require.NotEmpty(t, req2.Secret)
	require.NotEqual(t, req2.Secret, req.Secret)
	req3 := <-vm.ExecuctionRequests
	require.Equal(t, map[string]interface{}{"dataB": 2}, req3.Inputs)
	require.Equal(t, "instanceD", req3.InstanceHash)
	require.Equal(t, []byte{0xB}, req2.ParentHash)
	require.Equal(t, "taskD", req3.TaskKey)
	require.NotEmpty(t, req3.Secret)
	require.NotEqual(t, req3.Secret, req2.Secret)
}
