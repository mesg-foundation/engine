package daemon

func IsRunning() (running bool, err error) {
	container, err := Container()
	if err != nil {
		return
	}
	running = container != nil
	return
}
