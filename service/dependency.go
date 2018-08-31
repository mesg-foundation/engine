package service

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	// Key is the key of dependency.
	Key string `hash:"1"`

	// Image is the Docker image.
	Image string `hash:"name:2"`

	// Volumes are the Docker volumes.
	Volumes []string `hash:"name:3"`

	// VolumesFrom are the docker volumes-from from.
	VolumesFrom []string `hash:"name:4"`

	// Ports holds ports configuration for container.
	Ports []string `hash:"name:5"`

	// Command is the Docker command which will be executed when container started.
	Command string `hash:"name:6"`

	// service is the dependency's service.
	service *Service `hash:"-"`
}
