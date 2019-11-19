package executionsdk

import (
	"testing"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	"gotest.tools/assert"
)

func TestFilter(t *testing.T) {
	var tests = []struct {
		f     *api.StreamExecutionRequest_Filter
		e     *execution.Execution
		match bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			&api.StreamExecutionRequest_Filter{},
			&execution.Execution{},
			true,
		},
		{
			&api.StreamExecutionRequest_Filter{InstanceHash: hash.Int(1)},
			&execution.Execution{InstanceHash: hash.Int(1)},
			true,
		},
		{
			&api.StreamExecutionRequest_Filter{InstanceHash: hash.Int(1)},
			&execution.Execution{InstanceHash: hash.Int(2)},
			false,
		},
		{
			&api.StreamExecutionRequest_Filter{ExecutorHash: hash.Int(1)},
			&execution.Execution{ExecutorHash: hash.Int(1)},
			true,
		},
		{
			&api.StreamExecutionRequest_Filter{ExecutorHash: hash.Int(1)},
			&execution.Execution{ExecutorHash: hash.Int(2)},
			false,
		},
		{
			&api.StreamExecutionRequest_Filter{Statuses: []execution.Status{execution.Status_Created}},
			&execution.Execution{Status: execution.Status_Created},
			true,
		},
		{
			&api.StreamExecutionRequest_Filter{Statuses: []execution.Status{execution.Status_Created}},
			&execution.Execution{Status: execution.Status_InProgress},
			false,
		},
		{
			&api.StreamExecutionRequest_Filter{TaskKey: "0"},
			&execution.Execution{TaskKey: "0"},
			true,
		},
		{
			&api.StreamExecutionRequest_Filter{TaskKey: "*"},
			&execution.Execution{TaskKey: "0"},
			true,
		},
		{
			&api.StreamExecutionRequest_Filter{TaskKey: "0"},
			&execution.Execution{TaskKey: "1"},
			false,
		},
		{
			&api.StreamExecutionRequest_Filter{Tags: []string{"0"}},
			&execution.Execution{Tags: []string{"0"}},
			true,
		},
		{
			&api.StreamExecutionRequest_Filter{Tags: []string{"0", "1"}},
			&execution.Execution{Tags: []string{"0"}},
			false,
		},
	}

	for i, tt := range tests {
		assert.Equal(t, tt.match, match(tt.f, tt.e), i)
	}
}

func TestValidateFilter(t *testing.T) {
	var tests = []struct {
		f       *api.StreamExecutionRequest_Filter
		isError bool
	}{
		{
			nil,
			false,
		},
		{
			&api.StreamExecutionRequest_Filter{},
			false,
		},
		{
			&api.StreamExecutionRequest_Filter{
				TaskKey: "not-exist",
			},
			false,
		},
	}
	s := SDK{}
	for i, tt := range tests {
		if tt.isError {
			assert.Error(t, s.validateFilter(tt.f), "", i)
		} else {
			assert.NilError(t, s.validateFilter(tt.f), i)
		}
	}
}
