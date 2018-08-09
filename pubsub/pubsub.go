package pubsub

import (
	"sync"
)

var (
	listeners = make(map[string][]chan Message)
	mu        sync.Mutex
)

// Message sends subscribe/publish messages.
type Message interface{}

// Publish publishes a message to a channel.
func Publish(channel string, data Message) {
	mu.Lock()
	for _, listener := range listeners[channel] {
		if listener != nil {
			listener <- data
		}
	}
	mu.Unlock()
}

// Subscribe subscribes to the channel and returns listener for it.
func Subscribe(channel string) chan Message {
	listener := make(chan Message)
	mu.Lock()
	if listeners[channel] == nil {
		listeners[channel] = make([]chan Message, 0)
	}
	listeners[channel] = append(listeners[channel], listener)
	mu.Unlock()
	return listener
}

// Unsubscribe unsubscribes a subscription for listening channel.
func Unsubscribe(channel string, subscription chan Message) {
	mu.Lock()
	defer mu.Unlock()
	for i, s := range listeners[channel] {
		if s == subscription {
			listeners[channel] = append(listeners[channel][:i], listeners[channel][i+1:]...)
			return
		}
	}
}
