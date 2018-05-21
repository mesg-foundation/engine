package dependency

// Status return the status of the Docker Swarm Servicer
func Status(namespace string, name string) (status StatusType) {
	dockerService, err := Service(namespace, name)
	status = STOPPED
	if err == nil && dockerService.ID != "" {
		status = RUNNING
	}
	return
}
