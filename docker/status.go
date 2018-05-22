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
func Status(name []string) (status StatusType) {
	dockerService, err := FindService(name)
	status = STOPPED
	if err == nil && dockerService != nil && dockerService.ID != "" {
		status = RUNNING
	}
	return
}

// IsRunning returns true if the dependency is running, false otherwise
func IsRunning(name []string) (running bool) {
	running = Status(name) == RUNNING
	return
}

// IsStopped returns true if the dependency is stopped, false otherwise
func IsStopped(name []string) (running bool) {
	running = Status(name) == STOPPED
	return
}
