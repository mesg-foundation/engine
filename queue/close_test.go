package queue

import (
	"os"
	"testing"

	"github.com/stvp/assert"
)

func TestCloseEmptyQueue(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	queue := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}

	err := queue.connect()
	assert.Nil(t, err)

	err = queue.Close()
	assert.Nil(t, err)
}

func TestColseMultipleChannel(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	queue := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}
	err := queue.connect()
	assert.Nil(t, err)

	_, err = queue.createInternalChannel(Channel{Kind: Events, Name: "*"})
	assert.Nil(t, err)

	err = queue.Close()
	assert.Nil(t, err)
}

func TestColseMultipleChannelDisconnect(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	queue := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}
	err := queue.connect()
	assert.Nil(t, err)

	_, err = queue.createInternalChannel(Channel{Kind: Events, Name: "*"})
	assert.Nil(t, err)

	err = queue.disconnect()
	assert.Nil(t, err)

	err = queue.Close()
	assert.NotNil(t, err)
}
