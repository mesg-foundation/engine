package core

import (
	"testing"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestValidateTaskKey(t *testing.T) {
	s := &service.Service{
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	assert.Nil(t, validateTaskKey(s, ""))
	assert.Nil(t, validateTaskKey(s, "*"))
	assert.Nil(t, validateTaskKey(s, "test"))
	assert.NotNil(t, validateTaskKey(s, "xxx"))
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
	assert.Nil(t, validateOutputKey(s, "test", ""))
	assert.Nil(t, validateOutputKey(s, "test", "*"))
	assert.Nil(t, validateOutputKey(s, "test", "outputx"))
	assert.NotNil(t, validateOutputKey(s, "test", "xxx"))
	assert.Nil(t, validateOutputKey(s, "xxx", ""))
	assert.Nil(t, validateOutputKey(s, "xxx", "*"))
	assert.NotNil(t, validateOutputKey(s, "xxx", "outputX"))
	assert.NotNil(t, validateOutputKey(s, "xxx", "xxx"))
}

func TestIsSubscribedTask(t *testing.T) {
	x := &execution.Execution{Task: "task"}
	r := &ListenResultRequest{}
	assert.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskFilter: ""}
	assert.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskFilter: "*"}
	assert.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskFilter: "task"}
	assert.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskFilter: "xxx"}
	assert.False(t, isSubscribedTask(r, x))
}

func TestIsSubscribedOutput(t *testing.T) {
	x := &execution.Execution{Output: "output"}
	r := &ListenResultRequest{}
	assert.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputFilter: ""}
	assert.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputFilter: "*"}
	assert.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputFilter: "output"}
	assert.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputFilter: "xxx"}
	assert.False(t, isSubscribedOutput(r, x))
}

func TestIsSubscribedExecution(t *testing.T) {
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
			&ListenResultRequest{TagFilter: []string{"foo"}},
			[]result{
				{&execution.Execution{}, false},
				{&execution.Execution{Tags: []string{"foo"}}, true},
				{&execution.Execution{Tags: []string{"foo", "bar"}}, true},
				{&execution.Execution{Tags: []string{"none"}}, false},
			},
		},
		{
			&ListenResultRequest{TagFilter: []string{"foo", "bar"}},
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
			assert.Equal(t, r.valid, isSubscribedExecution(test.request, r.execution))
		}
	}
}
