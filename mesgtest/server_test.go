package mesgtest

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/api/service"
	"github.com/stvp/assert"
)

const token = "token"

func TestNewServer(t *testing.T) {
	service := NewServer()
	assert.NotNil(t, service)
	assert.NotNil(t, service.Socket())
}

type eventRequest struct {
	URL string `json:"url"`
}

func TestLastEmit(t *testing.T) {
	var (
		key     = "key"
		data    = eventRequest{"https://mesg.tech"}
		dataStr = jsonMarshal(t, data)
	)

	server := NewServer()
	assert.NotNil(t, server)

	server.service.EmitEvent(context.Background(), &service.EmitEventRequest{
		EventKey:  key,
		EventData: dataStr,
		Token:     token,
	})

	le := <-server.LastEmit()

	assert.Equal(t, key, le.Name())
	assert.Equal(t, token, le.Token())

	var data1 eventRequest
	assert.Nil(t, le.Data(&data1))
	assert.Equal(t, data.URL, data1.URL)
}

func jsonMarshal(t *testing.T, data interface{}) string {
	bytes, err := json.Marshal(data)
	assert.Nil(t, err)
	return string(bytes)
}

type taskRequest struct {
	URL string `json:"url"`
}

type taskResponse struct {
	Message string `json:"message"`
}

func TestExecute(t *testing.T) {
	var (
		task       = "task"
		key        = "success"
		reqData    = taskRequest{"https://mesg.com"}
		resData    = taskResponse{"ok"}
		reqDataStr = jsonMarshal(t, reqData)
		resDataStr = jsonMarshal(t, resData)
	)

	server := NewServer()
	assert.NotNil(t, server)

	var executionID string
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		executionID1, execution, err := server.Execute(task, reqData)
		assert.Nil(t, err)
		assert.Equal(t, executionID, execution.ID())
		assert.Equal(t, executionID, executionID1)
		assert.Equal(t, key, execution.Key())

		var data taskResponse
		assert.Nil(t, execution.Data(&data))
		assert.Equal(t, resData.Message, data.Message)
	}()

	stream := newTaskDataStream()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.service.ListenTask(&service.ListenTaskRequest{Token: token}, stream)
		assert.Nil(t, err)
	}()

	taskData := <-stream.taskC
	executionID = taskData.ExecutionID
	assert.Equal(t, task, taskData.TaskKey)
	assert.Equal(t, reqDataStr, taskData.InputData)

	_, err := server.service.SubmitResult(context.Background(), &service.SubmitResultRequest{
		ExecutionID: executionID,
		OutputKey:   key,
		OutputData:  resDataStr,
	})
	assert.Nil(t, err)

	stream.close()
	wg.Wait()
}

func TestListenToken(t *testing.T) {
	server := NewServer()
	assert.NotNil(t, server)

	stream := newTaskDataStream()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.service.ListenTask(&service.ListenTaskRequest{Token: token}, stream)
		assert.Nil(t, err)
	}()

	stream.close()
	wg.Wait()

	assert.Equal(t, token, server.ListenToken())
}
