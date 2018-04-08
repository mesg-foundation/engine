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
	queue.connect()
	err := queue.Close()
	assert.Nil(t, err)
}

func TestColseMultipleChannel(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	queue := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}
	queue.connect()
	queue.createInternalChannel(Channel{Kind: Events, Name: "*"})
	err := queue.Close()
	assert.Nil(t, err)
}

func TestColseMultipleChannelDisconnect(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	queue := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}
	queue.connect()
	queue.createInternalChannel(Channel{Kind: Events, Name: "*"})
	queue.disconnect()
	err := queue.Close()
	assert.NotNil(t, err)
}
