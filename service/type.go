package service

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
)

// Service represents a MESG service.
type Service struct {
	Hash          string       `hash:"-"`      // Hash is calculated from the combination of service's source and mesg.yml. It represents the service uniquely.
	Sid           string       `hash:"name:1"` // Sid is the service id. It needs to be unique and can be used to access to service.
	Name          string       `hash:"name:2"` // Name is the service name.
	Description   string       `hash:"name:3"` // Description is service description.
	Tasks         []Task       `hash:"name:4"` // Tasks are the list of tasks that service can execute.
	Events        []Event      `hash:"name:5"` // Events are the list of events that service can emit.
	Workflows     []Workflow   `hash:"name:5"` // Workflows are the list of Workflows that service can emit.
	Dependencies  []Dependency `hash:"name:6"` // Dependencies are the Docker containers that service can depend on.
	Configuration *Dependency  `hash:"name:8"` // Configuration of the service
	Repository    string       `hash:"name:7"` // Repository holds the service's repository url if it's living on a Git host.
	Source        string       `hash:"name:9"` // Source is the hash id of service's source code on IPFS.
	Owner         types.AccAddress
}

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	Key         string   `hash:"1"`      // Key is the key of dependency.
	Image       string   `hash:"name:2"` // Image is the Docker image.
	Volumes     []string `hash:"name:3"` // Volumes are the Docker volumes.
	VolumesFrom []string `hash:"name:4"` // VolumesFrom are the docker volumes-from from.
	Ports       []string `hash:"name:5"` // Ports holds ports configuration for container.
	Command     string   `hash:"name:6"` // Command is the Docker command which will be executed when container started.
	Args        []string `hash:"name:7"` // Argument holds the args to pass to the Docker container
	Env         []string `hash:"name:8"` // Env is a slice of environment variables in key=value format.
}

// Parameter describes task input parameters, output parameters of a task
// output and input parameters of an event.
type Parameter struct {
	Key         string      `hash:"name:1"` // Key is the key of parameter.
	Name        string      `hash:"name:2"` // Name is the name of parameter.
	Description string      `hash:"name:3"` // Description is the description of parameter.
	Type        string      `hash:"name:4"` // Type is the data type of parameter.
	Optional    bool        `hash:"name:5"` // Optional indicates if parameter is optional.
	Repeated    bool        `hash:"name:6"` // Repeated is to have an array of this parameter
	Object      []Parameter `hash:"name:7"` // Definition of the structure of the object when the type is object
}

// Event describes a service task.
type Event struct {
	Key         string      `hash:"name:1"` // Key is the key of event.
	Name        string      `hash:"name:2"` // Name is the name of event.
	Description string      `hash:"name:3"` // Description is the description of event.
	Data        []Parameter `hash:"name:4"` // Data holds the input parameters of event.
}

// Task describes a service task.
type Task struct {
	Key         string      `hash:"name:1"` // Key is the key of task.
	Name        string      `hash:"name:2"` // Name is the name of task.
	Description string      `hash:"name:3"` // Description is the description of task.
	Inputs      []Parameter `hash:"name:4"` // Inputs are the definition of the execution inputs of task.
	Outputs     []Parameter `hash:"name:5"` // Outputs are the definition of the execution results of task.
}

// Workflow ...
type Workflow struct {
	Trigger Trigger        `hash:"name:1"`
	Tasks   []WorkflowTask `hash:"name:2"`
}

// Trigger ...
type Trigger struct {
	ServiceHash string              `hash:"name:1"`
	EventKey    string              `hash:"name:2"`
	Filter      map[string][]string `hash:"name:3"`
}

// WorkflowTask ...
type WorkflowTask struct {
	ServiceHash string `hash:"name:1"`
	TaskKey     string `hash:"name:2"`
}

// NewService ...
func NewService() Service {
	return Service{}
}

// GetEvent returns event eventKey of service.
func (s Service) GetEvent(eventKey string) (Event, error) {
	for _, event := range s.Events {
		if event.Key == eventKey {
			return event, nil
		}
	}
	return Event{}, fmt.Errorf("EventNotFoundError{EventKey:    eventKey, ServiceName: s.Name}")
}

// GetDependency returns dependency dependencyKey or a not found error.
func (s Service) GetDependency(dependencyKey string) (Dependency, error) {
	for _, dep := range s.Dependencies {
		if dep.Key == dependencyKey {
			return dep, nil
		}
	}
	return Dependency{}, fmt.Errorf("dependency %s do not exist", dependencyKey)
}

// GetTask returns task taskKey of service.
func (s Service) GetTask(taskKey string) (Task, error) {
	for _, task := range s.Tasks {
		if task.Key == taskKey {
			return task, nil
		}
	}
	return Task{}, fmt.Errorf("TaskNotFoundError")
}
