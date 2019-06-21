package executionsdk

import (
	"testing"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/execution"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var tests = []struct {
		f     *Filter
		e     *execution.Execution
		match bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			&Filter{},
			&execution.Execution{},
			true,
		},
		{
			&Filter{ServiceHash: []byte{0}},
			&execution.Execution{ServiceHash: base58.Encode([]byte{0})},
			true,
		},
		{
			&Filter{ServiceHash: []byte{0}},
			&execution.Execution{ServiceHash: base58.Encode([]byte{1})},
			false,
		},
		{
			&Filter{Statuses: []execution.Status{execution.Created}},
			&execution.Execution{Status: execution.Created},
			true,
		},
		{
			&Filter{Statuses: []execution.Status{execution.Created}},
			&execution.Execution{Status: execution.InProgress},
			false,
		},
		{
			&Filter{TaskKey: "0"},
			&execution.Execution{TaskKey: "0"},
			true,
		},
		{
			&Filter{TaskKey: "*"},
			&execution.Execution{TaskKey: "0"},
			true,
		},
		{
			&Filter{TaskKey: "0"},
			&execution.Execution{TaskKey: "1"},
			false,
		},
		{
			&Filter{Tags: []string{"0"}},
			&execution.Execution{Tags: []string{"0"}},
			true,
		},
		{
			&Filter{Tags: []string{"0", "1"}},
			&execution.Execution{Tags: []string{"0"}},
			false,
		},
	}

	for i, tt := range tests {
		assert.Equal(t, tt.match, tt.f.Match(tt.e), i)
	}
}

func TestListener(t *testing.T) {
	topic := "test-topic"
	testExecution := &execution.Execution{TaskKey: "0"}
	ps := pubsub.New(0)
	el := NewListener(ps, topic, &Filter{TaskKey: "0"})

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
