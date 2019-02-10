package importer

// ServiceDefinition represents MESG services configurations.
type ServiceDefinition struct {
	// Name is the service name.
	Name string `yaml:"name" json:"name,omitempty"`

	// Sid is the service id. It must be unique.
	Sid string `yaml:"sid" json:"sid,omitempty"`

	// Description is service description.
	Description string `yaml:"description" json:"description,omitempty"`

	// Tasks are the list of tasks that service can execute.
	Tasks map[string]*Task `yaml:"tasks" json:"tasks,omitempty"`

	// Events are the list of events that service can emit.
	Events map[string]*Event `yaml:"events" json:"events,omitempty"`

	// Dependencies are the Docker containers that service can depend on.
	Dependencies map[string]*Dependency `yaml:"dependencies" json:"dependencies,omitempty"`

	// Configuration is the Docker container that service runs inside.
	Configuration *Dependency `yaml:"configuration" json:"configuration,omitempty"`

	// Repository holds the service's repository url if it's living on
	// a Git host.
	Repository string `yaml:"repository" json:"repository,omitempty"`
}

// Event describes a service task.
type Event struct {
	// Name is the name of event.
	Name string `yaml:"name" json:"name,omitempty"`

	// Description is the description of event.
	Description string `yaml:"description" json:"description,omitempty"`

	// Data holds the input parameters of event.
	Data map[string]*Parameter `yaml:"data" json:"data,omitempty"`
}

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	// Image is the Docker image.
	Image string `yaml:"image" json:"image,omitempty"`

	// Volumes are the Docker volumes.
	Volumes []string `yaml:"volumes" json:"volumes,omitempty"`

	// VolumesFrom are the docker volumes-from from.
	VolumesFrom []string `yaml:"volumesfrom" json:"volumefrom,omitempty"`

	// Ports holds ports configuration for container.
	Ports []string `yaml:"ports" json:"ports,omitempty"`

	// Command is the Docker command which will be executed when container started.
	Command string `yaml:"command" json:"command,omitempty"`

	// Args hold the args to pass to the Docker container
	Args []string `yaml:"args" json:"args,omitempty"`

	// Env is a slice of environment variables in key=value format.
	Env []string `yaml:"env" json:"env,omitempty"`
}

// Task describes a service task.
type Task struct {
	// Name is the name of task.
	Name string `yaml:"name" json:"name,omitempty"`

	// Description is the description of task.
	Description string `yaml:"description" json:"description,omitempty"`

	// Inputs are the definition of the execution inputs of task.
	Inputs map[string]*Parameter `yaml:"inputs" json:"inputs,omitempty"`

	// Outputs are the definition of the execution results of task.
	Outputs map[string]*Output `yaml:"outputs" json:"outputs,omitempty"`
}

// Output describes task output.
type Output struct {
	// Name is the name of task output.
	Name string `yaml:"name" json:"name,omitempty"`

	// Description is the description of task output.
	Description string `yaml:"description" json:"description,omitempty"`

	// Data holds the output parameters of a task output.
	Data map[string]*Parameter `yaml:"data" json:"data,omitempty"`
}

// Parameter describes task input parameters, output parameters of a task
// output and input parameters of an event.
type Parameter struct {
	// Name is the name of parameter.
	Name string `yaml:"name" json:"name,omitempty"`

	// Description is the description of parameter.
	Description string `yaml:"description" json:"description,omitempty"`

	// Type is the data type of parameter.
	Type string `yaml:"type" json:"type,omitempty"`

	// Optional indicates if parameter is optional.
	Optional bool `yaml:"optional" json:"optional,omitempty"`

	// Repeated is to have an array of this parameter
	Repeated bool `yaml:"repeated" json:"repeated,omitempty"`

	// Definition of the structure of the object when the type is object
	Object map[string]*Parameter `yaml:"object" json:"object,omitempty"`
}
