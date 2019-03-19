package definition

import (
	"github.com/mesg-foundation/core/service"
)

// mainServiceKey is key for main service.
const mainServiceKey = "service"

// Service represents MESG services configurations.
type Service struct {
	// Name is the service name.
	Name string `yaml:"name" json:"name,omitempty" validate:"required,printascii,min=1"`

	// Sid is the service id. It must be unique.
	Sid string `yaml:"sid" json:"sid,omitempty" validate:"omitempty,printascii,max=63,domain"`

	// Description is service description.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Repository holds the service's repository url if it's living on a git host.
	Repository string `yaml:"repository" json:"repository,omitempty" validate:"omitempty,uri"`

	// Tasks are the list of tasks that service can execute.
	Tasks map[string]Task `yaml:"tasks" json:"tasks,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Events are the list of events that service can emit.
	Events map[string]Event `yaml:"events" json:"events,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Configuration is the Docker container that service runs inside.
	Configuration Dependency `yaml:"configuration" json:"configuration,omitempty"`

	// Dependencies are the Docker containers that service can depend on.
	Dependencies map[string]Dependency `yaml:"dependencies" json:"dependencies,omitempty" validate:"dive,keys,printascii,ne=service,endkeys,required"`
}

// Event describes a service task.
type Event struct {
	// Name is the name of event.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of event.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Data holds the input inputs of event.
	Data map[string]Parameter `yaml:"data" json:"data,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	// Image is the Docker image.
	Image string `yaml:"image" json:"image,omitempty" validate:"printascii"`

	// Volumes are the Docker volumes.
	Volumes []string `yaml:"volumes" json:"volumes,omitempty" validate:"unique,dive,printascii"`

	// VolumesFrom are the docker volumes-from from.
	VolumesFrom []string `yaml:"volumesFrom" json:"volumesFrom,omitempty" validate:"unique,dive,printascii"`

	// Ports holds ports configuration for container.
	Ports []string `yaml:"ports" json:"ports,omitempty" validate:"unique,dive,portmap"`

	// Command is the Docker command which will be executed when container started.
	Command string `yaml:"command" json:"command,omitempty" validate:"printascii"`

	// Args hold the args to pass to the Docker container
	Args []string `yaml:"args" json:"args,omitempty" validate:"dive,printascii"`

	// Env is a slice of environment variables in key=value format.
	Env []string `yaml:"env" json:"env,omitempty" validate:"unique,dive,printascii,env"`
}

// Task describes a service task.
type Task struct {
	// Name is the name of task.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of task.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Parameters are the definition of the execution inputs of task.
	Inputs map[string]Parameter `yaml:"inputs" json:"inputs,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Outputs are the definition of the execution results of task.
	Outputs map[string]Output `yaml:"outputs" json:"outputs,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Output describes task output.
type Output struct {
	// Name is the name of task output.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of task output.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Data holds the output inputs of a task output.
	Data map[string]Parameter `yaml:"data" json:"data,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Parameter describes task input inputs, output inputs of a task
// output and input inputs of an event.
type Parameter struct {
	// Name is the name of input.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of input.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Type is the data type of input.
	Type string `yaml:"type" json:"type,omitempty" validate:"required,printascii,oneof=String Number Boolean Object Any"`

	// Optional indicates if input is optional.
	Optional bool `yaml:"optional" json:"optional,omitempty"`

	// Repeated is to have an array of this input
	Repeated bool `yaml:"repeated" json:"repeated,omitempty"`

	// Definition of the structure of the object when the type is object
	Object map[string]Parameter `yaml:"object" json:"object,omitempty" validate:"dive,keys,printascii,endkeys,required"`
}

func (s *Service) toService() *service.Service {
	return &service.Service{
		Sid:           s.Sid,
		Name:          s.Name,
		Description:   s.Description,
		Repository:    s.Repository,
		Tasks:         toServiceTasks(s.Tasks),
		Events:        toServiceEvents(s.Events),
		Configuration: toServiceDependency(mainServiceKey, s.Configuration),
		Dependencies:  toServiceDependencies(s.Dependencies),
	}
}

func toServiceTasks(tasks map[string]Task) []*service.Task {
	ts := make([]*service.Task, 0, len(tasks))
	for key, task := range tasks {
		t := &service.Task{
			Key:         key,
			Name:        task.Name,
			Description: task.Description,
			Inputs:      toServiceParameters(task.Inputs),
			Outputs:     make([]*service.Output, 0),
		}
		for key, output := range task.Outputs {
			t.Outputs = append(t.Outputs, &service.Output{
				Key:         key,
				Name:        output.Name,
				Description: output.Description,
				Data:        toServiceParameters(output.Data),
			})
		}
		ts = append(ts, t)
	}
	return ts
}

func toServiceEvents(events map[string]Event) []*service.Event {
	es := make([]*service.Event, 0, len(events))
	for key, event := range events {
		es = append(es, &service.Event{
			Key:         key,
			Name:        event.Name,
			Description: event.Description,
			Data:        toServiceParameters(event.Data),
		})
	}
	return es
}

func toServiceParameters(params map[string]Parameter) []*service.Parameter {
	ps := make([]*service.Parameter, 0, len(params))
	for key, param := range params {
		ps = append(ps, &service.Parameter{
			Key:         key,
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Repeated:    param.Repeated,
			Optional:    param.Optional,
			Object:      toServiceParameters(param.Object),
		})
	}
	return ps
}

func toServiceDependency(key string, dep Dependency) *service.Dependency {
	return &service.Dependency{
		Key:         key,
		Image:       dep.Image,
		Volumes:     dep.Volumes,
		VolumesFrom: dep.VolumesFrom,
		Ports:       dep.Ports,
		Command:     dep.Command,
	}
}

func toServiceDependencies(deps map[string]Dependency) []*service.Dependency {
	ds := make([]*service.Dependency, 0, len(deps))
	for key, dep := range deps {
		ds = append(ds, toServiceDependency(key, dep))
	}
	return ds
}
