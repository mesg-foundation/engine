package logger

import (
	"io/ioutil"
	"sync"
	"testing"

	mesg "github.com/mesg-foundation/go-service"
	"github.com/mesg-foundation/go-service/mesgtest"
	"github.com/stvp/assert"
)

const token = "token"
const endpoint = "endpoint"

func newServiceAndServer(t *testing.T) (*mesg.Service, *mesgtest.Server) {
	testServer := mesgtest.NewServer()
	service, err := mesg.New(
		mesg.DialOption(testServer.Socket()),
		mesg.TokenOption(token),
		mesg.EndpointOption(endpoint),
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

	service, server := newServiceAndServer(t)
	logger := New(service, LogOutputOption(ioutil.Discard))

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

	service, server := newServiceAndServer(t)
	logger := New(service)

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

	service, server := newServiceAndServer(t)
	logger := New(service)

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
	assert.Equal(t, mesgtest.ErrConnectionClosed{}, err)

	wg.Wait()
}
