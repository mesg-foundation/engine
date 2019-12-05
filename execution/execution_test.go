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
		parentResultHash = hash.Int(1)
		eventHash        = hash.Int(2)
		instanceHash     = hash.Int(3)
		executorHash     = hash.Int(4)
		processHash      = hash.Int(5)
		nodeKey          = "test"
		taskKey          = "key"
		tags             = []string{"tag"}
		inputs           = &types.Struct{
			Fields: map[string]*types.Value{
				"test": {Kind: &types.Value_StringValue{StringValue: "hello"}},
			},
		}
	)
	exec := New(processHash, instanceHash, parentResultHash, eventHash, nodeKey, taskKey, inputs, tags, executorHash)
	require.NotNil(t, exec)
	require.True(t, exec.Hash.Valid())
	require.Equal(t, processHash, exec.ProcessHash)
	require.Equal(t, eventHash, exec.EventHash)
	require.Equal(t, instanceHash, exec.InstanceHash)
	require.Equal(t, parentResultHash, exec.ParentResultHash)
	require.Equal(t, inputs, exec.Inputs)
	require.Equal(t, taskKey, exec.TaskKey)
	require.Equal(t, nodeKey, exec.NodeKey)
	require.Equal(t, tags, exec.Tags)
	require.Equal(t, executorHash, exec.ExecutorHash)
}

func TestExecutionHash(t *testing.T) {
	ids := make(map[string]bool)
	f := func(instanceHash, parentResultHash, eventID []byte, taskKey string, tags []string) bool {
		e := New(nil, instanceHash, parentResultHash, eventID, "", taskKey, nil, tags, nil)
		if ids[string(e.Hash)] {
			return false
		}
		ids[string(e.Hash)] = true
		return true
	}
	require.NoError(t, quick.Check(f, nil))
}
