package daemon

import (
	"errors"
	"testing"
	"time"

	"github.com/mesg-foundation/core/container"
	"github.com/stvp/assert"
)

func testForceAndWaitForFullStop() chan error {
	start := time.Now()
	timeout := time.Minute
	wait := make(chan error, 1)
	go func() {
		for {
			err := Stop()
			if err != nil {
				wait <- err
				return
			}
			status, err := Status()
			if err != nil {
				wait <- err
				return
			}
			if status == container.STOPPED {
				close(wait)
				return
			}
			diff := time.Now().Sub(start)
			if diff.Nanoseconds() >= int64(timeout) {
				wait <- errors.New("Wait too long for the MESG Core to fully stop, timeout reached")
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return wait
}

func TestIsNotRunning(t *testing.T) {
	<-testForceAndWaitForFullStop()
	status, err := Status()
	assert.Nil(t, err)
	assert.Equal(t, container.STOPPED, status)
}

func TestIsRunning(t *testing.T) {
	startForTest()
	status, err := Status()
	assert.Nil(t, err)
	assert.Equal(t, container.RUNNING, status)
}
