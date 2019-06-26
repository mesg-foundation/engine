package execution

import (
	"errors"
	"testing"
	"testing/quick"

	"github.com/mesg-foundation/core/hash"
	"github.com/stretchr/testify/require"
)

func TestNewFromService(t *testing.T) {
	var (
		parentHash = hash.Int(2)
		hash       = hash.Int(1)
		eventID    = "1"
		taskKey    = "key"
		tags       = []string{"tag"}
	)

	execution := New(hash, parentHash, eventID, taskKey, nil, tags)
	require.NotNil(t, execution)
	require.Equal(t, hash, execution.InstanceHash)
	require.Equal(t, parentHash, execution.ParentHash)
	require.Equal(t, eventID, execution.EventID)
	require.Equal(t, taskKey, execution.TaskKey)
	require.Equal(t, map[string]interface{}(nil), execution.Inputs)
	require.Equal(t, tags, execution.Tags)
	require.Equal(t, Created, execution.Status)
}

func TestExecute(t *testing.T) {
	e := New(nil, nil, "", "", nil, nil)
	require.NoError(t, e.Execute())
	require.Equal(t, InProgress, e.Status)
	require.Error(t, e.Execute())
}

func TestComplete(t *testing.T) {
	output := map[string]interface{}{"foo": "bar"}
	e := New(nil, nil, "", "", nil, nil)

	e.Execute()
	require.NoError(t, e.Complete(output))
	require.Equal(t, Completed, e.Status)
	require.Equal(t, output, e.Outputs)
	require.Error(t, e.Complete(nil))
}

func TestFailed(t *testing.T) {
	err := errors.New("test")
	e := New(nil, nil, "", "", nil, nil)
	e.Execute()
	require.NoError(t, e.Failed(err))
	require.Equal(t, Failed, e.Status)
	require.Equal(t, err.Error(), e.Error)
	require.Error(t, e.Failed(err))
}

func TestStatus(t *testing.T) {
	require.Equal(t, "created", Created.String())
	require.Equal(t, "in progress", InProgress.String())
	require.Equal(t, "completed", Completed.String())
	require.Equal(t, "failed", Failed.String())
}

func TestExecutionHash(t *testing.T) {
	ids := make(map[string]bool)

	f := func(instanceHash, parentHash []byte, eventID, taskKey, input string, tags []string) bool {
		e := New(instanceHash, parentHash, eventID, taskKey, map[string]interface{}{"input": input}, tags)
		if ids[string(e.Hash)] {
			return false
		}
		ids[string(e.Hash)] = true
		return true
	}

	require.NoError(t, quick.Check(f, nil))
}
