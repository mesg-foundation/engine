package logger

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/ilgooz/mesg-go/service"
	"github.com/ilgooz/mesg-go/service/servicetest"
	"github.com/stvp/assert"
)

const token = "token"
const endpoint = "endpoint"

func newServiceAndServer(t *testing.T) (*service.Service, *servicetest.Server) {
	testServer := servicetest.NewServer()
	service, err := service.New(
		service.DialOption(testServer.Socket()),
		service.TokenOption(token),
		service.EndpointOption(endpoint),
	)
	assert.Nil(t, err)
	assert.NotNil(t, service)
	return service, testServer
}

func TestListenSuccess(t *testing.T) {
	data := logRequest{
		ServiceID: "id",
		Data:      "data",
	}

	s, server := newServiceAndServer(t)
	go server.Start()

	l := New(s, LogOutputOption(ioutil.Discard))
	go l.Start()

	_, execution, err := server.Execute("log", data)
	assert.Nil(t, err)
	assert.Equal(t, "success", execution.Key())

	var resp successResponse
	assert.Nil(t, execution.Decode(&resp))
	assert.Equal(t, "ok", resp.Message)
}

func TestListenError(t *testing.T) {
	data := "data"

	s, server := newServiceAndServer(t)
	go server.Start()

	l := New(s)
	go l.Start()

	_, execution, err := server.Execute("log", data)
	assert.Nil(t, err)
	assert.Equal(t, "error", execution.Key())

	var resp errorResponse
	assert.Nil(t, execution.Decode(&resp))
	assert.Contains(t, "json", resp.Message)
}

func TestClose(t *testing.T) {
	data := "data"

	s, server := newServiceAndServer(t)
	go server.Start()

	l := New(s)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		assert.NotNil(t, l.Start())
		wg.Done()
	}()
	_, _, err := server.Execute("log", data)
	assert.Nil(t, err)
	assert.Nil(t, l.Close())

	_, _, err = server.Execute("log", data)
	assert.NotNil(t, err)

	wg.Wait()
}
