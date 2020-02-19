package execution

import (
	"errors"
	"testing"
	"testing/quick"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestNewFromService(t *testing.T) {
	var (
		parentHash = hash.Int(2)
		eventHash  = hash.Int(3)
		hash       = hash.Int(1)
		taskKey    = "key"
		tags       = []string{"tag"}
	)

	execution := New(nil, hash, parentHash, eventHash, "", taskKey, "", nil, tags, nil)
	require.NotNil(t, execution)
	require.Equal(t, hash, execution.InstanceHash)
	require.Equal(t, parentHash, execution.ParentHash)
	require.Equal(t, eventHash, execution.EventHash)
	require.Equal(t, taskKey, execution.TaskKey)
	require.Equal(t, (*types.Struct)(nil), execution.Inputs)
	require.Equal(t, tags, execution.Tags)
	require.Equal(t, Status_Created, execution.Status)
}

func TestExecute(t *testing.T) {
	e := New(nil, nil, nil, nil, "", "", "", nil, nil, nil)
	require.NoError(t, e.Execute())
	require.Equal(t, Status_InProgress, e.Status)
	require.Error(t, e.Execute())
}

func TestComplete(t *testing.T) {
	var output types.Struct
	e := New(nil, nil, nil, nil, "", "", "", nil, nil, nil)
	e.Execute()
	require.NoError(t, e.Complete(&output))
	require.Equal(t, Status_Completed, e.Status)
	require.Equal(t, &output, e.Outputs)
	require.Error(t, e.Complete(nil))
}

func TestFailed(t *testing.T) {
	err := errors.New("test")
	e := New(nil, nil, nil, nil, "", "", "", nil, nil, nil)
	e.Execute()
	require.NoError(t, e.Failed(err))
	require.Equal(t, Status_Failed, e.Status)
	require.Equal(t, err.Error(), e.Error)
	require.Error(t, e.Failed(err))
}

func TestExecutionHash(t *testing.T) {
	ids := make(map[string]bool)

	f := func(instanceHash, parentHash, eventID []byte, taskKey string, tags []string) bool {
		e := New(nil, instanceHash, parentHash, eventID, "", taskKey, "", nil, tags, nil)
		if ids[string(e.Hash)] {
			return false
		}
		ids[string(e.Hash)] = true
		return true
	}

	require.NoError(t, quick.Check(f, nil))
}
