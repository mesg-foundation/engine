package api

import (
	"fmt"
	"testing"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestValidateTaskKey(t *testing.T) {
	var (
		l = &ResultListener{}
		s = &service.Service{
			Tasks: []*service.Task{
				{
					Key: "test",
				},
			},
		}
	)

	l.taskKey = ""
	require.Nil(t, l.validateTaskKey(s))

	l.taskKey = "*"
	require.Nil(t, l.validateTaskKey(s))

	l.taskKey = "test"
	require.Nil(t, l.validateTaskKey(s))

	l.taskKey = "xxx"
	require.NotNil(t, l.validateTaskKey(s))
}

func TestIsSubscribedToTask(t *testing.T) {
	var (
		l = &ResultListener{}
		x = &execution.Execution{TaskKey: "task"}
	)

	l.taskKey = ""
	require.True(t, l.isSubscribedToTask(x))

	l.taskKey = "*"
	require.True(t, l.isSubscribedToTask(x))

	l.taskKey = "task"
	require.True(t, l.isSubscribedToTask(x))

	l.taskKey = "xxx"
	require.False(t, l.isSubscribedToTask(x))
}

func TestIsSubscribedToTags(t *testing.T) {
	l := &ResultListener{}

	type result struct {
		execution *execution.Execution
		valid     bool
	}
	tests := []struct {
		tags    []string
		results []result
	}{
		{
			[]string{},
			[]result{
				{&execution.Execution{}, true},
				{&execution.Execution{Tags: []string{"foo"}}, true},
				{&execution.Execution{Tags: []string{"foo", "bar"}}, true},
				{&execution.Execution{Tags: []string{"none"}}, true},
			},
		},
		{
			[]string{"foo"},
			[]result{
				{&execution.Execution{}, false},
				{&execution.Execution{Tags: []string{"foo"}}, true},
				{&execution.Execution{Tags: []string{"foo", "bar"}}, true},
				{&execution.Execution{Tags: []string{"none"}}, false},
			},
		},
		{
			[]string{"foo", "bar"},
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
			l.tagFilters = test.tags
			require.Equal(t, r.valid, l.isSubscribedToTags(r.execution))
		}
	}
}

func TestIsSubscribed(t *testing.T) {
	type test struct {
		taskFilter string
		tagFilters []string
		e          execution.Execution
	}
	subscribeToTask := func(x *test) *test {
		return x
	}
	notSubscribeToTask := func(x *test) *test {
		x.taskFilter = "foo"
		return x
	}
	subscribeToTags := func(x *test) *test {
		return x
	}
	notSubscribeToTags := func(x *test) *test {
		x.tagFilters = []string{"foo"}
		return x
	}
	tests := []struct {
		t     *test
		valid bool
		msg   string
	}{
		{notSubscribeToTags(notSubscribeToTask(&test{})), false, "[]"},
		{notSubscribeToTags(subscribeToTask(&test{})), false, "[task]"},
		{subscribeToTags(notSubscribeToTask(&test{})), false, "[tags]"},
		{subscribeToTags(subscribeToTask(&test{})), true, "[tags, task]"},
	}
	for _, test := range tests {
		l := &ResultListener{
			taskKey:    test.t.taskFilter,
			tagFilters: test.t.tagFilters,
		}
		fmt.Println(l.taskKey, l.tagFilters)
		require.Equal(t, test.valid, l.isSubscribed(&test.t.e), test.msg)
	}
}
