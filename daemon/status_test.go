package daemon

import (
	"testing"

	"github.com/stvp/assert"
)

func TestIsNotRunning(t *testing.T) {
	<-WaitForFullyStop()
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
	<-WaitForRunning()
	Stop()
	err := <-WaitForFullyStop()
	assert.Nil(t, err)
	runs, err := IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, false, runs)
}
