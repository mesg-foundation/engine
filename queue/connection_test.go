package queue

import (
	"os"
	"testing"

	"github.com/stvp/assert"
)

func TestConnect(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	queue := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}
	err := queue.connect()
	assert.Nil(t, err)
}

func TestDoubleConnect(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	queue := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}
	queue.connect()
	err := queue.connect()
	assert.Nil(t, err)
}

func TestDisconnect(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	queue := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}
	queue.connect()
	err := queue.disconnect()
	assert.Nil(t, err)
}

func TestDoubleDisconnect(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	queue := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}
	queue.connect()
	queue.disconnect()
	err := queue.disconnect()
	assert.NotNil(t, err)
}
