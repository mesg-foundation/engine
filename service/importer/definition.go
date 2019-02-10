package importer

// ConfigurationDependencyKey is the reserved key of the service's configuration in the dependencies array.
const ConfigurationDependencyKey = "service"

// ServiceDefinition represents MESG services configurations.
type ServiceDefinition struct {
	// Name is the service name.
	Name string `yaml:"name" json:"name,omitempty" validate:"required,printascii,min=1"`

	// Sid is the service id. It must be unique.
	Sid string `yaml:"sid" json:"sid,omitempty" validate:"printascii,max=39"`

	// Description is service description.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Tasks are the list of tasks that service can execute.
	Tasks map[string]*Task `yaml:"tasks" json:"tasks,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Events are the list of events that service can emit.
	Events map[string]*Event `yaml:"events" json:"events,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Dependencies are the Docker containers that service can depend on.
	Dependencies map[string]*Dependency `yaml:"dependencies" json:"dependencies,omitempty" validate:"dive,keys,printascii,ne=service,endkeys,required"`

	// Configuration is the Docker container that service runs inside.
	Configuration *Dependency `yaml:"configuration" json:"configuration,omitempty"`

	// Repository holds the service's repository url if it's living on
	// a Git host.
	Repository string `yaml:"repository" json:"repository,omitempty" validate:"omitempty,uri"`
}

// Event describes a service task.
type Event struct {
	// Name is the name of event.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of event.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Data holds the input parameters of event.
	Data map[string]*Parameter `yaml:"data" json:"data,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	// Image is the Docker image.
	Image string `yaml:"image" json:"image,omitempty" validate:"printascii"`

	// Volumes are the Docker volumes.
	Volumes []string `yaml:"volumes" json:"volumes,omitempty" validate:"unique,dive,printascii"`

	// VolumesFrom are the docker volumes-from from.
	VolumesFrom []string `yaml:"volumesfrom" json:"volumefrom,omitempty" validate:"unique,dive,printascii"`

	// Ports holds ports configuration for container.
	Ports []string `yaml:"ports" json:"ports,omitempty" validate:"unique,dive,numeric"`

	// Command is the Docker command which will be executed when container started.
	Command string `yaml:"command" json:"command,omitempty" validate:"printascii"`

	// Args hold the args to pass to the Docker container
	Args []string `yaml:"args" json:"args,omitempty" validate:"dive,printascii"`

	// Env is a slice of environment variables in key=value format.
	Env []string `yaml:"env" json:"env,omitempty" validate:"unique,dive,printascii"`
}

// Task describes a service task.
type Task struct {
	// Name is the name of task.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of task.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Inputs are the definition of the execution inputs of task.
	Inputs map[string]*Parameter `yaml:"inputs" json:"inputs,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Outputs are the definition of the execution results of task.
	Outputs map[string]*Output `yaml:"outputs" json:"outputs,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Output describes task output.
type Output struct {
	// Name is the name of task output.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of task output.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Data holds the output parameters of a task output.
	Data map[string]*Parameter `yaml:"data" json:"data,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Parameter describes task input parameters, output parameters of a task
// output and input parameters of an event.
type Parameter struct {
	// Name is the name of parameter.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of parameter.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Type is the data type of parameter.
	Type string `yaml:"type" json:"type,omitempty" validate:"required,printascii,oneof=String Number Boolean Object Any"`

	// Optional indicates if parameter is optional.
	Optional bool `yaml:"optional" json:"optional,omitempty"`

	// Repeated is to have an array of this parameter
	Repeated bool `yaml:"repeated" json:"repeated,omitempty"`

	// Definition of the structure of the object when the type is object
	Object map[string]*Parameter `yaml:"object" json:"object,omitempty" validate:"dive,keys,printascii,endkeys,required"`
}
