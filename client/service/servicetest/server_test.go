package servicetest

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"github.com/stretchr/testify/require"
)

const token = "token"

func TestNewServer(t *testing.T) {
	service := NewServer()
	require.NotNil(t, service)
	require.NotNil(t, service.Socket())
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
	require.NotNil(t, server)

	server.service.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		EventKey:  key,
		EventData: dataStr,
		Token:     token,
	})

	le := <-server.LastEmit()

	require.Equal(t, key, le.Name())
	require.Equal(t, token, le.Token())

	var data1 eventRequest
	require.Nil(t, le.Data(&data1))
	require.Equal(t, data.URL, data1.URL)
}

func jsonMarshal(t *testing.T, data interface{}) string {
	bytes, err := json.Marshal(data)
	require.NoError(t, err)
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
	require.NotNil(t, server)

	var executionID string
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		executionID1, execution, err := server.Execute(task, reqData)
		require.NoError(t, err)
		require.Equal(t, executionID, execution.ID())
		require.Equal(t, executionID, executionID1)
		require.Equal(t, key, execution.Key())

		var data taskResponse
		require.Nil(t, execution.Data(&data))
		require.Equal(t, resData.Message, data.Message)
	}()

	stream := newTaskDataStream()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.service.ListenTask(&serviceapi.ListenTaskRequest{Token: token}, stream)
		require.NoError(t, err)
	}()

	taskData := <-stream.taskC
	executionID = taskData.ExecutionID
	require.Equal(t, task, taskData.TaskKey)
	require.Equal(t, reqDataStr, taskData.InputData)

	_, err := server.service.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionID: executionID,
		OutputKey:   key,
		OutputData:  resDataStr,
	})
	require.NoError(t, err)

	stream.close()
	wg.Wait()
}

func TestListenToken(t *testing.T) {
	server := NewServer()
	require.NotNil(t, server)

	stream := newTaskDataStream()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.service.ListenTask(&serviceapi.ListenTaskRequest{Token: token}, stream)
		require.NoError(t, err)
	}()

	stream.close()
	wg.Wait()

	require.Equal(t, token, server.ListenToken())
}
