package queue_test

import (
	"testing"

	"github.com/mesg-foundation/application/queue"
	"github.com/stvp/assert"
)

func TestCreateChannel(t *testing.T) {
	q := queue.Queue{
		URL: "amqp://guest:guest@localhost:5672/",
	}
	ch := queue.Channel{
		Kind: queue.Events,
		Name: "test",
	}
	err := q.Publish("TestCreateChannel", []queue.Channel{ch}, map[string]string{
		"foo": "bar",
	})
	assert.Nil(t, err)
	q.Close()
}
