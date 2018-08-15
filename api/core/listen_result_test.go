package core

import (
	"testing"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestValidateTaskKey(t *testing.T) {
	s := &service.Service{
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	require.Nil(t, validateTaskKey(s, ""))
	require.Nil(t, validateTaskKey(s, "*"))
	require.Nil(t, validateTaskKey(s, "test"))
	require.NotNil(t, validateTaskKey(s, "xxx"))
}

func TestValidateOutputKey(t *testing.T) {
	s := &service.Service{
		Tasks: map[string]*service.Task{
			"test": {
				Outputs: map[string]*service.Output{
					"outputx": {},
				},
			},
		},
	}
	require.Nil(t, validateOutputKey(s, "test", ""))
	require.Nil(t, validateOutputKey(s, "test", "*"))
	require.Nil(t, validateOutputKey(s, "test", "outputx"))
	require.NotNil(t, validateOutputKey(s, "test", "xxx"))
	require.Nil(t, validateOutputKey(s, "xxx", ""))
	require.Nil(t, validateOutputKey(s, "xxx", "*"))
	require.NotNil(t, validateOutputKey(s, "xxx", "outputX"))
	require.NotNil(t, validateOutputKey(s, "xxx", "xxx"))
}

func TestIsSubscribedTask(t *testing.T) {
	x := &execution.Execution{Task: "task"}
	r := &ListenResultRequest{}
	require.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskFilter: ""}
	require.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskFilter: "*"}
	require.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskFilter: "task"}
	require.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskFilter: "xxx"}
	require.False(t, isSubscribedTask(r, x))
}

func TestIsSubscribedOutput(t *testing.T) {
	x := &execution.Execution{Output: "output"}
	r := &ListenResultRequest{}
	require.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputFilter: ""}
	require.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputFilter: "*"}
	require.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputFilter: "output"}
	require.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputFilter: "xxx"}
	require.False(t, isSubscribedOutput(r, x))
}
