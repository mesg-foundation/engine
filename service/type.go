package service

import "github.com/mesg-foundation/engine/hash"

// MainServiceKey is key for main service.
const MainServiceKey = "service"

// StatusType of the service.
type StatusType uint

// Possible statuses for service.
const (
	UNKNOWN StatusType = iota
	STOPPED
	STARTING
	PARTIAL
	RUNNING
)

func (s StatusType) String() string {
	switch s {
	case STOPPED:
		return "STOPPED"
	case STARTING:
		return "STARTING"
	case PARTIAL:
		return "PARTIAL"
	case RUNNING:
		return "RUNNING"
	default:
		return "UNKNOWN"
	}
}

// WARNING about hash tags on Service type and its inner types:
// * never change the name attr of hash tag. use an incremented value for
// name attr when a new configuration field added to Service.
// * don't increment the value of name attr if corresponding field's name
// changed but its behavior remains the same.
// * this is required for not breaking Service IDs unless there is a behavioral
// change.

// Service represents a MESG service.
type Service struct {
	// Hash is calculated from the combination of service's source and mesg.yml.
	// It represents the service uniquely.
	Hash hash.Hash `hash:"-" validate:"required"`

	// Sid is the service id.
	// It needs to be unique and can be used to access to service.
	Sid string `hash:"name:1"  validate:"required,printascii,max=63,domain"`

	// Name is the service name.
	Name string `hash:"name:2" validate:"required,printascii"`

	// Description is service description.
	Description string `hash:"name:3" validate:"printascii"`

	// Tasks are the list of tasks that service can execute.
	Tasks []*Task `hash:"name:4" validate:"dive,required"`

	// Events are the list of events that service can emit.
	Events []*Event `hash:"name:5" validate:"dive,required"`

	// Dependencies are the Docker containers that service can depend on.
	Dependencies []*Dependency `hash:"name:6" validate:"dive,required"`

	// Configuration of the service
	Configuration *Dependency `hash:"name:8" validate:"required"`

	// Repository holds the service's repository url if it's living on
	// a Git host.
	Repository string `hash:"name:7" validate:"omitempty,uri"`

	// Source is the hash id of service's source code on IPFS.
	Source string `hash:"name:9" validate:"required,printascii"`
}

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
	Ports []string `hash:"name:5" validate:"unique,dive,portmap"`

	// Command is the Docker command which will be executed when container started.
	Command string `hash:"name:6" validate:"printascii"`

	// Argument holds the args to pass to the Docker container
	Args []string `hash:"name:7" validate:"dive,printascii"`

	// Env is a slice of environment variables in key=value format.
	Env []string `hash:"name:8" validate:"unique,dive,env"`
}

// Task describes a service task.
type Task struct {
	// Key is the key of task.
	Key string `hash:"name:1" validate:"printascii"`

	// Name is the name of task.
	Name string `hash:"name:2" validate:"printascii"`

	// Description is the description of task.
	Description string `hash:"name:3" validate:"printascii"`

	// Inputs are the definition of the execution inputs of task.
	Inputs []*Parameter `hash:"name:4" validate:"dive,required"`

	// Outputs are the definition of the execution results of task.
	Outputs []*Parameter `hash:"name:5" validate:"dive,required"`
}

// Event describes a service task.
type Event struct {
	// Key is the key of event.
	Key string `hash:"name:1" validate:"printascii"`

	// Name is the name of event.
	Name string `hash:"name:2" validate:"printascii"`

	// Description is the description of event.
	Description string `hash:"name:3" validate:"printascii"`

	// Data holds the input parameters of event.
	Data []*Parameter `hash:"name:4" validate:"dive,required"`
}

// Parameter describes task input parameters, output parameters of a task
// output and input parameters of an event.
type Parameter struct {
	// Key is the key of parameter.
	Key string `hash:"name:1" validate:"printascii"`

	// Name is the name of parameter.
	Name string `hash:"name:2" validate:"printascii"`

	// Description is the description of parameter.
	Description string `hash:"name:3" validate:"printascii"`

	// Type is the data type of parameter.
	Type string `hash:"name:4" validate:"required,printascii,oneof=String Number Boolean Object Any"`

	// Optional indicates if parameter is optional.
	Optional bool `hash:"name:5"`

	// Repeated is to have an array of this parameter
	Repeated bool `hash:"name:6"`

	// Definition of the structure of the object when the type is object
	Object []*Parameter `hash:"name:7" validate:"unique,dive,required"`
}
