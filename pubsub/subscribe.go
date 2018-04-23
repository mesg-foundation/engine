package pubsub

// Subscribe to a specific event associated to the service
func Subscribe(key string) (res chan Message) {
	res = make(chan Message)
	mu.Lock()
	defer mu.Unlock()
	if listeners[key] == nil {
		listeners[key] = make([]chan Message, 0)
	}
	listeners[key] = append(listeners[key], res)
	return
}
