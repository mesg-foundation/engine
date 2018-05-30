package client

import (
	"testing"

	"github.com/mesg-foundation/core/api/core"
	"github.com/stvp/assert"
)

func TestValidEventFromAny(t *testing.T) {
	wf := &Workflow{
		OnEvent: &Event{Name: "*"},
	}
	assert.True(t, wf.validEvent(&core.EventData{EventKey: "xxx"}))
}

func TestValidEventFromValue(t *testing.T) {
	wf := &Workflow{
		OnEvent: &Event{Name: "xxx"},
	}
	assert.True(t, wf.validEvent(&core.EventData{EventKey: "xxx"}))
	assert.False(t, wf.validEvent(&core.EventData{EventKey: "yyy"}))
}

func TestValidResultFromAnyNameAndAnyOutput(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{Name: "*", Output: "*"},
	}
	assert.True(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "xxx"}))
}

func TestValidResultFromAnyNameAndNotAnyOutput(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{Name: "*", Output: "xxx"},
	}
	assert.True(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "xxx"}))
	assert.False(t, wf.validResult(&core.ResultData{TaskKey: "yyy", OutputKey: "yyy"}))
}

func TestValidResultFromNotAnyNameAndAnyOutput(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{Name: "xxx", Output: "*"},
	}
	assert.True(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "xxx"}))
	assert.False(t, wf.validResult(&core.ResultData{TaskKey: "yyy", OutputKey: "yyy"}))
}

func TestValidResultFromNotAnyNameAndNotAnyOutput(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{Name: "xxx", Output: "yyy"},
	}
	assert.True(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "yyy"}))
	assert.False(t, wf.validResult(&core.ResultData{TaskKey: "yyy", OutputKey: "yyy"}))
	assert.False(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "xxx"}))
	assert.False(t, wf.validResult(&core.ResultData{TaskKey: "yyy", OutputKey: "xxx"}))
}

func TestInvalidListenResult(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{},
	}
	assert.NotNil(t, listenResults(wf))
}

func TestInvalidListenEvent(t *testing.T) {
	wf := &Workflow{
		OnEvent: &Event{},
	}
	assert.NotNil(t, listenEvents(wf))
}

func TestInvalidWorkflowWithNoExecute(t *testing.T) {
	wf := Workflow{OnEvent: &Event{}}
	err := wf.Start()
	assert.Equal(t, err.Error(), "A workflow needs a task")
}

func TestInvalidWorkflowWithNoEvent(t *testing.T) {
	wf := Workflow{Execute: &Task{}}
	err := wf.Start()
	assert.Equal(t, err.Error(), "A workflow needs an event OnEvent or OnResult")
}
