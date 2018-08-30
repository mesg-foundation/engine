package service

// Service represents a MESG service.
type Service struct {
	// ID is the unique id of service.
	ID string `hash:"-" yaml:"-"`

	// Name is the service name.
	Name string `hash:"name:1" yaml:"name"`

	// Description is service description.
	Description string `hash:"name:2" yaml:"description"`

	// Tasks are the list of tasks that service can execute.
	Tasks map[string]*Task `hash:"name:3" yaml:"tasks"`

	// Events are the list of events that service can emit.
	Events map[string]*Event `hash:"name:4" yaml:"events"`

	// Dependencies are the Docker containers that service can depend on.
	Dependencies map[string]*Dependency `hash:"name:5" yaml:"dependencies"`

	// Configuration is the Docker container that service runs inside.
	Configuration *Dependency `hash:"name:6" yaml:"configuration"`

	// Repository holds the service's repository url if it's living on
	// a Git host.
	Repository string `hash:"name:7" yaml:"repository"`
}
