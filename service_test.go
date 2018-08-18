package mesg

import (
	"bufio"
	"io"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/mesg-foundation/go-service/mesgtest"
	"github.com/stvp/assert"
)

const token = "token"
const endpoint = "endpoint"

type eventRequest struct {
	URL string `json:"url"`
}

func newServiceAndServer(t *testing.T) (*Service, *mesgtest.Server) {
	testServer := mesgtest.NewServer()

	service, err := New(
		DialOption(testServer.Socket()),
		TokenOption(token),
		EndpointOption(endpoint),
		LogOutputOption(ioutil.Discard),
	)

	assert.Nil(t, err)
	assert.NotNil(t, service)

	return service, testServer
}

func TestEmit(t *testing.T) {
	var (
		event = "request"
		data  = eventRequest{"https://mesg.tech"}
	)

	service, server := newServiceAndServer(t)
	go server.Start()

	go func() { assert.Nil(t, service.Emit(event, data)) }()
	le := <-server.LastEmit()

	assert.Equal(t, event, le.Name())
	assert.Equal(t, token, le.Token())

	var data1 eventRequest
	assert.Nil(t, le.Data(&data1))
	assert.Equal(t, data.URL, data1.URL)
}

type taskRequest struct {
	URL string `json:"url"`
}

type taskResponse struct {
	Message string `json:"message"`
}

func TestListen(t *testing.T) {
	var (
		task    = "send"
		key     = "success"
		reqData = taskRequest{"https://mesg.com"}
		resData = taskResponse{"ok"}
	)

	service, server := newServiceAndServer(t)
	go server.Start()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := service.Listen(
			Task(task, func(execution *Execution) (string, Data) {
				var data2 taskRequest
				assert.Nil(t, execution.Data(&data2))
				assert.Equal(t, reqData.URL, data2.URL)

				return key, resData
			}),
		)

		assert.NotNil(t, err)
	}()

	id, execution, err := server.Execute(task, reqData)
	assert.Nil(t, err)
	assert.Equal(t, id, execution.ID())
	assert.Equal(t, key, execution.Key())
	assert.Equal(t, token, server.ListenToken())

	var data1 taskResponse
	assert.Nil(t, execution.Data(&data1))
	assert.Equal(t, resData.Message, data1.Message)

	service.Close()
	wg.Wait()
}

func TestMultipleListenCall(t *testing.T) {
	var (
		taskKey = "1"
		data    = taskRequest{"https://mesg.com"}
	)

	service, server := newServiceAndServer(t)
	go server.Start()

	makeSureListeningC := make(chan struct{}, 0)
	taskable := Task(taskKey, func(*Execution) (string, Data) {
		close(makeSureListeningC)
		return "", ""
	})

	go service.Listen(taskable)
	server.Execute(taskKey, data)
	<-makeSureListeningC

	assert.Equal(t, service.Listen(taskable).Error(), errAlreadyListening{}.Error())
}

func TestNonExistentTaskExecutionRequest(t *testing.T) {
	var (
		taskKey            = "1"
		nonExistentTaskKey = "2"
		data               = taskRequest{"https://mesg.com"}
	)

	server := mesgtest.NewServer()
	go server.Start()

	reader, writer := io.Pipe()
	service, _ := New(
		DialOption(server.Socket()),
		TokenOption(token),
		EndpointOption(endpoint),
		LogOutputOption(writer),
	)

	go service.Listen(Task(taskKey, func(*Execution) (string, Data) { return "", "" }))
	go server.Execute(nonExistentTaskKey, data)

	line, _, err := bufio.NewReader(reader).ReadLine()
	assert.Nil(t, err)
	assert.Contains(t, errNonExistentTask{name: nonExistentTaskKey}.Error(), string(line))
}
