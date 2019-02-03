package logger

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/client/service"
	"github.com/mesg-foundation/core/client/service/servicetest"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	require.NotNil(t, s)
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
	require.NoError(t, err)
	require.Equal(t, "success", execution.Key())

	var resp successResponse
	require.Nil(t, execution.Data(&resp))
	require.Equal(t, "ok", resp.Message)
}

func TestListenError(t *testing.T) {
	data := "data"

	s, server := newServiceAndServer(t)
	logger := New(s)

	go server.Start()
	go logger.Start()

	_, execution, err := server.Execute("log", data)
	require.NoError(t, err)
	require.Equal(t, "error", execution.Key())

	var resp errorResponse
	require.Nil(t, execution.Data(&resp))
	require.Contains(t, "json", resp.Message)
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
		require.Nil(t, logger.Start())
	}()

	_, _, err := server.Execute("log", data)
	require.NoError(t, err)
	require.Nil(t, logger.Close())

	_, _, err = server.Execute("log", data)
	require.Equal(t, servicetest.ErrConnectionClosed{}, err)

	wg.Wait()
}
