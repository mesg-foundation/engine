package service

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	// Key is the key of dependency.
	Key string `hash:"1" validate:"printascii"`

	// Image is the Docker image.
	Image string `hash:"name:2" validate:"printascii"`

	// Volumes are the Docker volumes.
	Volumes []string `hash:"name:3" validate:"unique,dive,printascii"`

	// VolumesFrom are the docker volumes-from from.
	VolumesFrom []string `hash:"name:4" validate:"unique,dive,printascii"`

	// Ports holds ports configuration for container.
	Ports []string `hash:"name:5"`

	// Command is the Docker command which will be executed when container started.
	Command string `hash:"name:6" validate:"printascii"`

	// Argument holds the args to pass to the Docker container
	Args []string `hash:"name:7" validate:"dive,printascii"`

	// Env is a slice of environment variables in key=value format.
	Env []string `hash:"name:8" validate:"unique,dive,printascii"`
}
