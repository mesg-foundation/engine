package daemon

import (
	"testing"

	"github.com/stvp/assert"
)

func TestIsNotRunning(t *testing.T) {
	runs, err := IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, false, runs)
}

func TestIsRunning(t *testing.T) {
	Start()
	defer Stop()
	err := <-WaitForRunning()
	assert.Nil(t, err)

	runs, err := IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, true, runs)
}

func TestIsNotRunningAfterStop(t *testing.T) {
	Start()
	Stop()
	err := <-WaitForStopped()
	assert.Nil(t, err)

	runs, err := IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, false, runs)
}
