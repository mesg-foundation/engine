package pubsub

// Subscribe to a channel
func Subscribe(channel string) chan Message {
	listener := make(chan Message)
	mu.Lock()
	defer mu.Unlock()
	if listeners[channel] == nil {
		listeners[channel] = make([]chan Message, 0)
	}
	listeners[channel] = append(listeners[channel], listener)
	return listener
}
