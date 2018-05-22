package docker

// StatusType of the service
type StatusType uint

// status for services
const (
	STOPPED StatusType = 0
	RUNNING StatusType = 1
	PARTIAL StatusType = 2
)

// Status return the status of the Docker Swarm Servicer
func Status(namespace string, name string) (status StatusType) {
	dockerService, err := Service(namespace, name)
	status = STOPPED
	if err == nil && dockerService.ID != "" {
		status = RUNNING
	}
	return
}

// IsRunning returns true if the dependency is running, false otherwise
func IsRunning(namespace string, name string) (running bool) {
	running = Status(namespace, name) == RUNNING
	return
}

// IsStopped returns true if the dependency is stopped, false otherwise
func IsStopped(namespace string, name string) (running bool) {
	running = Status(namespace, name) == STOPPED
	return
}
