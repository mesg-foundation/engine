package api

import (
	"testing"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/execution"
	"github.com/stretchr/testify/assert"
)

func TestExecutionFilter(t *testing.T) {
	var tests = []struct {
		f     *ExecutionFilter
		e     *execution.Execution
		match bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			&ExecutionFilter{},
			&execution.Execution{},
			true,
		},
		{
			&ExecutionFilter{Status: execution.Created},
			&execution.Execution{Status: execution.Created},
			true,
		},
		{
			&ExecutionFilter{Status: execution.Created},
			&execution.Execution{Status: execution.InProgress},
			false,
		},
		{
			&ExecutionFilter{TaskKey: "0"},
			&execution.Execution{TaskKey: "0"},
			true,
		},
		{
			&ExecutionFilter{TaskKey: "*"},
			&execution.Execution{TaskKey: "0"},
			true,
		},
		{
			&ExecutionFilter{TaskKey: "0"},
			&execution.Execution{TaskKey: "1"},
			false,
		},
		{
			&ExecutionFilter{OutputKey: "0"},
			&execution.Execution{OutputKey: "0"},
			true,
		},
		{
			&ExecutionFilter{OutputKey: "*"},
			&execution.Execution{OutputKey: "0"},
			true,
		},
		{
			&ExecutionFilter{OutputKey: "0"},
			&execution.Execution{OutputKey: "1"},
			false,
		},
		{
			&ExecutionFilter{Tags: []string{"0"}},
			&execution.Execution{Tags: []string{"0"}},
			true,
		},
		{
			&ExecutionFilter{Tags: []string{"0", "1"}},
			&execution.Execution{Tags: []string{"0"}},
			false,
		},
	}

	for i, tt := range tests {
		assert.Equal(t, tt.match, tt.f.Match(tt.e), i)
	}
}

func TestExecutionListener(t *testing.T) {
	topic := "test-topic"
	testExecution := &execution.Execution{TaskKey: "0"}
	ps := pubsub.New(0)
	el := NewExecutionListener(ps, topic, &ExecutionFilter{TaskKey: "0"})

	go func() {
		ps.Pub(&execution.Execution{TaskKey: "1"}, topic)
		ps.Pub(testExecution, topic)
	}()
	go el.Listen()

	recvExecution := <-el.C
	assert.Equal(t, testExecution, recvExecution)

	el.Close()
	_, ok := <-el.C
	assert.False(t, ok)
}
