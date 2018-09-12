package importer

// ServiceDefinition represents MESG services configurations.
type ServiceDefinition struct {
	// Name is the service name.
	Name string `yaml:"name"`

	// Description is service description.
	Description string `yaml:"description"`

	// Tasks are the list of tasks that service can execute.
	Tasks map[string]*Task `yaml:"tasks"`

	// Events are the list of events that service can emit.
	Events map[string]*Event `yaml:"events"`

	// Dependencies are the Docker containers that service can depend on.
	Dependencies map[string]*Dependency `yaml:"dependencies"`

	// Configuration is the Docker container that service runs inside.
	Configuration *Dependency `yaml:"configuration"`

	// Repository holds the service's repository url if it's living on
	// a Git host.
	Repository string `yaml:"repository"`
}

// Event describes a service task.
type Event struct {
	// Name is the name of event.
	Name string `yaml:"name"`

	// Description is the description of event.
	Description string `yaml:"description"`

	// Data holds the input parameters of event.
	Data map[string]*Parameter `yaml:"data"`
}

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	// Image is the Docker image.
	Image string `yaml:"image"`

	// Volumes are the Docker volumes.
	Volumes []string `yaml:"volumes"`

	// VolumesFrom are the docker volumes-from from.
	VolumesFrom []string `yaml:"volumesfrom"`

	// Ports holds ports configuration for container.
	Ports []string `yaml:"ports"`

	// Command is the Docker command which will be executed when container started.
	Command string `yaml:"command"`
}

// Task describes a service task.
type Task struct {
	// Name is the name of task.
	Name string `yaml:"name"`

	// Description is the description of task.
	Description string `yaml:"description"`

	// Inputs are the definition of the execution inputs of task.
	Inputs map[string]*Parameter `yaml:"inputs"`

	// Outputs are the definition of the execution results of task.
	Outputs map[string]*Output `yaml:"outputs"`
}

// Output describes task output.
type Output struct {
	// Name is the name of task output.
	Name string `yaml:"name"`

	// Description is the description of task output.
	Description string `yaml:"description"`

	// Data holds the output parameters of a task output.
	Data map[string]*Parameter `yaml:"data"`
}

// Parameter describes task input parameters, output parameters of a task
// output and input parameters of an event.
type Parameter struct {
	// Name is the name of parameter.
	Name string `yaml:"name"`

	// Description is the description of parameter.
	Description string `yaml:"description"`

	// Type is the data type of parameter.
	Type string `yaml:"type"`

	// Optional indicates if parameter is optional.
	Optional bool `yaml:"optional"`
}
