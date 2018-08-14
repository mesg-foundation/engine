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
	defer mu.Unlock()
	for _, listener := range listeners[channel] {
		listener <- data
	}
}

// Subscribe subscribes to the channel and returns listener for it.
func Subscribe(channel string) chan Message {
	mu.Lock()
	defer mu.Unlock()
	listener := make(chan Message)
	if listeners[channel] == nil {
		listeners[channel] = make([]chan Message, 0)
	}
	listeners[channel] = append(listeners[channel], listener)
	return listener
}

// Unsubscribe unsubscribes a listener from listening channel.
func Unsubscribe(channel string, listener chan Message) {
	mu.Lock()
	defer mu.Unlock()
	for i, l := range listeners[channel] {
		if l == listener {
			listeners[channel] = append(listeners[channel][:i], listeners[channel][i+1:]...)
			if len(listeners[channel]) == 0 {
				listeners[channel] = nil
			}
			return
		}
	}
}
