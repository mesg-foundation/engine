package api

import (
	"testing"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestValidateTaskKey(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	ln := newResultListener(a)

	s, _ := service.FromService(&service.Service{
		Tasks: []*service.Task{
			{
				Key: "test",
			},
		},
	}, service.ContainerOption(a.container))

	ln.taskKey = ""
	require.Nil(t, ln.validateTaskKey(s))

	ln.taskKey = "*"
	require.Nil(t, ln.validateTaskKey(s))

	ln.taskKey = "test"
	require.Nil(t, ln.validateTaskKey(s))

	ln.taskKey = "xxx"
	require.NotNil(t, ln.validateTaskKey(s))
}

func TestValidateOutputKey(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	ln := newResultListener(a)

	s, _ := service.FromService(&service.Service{
		Tasks: []*service.Task{
			{
				Key: "test",
				Outputs: []*service.Output{
					{
						Key: "outputx",
					},
				},
			},
		},
	}, service.ContainerOption(a.container))

	ln.taskKey = "test"
	ln.outputKey = ""
	require.Nil(t, ln.validateOutputKey(s))

	ln.taskKey = "test"
	ln.outputKey = "*"
	require.Nil(t, ln.validateOutputKey(s))

	ln.taskKey = "test"
	ln.outputKey = "outputx"
	require.Nil(t, ln.validateOutputKey(s))

	ln.taskKey = "test"
	ln.outputKey = "xxx"
	require.NotNil(t, ln.validateOutputKey(s))

	ln.taskKey = "xxx"
	ln.outputKey = ""
	require.Nil(t, ln.validateOutputKey(s))

	ln.taskKey = "xxx"
	ln.outputKey = "*"
	require.Nil(t, ln.validateOutputKey(s))

	ln.taskKey = "xxx"
	ln.outputKey = "outputX"
	require.NotNil(t, ln.validateOutputKey(s))

	ln.taskKey = "xxx"
	ln.outputKey = "xxx"
	require.NotNil(t, ln.validateOutputKey(s))
}

func TestIsSubscribedToTask(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	ln := newResultListener(a)

	x := &execution.Execution{TaskKey: "task"}

	ln.taskKey = ""
	require.True(t, ln.isSubscribedToTask(x))

	ln.taskKey = "*"
	require.True(t, ln.isSubscribedToTask(x))

	ln.taskKey = "task"
	require.True(t, ln.isSubscribedToTask(x))

	ln.taskKey = "xxx"
	require.False(t, ln.isSubscribedToTask(x))
}

func TestIsSubscribedToOutput(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	ln := newResultListener(a)

	x := &execution.Execution{Output: "output"}

	ln.outputKey = ""
	require.True(t, ln.isSubscribedToOutput(x))

	ln.outputKey = "*"
	require.True(t, ln.isSubscribedToOutput(x))

	ln.outputKey = "output"
	require.True(t, ln.isSubscribedToOutput(x))

	ln.outputKey = "xxx"
	require.False(t, ln.isSubscribedToOutput(x))
}

func TestIsSubscribedToTags(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	ln := newResultListener(a)

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
			ln.tagFilters = test.tags
			require.Equal(t, r.valid, ln.isSubscribedToTags(r.execution))
		}
	}
}

func TestIsSubscribed(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	ln := newResultListener(a)

	type test struct {
		taskFilter, outputFilter string
		tagFilters               []string
		e                        execution.Execution
	}
	subscribeToOutput := func(x *test) *test {
		return x
	}
	notSubscribeToOutput := func(x *test) *test {
		x.outputFilter = "foo"
		return x
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
		{notSubscribeToTags(notSubscribeToTask(notSubscribeToOutput(&test{}))), false, "[]"},
		{notSubscribeToTags(notSubscribeToTask(subscribeToOutput(&test{}))), false, "[output]"},
		{notSubscribeToTags(subscribeToTask(notSubscribeToOutput(&test{}))), false, "[task]"},
		{notSubscribeToTags(subscribeToTask(subscribeToOutput(&test{}))), false, "[task, output]"},
		{subscribeToTags(notSubscribeToTask(notSubscribeToOutput(&test{}))), false, "[tags]"},
		{subscribeToTags(notSubscribeToTask(subscribeToOutput(&test{}))), false, "[tags, output]"},
		{subscribeToTags(subscribeToTask(notSubscribeToOutput(&test{}))), false, "[tags, task]"},
		{subscribeToTags(subscribeToTask(subscribeToOutput(&test{}))), true, "[tags, task, output]"},
	}
	for _, test := range tests {
		ln.taskKey = test.t.taskFilter
		ln.outputKey = test.t.outputFilter
		ln.tagFilters = test.t.tagFilters
		require.Equal(t, test.valid, ln.isSubscribed(&test.t.e), test.msg)
	}
}
