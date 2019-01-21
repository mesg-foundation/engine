package logger

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/client/service"
	"github.com/mesg-foundation/core/client/service/servicetest"
	"github.com/stvp/assert"
)

const token = "token"
const endpoint = "endpoint"

func newServiceAndServer(t *testing.T) (*service.Service, *servicetest.Server) {
	testServer := servicetest.NewServer()
	s, err := service.New(
		service.DialOption(testServer.Socket()),
		service.TokenOption(token),
		service.EndpointOption(endpoint),
	)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	return s, testServer
}

func TestListenSuccess(t *testing.T) {
	data := logRequest{
		ServiceID: "id",
		Data:      "data",
	}

	s, server := newServiceAndServer(t)
	logger := New(s, LogOutputOption(ioutil.Discard))

	go server.Start()
	go logger.Start()

	_, execution, err := server.Execute("log", data)
	assert.Nil(t, err)
	assert.Equal(t, "success", execution.Key())

	var resp successResponse
	assert.Nil(t, execution.Data(&resp))
	assert.Equal(t, "ok", resp.Message)
}

func TestListenError(t *testing.T) {
	data := "data"

	s, server := newServiceAndServer(t)
	logger := New(s)

	go server.Start()
	go logger.Start()

	_, execution, err := server.Execute("log", data)
	assert.Nil(t, err)
	assert.Equal(t, "error", execution.Key())

	var resp errorResponse
	assert.Nil(t, execution.Data(&resp))
	assert.Contains(t, "json", resp.Message)
}

func TestClose(t *testing.T) {
	data := "data"

	s, server := newServiceAndServer(t)
	logger := New(s)

	go server.Start()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, logger.Start())
	}()

	_, _, err := server.Execute("log", data)
	assert.Nil(t, err)
	assert.Nil(t, logger.Close())

	_, _, err = server.Execute("log", data)
	assert.Equal(t, servicetest.ErrConnectionClosed{}, err)

	wg.Wait()
}
