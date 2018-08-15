package pubsub

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type messageStructTest struct {
	test string
}

func TestPublish(t *testing.T) {
	key := "TestPublish"
	data := messageStructTest{test: "TestPublish"}

	res := Subscribe(key)
	go Publish(key, data)
	x := <-res
	require.Equal(t, x, data)
}

func TestPublishMultipleListeners(t *testing.T) {
	key := "TestPublishMultipleListeners"
	data := messageStructTest{test: "TestPublishMultipleListeners"}
	res1 := Subscribe(key)
	res2 := Subscribe(key)
	go Publish(key, data)
	x := <-res1
	y := <-res2
	require.Equal(t, x, data)
	require.Equal(t, y, data)
}

func TestSubscribe(t *testing.T) {
	key := "TestSubscribe"
	res := Subscribe(key)
	require.NotNil(t, res)
	require.Equal(t, len(listeners[key]), 1)
}

func TestSubscribeMultipleTimes(t *testing.T) {
	key := "TestSubscribeMultipleTimes"
	Subscribe(key)
	Subscribe(key)
	require.Equal(t, len(listeners[key]), 2)
}

func TestUnsubscribe(t *testing.T) {
	key := "TestUnsubscribe"
	channel := Subscribe(key)
	channel1 := Subscribe(key)
	Unsubscribe(key, channel)
	assert.Equal(t, len(listeners[key]), 1)
	assert.NotNil(t, listeners[key])
	Unsubscribe(key, channel1)
	assert.Equal(t, len(listeners[key]), 0)
	assert.Nil(t, listeners[key])
}
