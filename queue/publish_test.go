package queue

import (
	"os"
	"testing"

	"github.com/stvp/assert"
)

func TestCreateInternalChannel(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	q := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}

	channel := Channel{
		Kind: Events,
		Name: "*",
	}
	ch, err := q.createInternalChannel(channel)

	assert.Nil(t, err)
	assert.NotNil(t, ch)
	assert.Equal(t, len(q.channels), 1)
	assert.Equal(t, q.channels[channel.namespace()], ch)
}

func TestPublishChannelData(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	type TestPublishChannelDataType struct {
		Foo string
		Bar int
	}
	q := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}

	channel := Channel{
		Kind: Events,
		Name: "test",
	}

	err := q.publish("TestPublishChannelData", channel, TestPublishChannelDataType{
		Foo: "test",
		Bar: 1,
	})
	assert.Nil(t, err)
}

func TestPublishChannelsData(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	type TestPublishChannelsDataType struct {
		Foo string
		Bar int
	}
	q := Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}

	channels := []Channel{
		Channel{
			Kind: Events,
			Name: "test",
		},
	}

	err := q.Publish("TestPublishChannelsData", channels, TestPublishChannelsDataType{
		Foo: "test",
		Bar: 1,
	})
	assert.Nil(t, err)
}
