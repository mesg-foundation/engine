package mesg

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/api/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const token = "token"
const endpoint = "endpoint"

type requestData struct {
	Message string `json:"message"`
}

type replyData struct {
	Message string `json:"message"`
}

var errClosedConn = errors.New("closed connection")

func TestListenTasksAndReply(t *testing.T) {
	taskKey := "task"
	executionID := "executionID"
	inputData := requestData{Message: "inputData"}
	outputData := replyData{Message: "outputData"}
	outputKey := "outputKey"
	inputDataBytes, err := json.Marshal(inputData)
	assert.Nil(t, err)
	inputDataStr := string(inputDataBytes)
	outputDataBytes, err := json.Marshal(outputData)
	assert.Nil(t, err)
	outputDataStr := string(outputDataBytes)

	srv, err := NewService(
		ServiceTokenOption(token),
		ServiceEndpointOption(endpoint),
	)
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	taskC := make(chan *service.TaskData, 0)
	submitC := make(chan *service.SubmitResultRequest, 0)
	srv.Client = &testClient{
		stream:  &taskDataStream{taskC: taskC},
		submitC: submitC,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		err = srv.ListenTasks(
			NewTask(taskKey, func(req *Request) {
				var data requestData
				assert.Nil(t, req.Get(&data))
				assert.Equal(t, inputData.Message, data.Message)
				assert.Equal(t, executionID, req.executionID)

				err := req.Reply(outputKey, outputData)
				assert.Nil(t, err)

				wg.Done()
			}),
		)
		assert.Equal(t, errClosedConn, err)
		wg.Done()
	}()

	taskC <- &service.TaskData{
		ExecutionID: executionID,
		TaskKey:     taskKey,
		InputData:   inputDataStr,
	}
	close(taskC)

	reply := <-submitC
	assert.Equal(t, outputKey, reply.OutputKey)
	assert.Equal(t, outputDataStr, reply.OutputData)
	assert.Equal(t, executionID, reply.ExecutionID)

	wg.Wait()
}

func TestListenMultipleTasks(t *testing.T) {
	taskKey := "task"
	taskKey1 := "task2"
	executionID := "executionID"
	inputData := requestData{Message: "inputData"}
	inputDataBytes, err := json.Marshal(inputData)
	assert.Nil(t, err)
	inputDataStr := string(inputDataBytes)

	srv, err := NewService(
		ServiceTokenOption(token),
		ServiceEndpointOption(endpoint),
	)
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	taskC := make(chan *service.TaskData, 0)
	srv.Client = &testClient{
		stream: &taskDataStream{taskC: taskC},
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		err = srv.ListenTasks(
			NewTask(taskKey, func(req *Request) {
				assert.NotNil(t, req)
				wg.Done()
			}),
			NewTask(taskKey1, func(req *Request) {
				assert.NotNil(t, req)
				wg.Done()
			}),
		)
		assert.Equal(t, errClosedConn, err)
		wg.Done()
	}()

	taskC <- &service.TaskData{
		ExecutionID: executionID,
		TaskKey:     taskKey,
		InputData:   inputDataStr,
	}
	taskC <- &service.TaskData{
		ExecutionID: executionID,
		TaskKey:     taskKey1,
		InputData:   inputDataStr,
	}
	close(taskC)

	wg.Wait()
}

type eventData struct {
	Message string `json:"message"`
}

func TestEmitEvent(t *testing.T) {
	name := "name"
	eventData := eventData{Message: "eventData"}
	eventDataBytes, err := json.Marshal(eventData)
	assert.Nil(t, err)
	eventDataStr := string(eventDataBytes)

	srv, err := NewService(
		ServiceTokenOption(token),
		ServiceEndpointOption(endpoint),
	)
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	emitC := make(chan *service.EmitEventRequest, 0)
	srv.Client = &testClient{
		emitC: emitC,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		assert.Nil(t, srv.EmitEvent(name, eventData))
		wg.Done()
	}()

	data := <-emitC
	assert.Equal(t, name, data.EventKey)
	assert.Equal(t, eventDataStr, data.EventData)
	assert.Equal(t, token, data.Token)

	wg.Wait()
}

type testClient struct {
	stream  service.Service_ListenTaskClient
	emitC   chan *service.EmitEventRequest
	submitC chan *service.SubmitResultRequest
}

func (t *testClient) EmitEvent(ctx context.Context, in *service.EmitEventRequest,
	opts ...grpc.CallOption) (*service.EmitEventReply, error) {
	t.emitC <- in
	return nil, nil
}

func (t *testClient) ListenTask(ctx context.Context,
	in *service.ListenTaskRequest,
	opts ...grpc.CallOption) (service.Service_ListenTaskClient, error) {
	return t.stream, nil
}

func (t *testClient) SubmitResult(ctx context.Context,
	in *service.SubmitResultRequest,
	opts ...grpc.CallOption) (*service.SubmitResultReply, error) {
	t.submitC <- in
	return nil, nil
}

type taskDataStream struct {
	taskC chan *service.TaskData
	grpc.ClientStream
}

func (s taskDataStream) Recv() (*service.TaskData, error) {
	data, ok := <-s.taskC
	if !ok {
		return nil, errClosedConn
	}
	return data, nil
}
