package execution

import (
	"testing"
	"testing/quick"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var (
		parentHash   = hash.Int(1)
		eventHash    = hash.Int(2)
		instanceHash = hash.Int(3)
		executorHash = hash.Int(4)
		processHash  = hash.Int(5)
		stepID       = "test"
		taskKey      = "key"
		tags         = []string{"tag"}
		inputs       = &types.Struct{
			Fields: map[string]*types.Value{
				"test": &types.Value{Kind: &types.Value_StringValue{StringValue: "hello"}},
			},
		}
		execReq    *ExecutionRequest
		execRes    *ExecutionResult
		execResErr *ExecutionResult
	)
	t.Run("NewRequest", func(t *testing.T) {
		execReq = NewRequest(processHash, instanceHash, parentHash, eventHash, stepID, taskKey, inputs, tags, executorHash)
		require.NotNil(t, execReq)
		require.True(t, execReq.Hash.Valid())
		require.Equal(t, processHash, execReq.ProcessHash)
		require.Equal(t, eventHash, execReq.EventHash)
		require.Equal(t, instanceHash, execReq.InstanceHash)
		require.Equal(t, parentHash, execReq.ParentHash)
		require.Equal(t, inputs, execReq.Inputs)
		require.Equal(t, taskKey, execReq.TaskKey)
		require.Equal(t, stepID, execReq.StepID)
		require.Equal(t, tags, execReq.Tags)
		require.Equal(t, executorHash, execReq.ExecutorHash)
	})
	t.Run("NewResultWithOutputs", func(t *testing.T) {
		execRes = NewResultWithOutputs(execReq.Hash, &types.Struct{
			Fields: map[string]*types.Value{
				"test": &types.Value{Kind: &types.Value_StringValue{StringValue: "hello"}},
			},
		})
		require.True(t, execRes.Hash.Valid())
		require.True(t, execReq.Hash.Equal(execRes.RequestHash))
	})
	t.Run("NewResultWithError", func(t *testing.T) {
		execResErr = NewResultWithError(execReq.Hash, "error string")
		require.True(t, execResErr.Hash.Valid())
		require.True(t, execReq.Hash.Equal(execResErr.RequestHash))
	})
	t.Run("ToExecution", func(t *testing.T) {
		t.Run("InProgress", func(t *testing.T) {
			exec := ToExecution(execReq, nil)
			require.NotNil(t, execReq)
			require.True(t, exec.Hash.Valid())
			require.Equal(t, processHash, exec.ProcessHash)
			require.Equal(t, eventHash, exec.EventHash)
			require.Equal(t, instanceHash, exec.InstanceHash)
			require.Equal(t, parentHash, exec.ParentHash)
			require.Equal(t, inputs, exec.Inputs)
			require.Equal(t, taskKey, exec.TaskKey)
			require.Equal(t, stepID, exec.StepID)
			require.Equal(t, tags, exec.Tags)
			require.Equal(t, executorHash, exec.ExecutorHash)
			require.Nil(t, exec.Outputs)
			require.Empty(t, exec.Error)
			require.Equal(t, exec.Status, Status_InProgress)
		})
		t.Run("Completed", func(t *testing.T) {
			exec := ToExecution(execReq, execRes)
			require.NotNil(t, exec.Outputs)
			require.Empty(t, exec.Error)
			require.Equal(t, exec.Status, Status_Completed)
		})
		t.Run("Failed", func(t *testing.T) {
			exec := ToExecution(execReq, execResErr)
			require.Nil(t, exec.Outputs)
			require.Equal(t, "error string", exec.Error)
			require.Equal(t, exec.Status, Status_Failed)
		})
	})
}

func TestExecutionHash(t *testing.T) {
	ids := make(map[string]bool)

	f := func(instanceHash, parentHash, eventID []byte, taskKey string, tags []string) bool {
		e := NewRequest(nil, instanceHash, parentHash, eventID, "", taskKey, nil, tags, nil)
		if ids[string(e.Hash)] {
			return false
		}
		ids[string(e.Hash)] = true
		return true
	}

	require.NoError(t, quick.Check(f, nil))
}
