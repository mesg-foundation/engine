package pubsub

import (
	"testing"

	"github.com/stvp/assert"
)

type MessageStructTest struct{}

func TestPublish(t *testing.T) {
	key := "TestPublish"
	data := MessageStructTest{}

	res := Subscribe(key)
	go Publish(key, data)
	x := <-res
	assert.Equal(t, x, data)
}

func TestPublishMultipleListeners(t *testing.T) {
	key := "TestPublishMultipleListeners"
	data := MessageStructTest{}
	res1 := Subscribe(key)
	res2 := Subscribe(key)
	go Publish(key, data)
	x := <-res1
	y := <-res2
	assert.Equal(t, x, data)
	assert.Equal(t, y, data)
}
