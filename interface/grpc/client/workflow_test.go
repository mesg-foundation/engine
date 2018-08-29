package client

import (
	"testing"

	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/stretchr/testify/require"
)

func TestValidEventFromAny(t *testing.T) {
	wf := &Workflow{
		OnEvent: &Event{Name: "*"},
	}
	require.True(t, wf.validEvent(&core.EventData{EventKey: "xxx"}))
}

func TestValidEventFromValue(t *testing.T) {
	wf := &Workflow{
		OnEvent: &Event{Name: "xxx"},
	}
	require.True(t, wf.validEvent(&core.EventData{EventKey: "xxx"}))
	require.False(t, wf.validEvent(&core.EventData{EventKey: "yyy"}))
}

func TestValidResultFromAnyNameAndAnyOutput(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{Name: "*", Output: "*"},
	}
	require.True(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "xxx"}))
}

func TestValidResultFromAnyNameAndNotAnyOutput(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{Name: "*", Output: "xxx"},
	}
	require.True(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "xxx"}))
	require.False(t, wf.validResult(&core.ResultData{TaskKey: "yyy", OutputKey: "yyy"}))
}

func TestValidResultFromNotAnyNameAndAnyOutput(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{Name: "xxx", Output: "*"},
	}
	require.True(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "xxx"}))
	require.False(t, wf.validResult(&core.ResultData{TaskKey: "yyy", OutputKey: "yyy"}))
}

func TestValidResultFromNotAnyNameAndNotAnyOutput(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{Name: "xxx", Output: "yyy"},
	}
	require.True(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "yyy"}))
	require.False(t, wf.validResult(&core.ResultData{TaskKey: "yyy", OutputKey: "yyy"}))
	require.False(t, wf.validResult(&core.ResultData{TaskKey: "xxx", OutputKey: "xxx"}))
	require.False(t, wf.validResult(&core.ResultData{TaskKey: "yyy", OutputKey: "xxx"}))
}

func TestInvalidListenResult(t *testing.T) {
	wf := &Workflow{
		OnResult: &Result{},
	}
	require.NotNil(t, listenResults(wf))
}

func TestInvalidListenEvent(t *testing.T) {
	wf := &Workflow{
		OnEvent: &Event{},
	}
	require.NotNil(t, listenEvents(wf))
}

func TestInvalidWorkflowWithNoExecute(t *testing.T) {
	wf := Workflow{OnEvent: &Event{}}
	err := wf.Start()
	require.Equal(t, err.Error(), "A workflow needs a task")
}

func TestInvalidWorkflowWithNoEvent(t *testing.T) {
	wf := Workflow{Execute: &Task{}}
	err := wf.Start()
	require.Equal(t, err.Error(), "A workflow needs an event OnEvent or OnResult")
}
