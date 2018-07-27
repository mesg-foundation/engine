package application

import (
	"sync"
	"testing"

	"github.com/stvp/assert"
)

func TestWhenResult(t *testing.T) {
	resultServiceID := "1"
	taskServiceID := "2"
	task := "3"

	app, server := newApplicationAndServer(t)
	go server.Start()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		stream, err := app.
			WhenResult(resultServiceID).
			Execute(taskServiceID, task)
		assert.Nil(t, err)
		assert.NotNil(t, stream)
	}()

	rl := server.LastResultListen()
	assert.Equal(t, resultServiceID, rl.ServiceID())
	assert.Equal(t, "*", rl.KeyFilter())
	assert.Equal(t, "*", rl.TaskFilter())

	wg.Wait()
}
