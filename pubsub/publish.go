package pubsub

// Publish a specific event associated to a service
func Publish(key string, data Message) {
	mu.Lock()
	defer mu.Unlock()
	for _, event := range listeners[key] {
		if event != nil {
			event <- data
		}
	}
}
