package application

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/ilgooz/mesg-go/application/applicationtest"
	"github.com/stvp/assert"
)

func newApplicationAndServer(t *testing.T) (*Application, *applicationtest.Server) {
	testServer := applicationtest.NewServer()

	app, err := New(
		DialOption(testServer.Socket()),
		EndpointOption(endpoint),
		LogOutputOption(ioutil.Discard),
	)

	assert.Nil(t, err)
	assert.NotNil(t, app)

	return app, testServer
}
func TestExecute(t *testing.T) {
	serviceID := "1"
	task := "2"
	reqData := taskRequest{"https://mesg.tech"}

	app, server := newApplicationAndServer(t)
	go server.Start()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		execution := server.LastExecute()
		assert.Equal(t, serviceID, execution.ServiceID())
		assert.Equal(t, task, execution.Task())

		var data taskRequest
		assert.Nil(t, execution.Decode(&data))
		assert.Equal(t, reqData.URL, data.URL)
	}()

	executionID, err := app.Execute(serviceID, task, reqData)
	assert.Nil(t, err)
	assert.True(t, executionID != "")

	wg.Wait()
}
