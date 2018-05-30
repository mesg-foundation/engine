package client

import (
	"testing"

	"github.com/mesg-foundation/core/api/core"
	"github.com/stvp/assert"
)

func TestProcessEventWithInvalidEventData(t *testing.T) {
	wf := &Workflow{
		Execute: &Task{
			Name:      "TestProcessEventWithInvalidEventData",
			ServiceID: "xxx",
		},
	}
	data := &core.EventData{
		EventKey:  "EventX",
		EventData: "",
	}
	err := wf.Execute.processEvent(wf, data)
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestProcessResulsWithInvalidEventData(t *testing.T) {
	wf := &Workflow{
		Execute: &Task{
			Name:      "TestProcessResulsWithInvalidEventData",
			ServiceID: "xxx",
		},
	}
	data := &core.ResultData{
		ExecutionID: "xxx",
		OutputData:  "",
		OutputKey:   "outputx",
		TaskKey:     "taskx",
	}
	err := wf.Execute.processResult(wf, data)
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestConvertData(t *testing.T) {
	task := &Task{
		Inputs: func(interface{}) interface{} {
			return "bar"
		},
	}
	res, err := task.convertData("foo")
	assert.Nil(t, err)
	assert.Equal(t, res, "\"bar\"")
}

func TestConvertDataObject(t *testing.T) {
	task := &Task{
		Inputs: func(d interface{}) interface{} {
			return d
		},
	}
	res, err := task.convertData(map[string]interface{}{
		"foo":    "bar",
		"number": 42,
	})
	assert.Nil(t, err)
	assert.Equal(t, res, "{\"foo\":\"bar\",\"number\":42}")
}

func TestConvertDataWithNull(t *testing.T) {
	task := &Task{
		Inputs: func(d interface{}) interface{} {
			return nil
		},
	}
	res, err := task.convertData("xxx")
	assert.Nil(t, err)
	assert.Equal(t, res, "null")
}
