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

func TestIsSubscribedToTask(t *testing.T) {
	x := &execution.Execution{Task: "task"}
	r := &ListenResultRequest{}
	require.True(t, isSubscribedToTask(r, x))

	r = &ListenResultRequest{TaskFilter: ""}
	require.True(t, isSubscribedToTask(r, x))

	r = &ListenResultRequest{TaskFilter: "*"}
	require.True(t, isSubscribedToTask(r, x))

	r = &ListenResultRequest{TaskFilter: "task"}
	require.True(t, isSubscribedToTask(r, x))

	r = &ListenResultRequest{TaskFilter: "xxx"}
	require.False(t, isSubscribedToTask(r, x))
}

func TestIsSubscribedToOutput(t *testing.T) {
	x := &execution.Execution{Output: "output"}
	r := &ListenResultRequest{}
	require.True(t, isSubscribedToOutput(r, x))

	r = &ListenResultRequest{OutputFilter: ""}
	require.True(t, isSubscribedToOutput(r, x))

	r = &ListenResultRequest{OutputFilter: "*"}
	require.True(t, isSubscribedToOutput(r, x))

	r = &ListenResultRequest{OutputFilter: "output"}
	require.True(t, isSubscribedToOutput(r, x))

	r = &ListenResultRequest{OutputFilter: "xxx"}
	require.False(t, isSubscribedToOutput(r, x))
}

func TestIsSubscribedToTags(t *testing.T) {
	type result struct {
		execution *execution.Execution
		valid     bool
	}
	tests := []struct {
		request *ListenResultRequest
		results []result
	}{
		{
			&ListenResultRequest{},
			[]result{
				{&execution.Execution{}, true},
				{&execution.Execution{Tags: []string{"foo"}}, true},
				{&execution.Execution{Tags: []string{"foo", "bar"}}, true},
				{&execution.Execution{Tags: []string{"none"}}, true},
			},
		},
		{
			&ListenResultRequest{TagFilters: []string{"foo"}},
			[]result{
				{&execution.Execution{}, false},
				{&execution.Execution{Tags: []string{"foo"}}, true},
				{&execution.Execution{Tags: []string{"foo", "bar"}}, true},
				{&execution.Execution{Tags: []string{"none"}}, false},
			},
		},
		{
			&ListenResultRequest{TagFilters: []string{"foo", "bar"}},
			[]result{
				{&execution.Execution{}, false},
				{&execution.Execution{Tags: []string{"foo"}}, false},
				{&execution.Execution{Tags: []string{"foo", "bar"}}, true},
				{&execution.Execution{Tags: []string{"none"}}, false},
			},
		},
	}
	for _, test := range tests {
		for _, r := range test.results {
			require.Equal(t, r.valid, isSubscribedToTags(test.request, r.execution))
		}
	}
}
