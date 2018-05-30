package daemon

import (
	"errors"
	"testing"
	"time"

	"github.com/stvp/assert"
)

func testForceAndWaitForFullStop() (wait chan error) {
	start := time.Now()
	timeout := time.Minute
	wait = make(chan error, 1)
	go func() {
		for {
			err := Stop()
			if err != nil {
				wait <- err
				return
			}
			stopped, err := IsStopped()
			if err != nil {
				wait <- err
				return
			}
			if stopped == true {
				close(wait)
				return
			}
			diff := time.Now().Sub(start)
			if diff.Nanoseconds() >= int64(timeout) {
				wait <- errors.New("Wait too long for the daemon to fully stop, timeout reached")
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return
}

func TestIsNotRunning(t *testing.T) {
	<-testForceAndWaitForFullStop()
	runs, err := IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, false, runs)
	stopped, err := IsStopped()
	assert.Nil(t, err)
	assert.Equal(t, true, stopped)
}

func TestIsRunning(t *testing.T) {
	startForTest()
	runs, err := IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, true, runs)
	stopped, err := IsStopped()
	assert.Nil(t, err)
	assert.Equal(t, false, stopped)
}
