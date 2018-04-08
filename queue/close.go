package queue

// Close see "./interface.go"
func (queue *Queue) Close() (err error) {
	for _, ch := range queue.channels {
		err = ch.Close()
		if err == nil {
			break
		}
	}
	if err == nil {
		err = queue.disconnect()
	}
	return
}
