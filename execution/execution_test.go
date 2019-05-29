package execution

import (
	"errors"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"
)

func TestNewFromService(t *testing.T) {
	var (
		hash       = "a"
		parentHash = "b"
		eventID    = "1"
		taskKey    = "key"
		tags       = []string{"tag"}
	)

	execution := New(hash, parentHash, eventID, taskKey, nil, tags)
	require.NotNil(t, execution)
	require.Equal(t, hash, execution.ServiceHash)
	require.Equal(t, parentHash, execution.ServiceParentHash)
	require.Equal(t, eventID, execution.EventID)
	require.Equal(t, taskKey, execution.TaskKey)
	require.Equal(t, map[string]interface{}(nil), execution.Inputs)
	require.Equal(t, tags, execution.Tags)
	require.Equal(t, Created, execution.Status)
}

func TestExecute(t *testing.T) {
	e := New("", "", "", "", nil, nil)
	require.NoError(t, e.Execute())
	require.Equal(t, InProgress, e.Status)
	require.Error(t, e.Execute())
}

func TestComplete(t *testing.T) {
	output := map[string]interface{}{"foo": "bar"}
	e := New("", "", "", "", nil, nil)

	e.Execute()
	require.NoError(t, e.Complete("test", output))
	require.Equal(t, Completed, e.Status)
	require.Equal(t, output, e.OutputData)
	require.Equal(t, "test", e.OutputKey)
	require.Error(t, e.Complete("", nil))
}

func TestFailed(t *testing.T) {
	err := errors.New("test")
	e := New("", "", "", "", nil, nil)
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

	f := func(service, parentService, eventID, taskKey, input string, tags []string) bool {
		e := New(service, parentService, eventID, taskKey, map[string]interface{}{"input": input}, tags)
		if ids[string(e.Hash)] {
			return false
		}
		ids[string(e.Hash)] = true
		return true
	}

	require.NoError(t, quick.Check(f, nil))
}
