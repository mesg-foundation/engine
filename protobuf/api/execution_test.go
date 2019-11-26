package api

import (
	"testing"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var tests = []struct {
		f     *StreamExecutionRequest_Filter
		e     *execution.Execution
		match bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			&StreamExecutionRequest_Filter{},
			&execution.Execution{},
			true,
		},
		{
			&StreamExecutionRequest_Filter{InstanceHash: hash.Int(1)},
			&execution.Execution{InstanceHash: hash.Int(1)},
			true,
		},
		{
			&StreamExecutionRequest_Filter{InstanceHash: hash.Int(1)},
			&execution.Execution{InstanceHash: hash.Int(2)},
			false,
		},
		{
			&StreamExecutionRequest_Filter{ExecutorHash: hash.Int(1)},
			&execution.Execution{ExecutorHash: hash.Int(1)},
			true,
		},
		{
			&StreamExecutionRequest_Filter{ExecutorHash: hash.Int(1)},
			&execution.Execution{ExecutorHash: hash.Int(2)},
			false,
		},
		{
			&StreamExecutionRequest_Filter{Statuses: []execution.Status{execution.Status_Created}},
			&execution.Execution{Status: execution.Status_Created},
			true,
		},
		{
			&StreamExecutionRequest_Filter{Statuses: []execution.Status{execution.Status_Created}},
			&execution.Execution{Status: execution.Status_InProgress},
			false,
		},
		{
			&StreamExecutionRequest_Filter{TaskKey: "0"},
			&execution.Execution{TaskKey: "0"},
			true,
		},
		{
			&StreamExecutionRequest_Filter{TaskKey: "*"},
			&execution.Execution{TaskKey: "0"},
			true,
		},
		{
			&StreamExecutionRequest_Filter{TaskKey: "0"},
			&execution.Execution{TaskKey: "1"},
			false,
		},
		{
			&StreamExecutionRequest_Filter{Tags: []string{"0"}},
			&execution.Execution{Tags: []string{"0"}},
			true,
		},
		{
			&StreamExecutionRequest_Filter{Tags: []string{"0", "1"}},
			&execution.Execution{Tags: []string{"0"}},
			false,
		},
	}

	for i, tt := range tests {
		assert.Equal(t, tt.match, tt.f.Match(tt.e), i)
	}
}

func TestValidateFilter(t *testing.T) {
	var tests = []struct {
		f       *StreamExecutionRequest_Filter
		isError bool
	}{
		{
			nil,
			false,
		},
		{
			&StreamExecutionRequest_Filter{},
			false,
		},
		{
			&StreamExecutionRequest_Filter{
				TaskKey: "not-exist",
			},
			false,
		},
	}
	for i, tt := range tests {
		if tt.isError {
			assert.Error(t, tt.f.Validate(), "", i)
		} else {
			assert.NoError(t, tt.f.Validate(), i)
		}
	}
}
