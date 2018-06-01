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
			"test": &service.Task{},
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
			"test": &service.Task{
				Outputs: map[string]*service.Output{
					"outputx": &service.Output{},
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

func TestIncludedIn(t *testing.T) {
	assert.False(t, includedIn([]string{}, ""))
	assert.True(t, includedIn([]string{""}, ""))
	assert.False(t, includedIn([]string{"a"}, ""))
	assert.True(t, includedIn([]string{"a"}, "a"))
	assert.False(t, includedIn([]string{""}, "a"))
	assert.True(t, includedIn([]string{"a", "b"}, "a"))
	assert.True(t, includedIn([]string{"a", "b"}, "b"))
	assert.False(t, includedIn([]string{"a", "b"}, "c"))
}

func TestIsSubscribedTask(t *testing.T) {
	x := &execution.Execution{Task: "task"}
	r := &ListenResultRequest{}
	assert.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskKey: ""}
	assert.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskKey: "*"}
	assert.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskKey: "task"}
	assert.True(t, isSubscribedTask(r, x))

	r = &ListenResultRequest{TaskKey: "xxx"}
	assert.False(t, isSubscribedTask(r, x))
}

func TestIsSubscribedOutput(t *testing.T) {
	x := &execution.Execution{Output: "output"}
	r := &ListenResultRequest{}
	assert.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputKey: ""}
	assert.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputKey: "*"}
	assert.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputKey: "output"}
	assert.True(t, isSubscribedOutput(r, x))

	r = &ListenResultRequest{OutputKey: "xxx"}
	assert.False(t, isSubscribedOutput(r, x))
}
