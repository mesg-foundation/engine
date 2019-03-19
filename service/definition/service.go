package definition

import (
	"github.com/mesg-foundation/core/service"
)

// mainServiceKey is the reserved key for main service.
const mainServiceKey = "service"

// Service represents the service's definition used in mesg.yml files.
type Service struct {
	// Name of the service.
	Name string `yaml:"name" json:"name,omitempty" validate:"required,printascii,min=1"`

	// Sid of the service. It must be unique.
	Sid string `yaml:"sid" json:"sid,omitempty" validate:"omitempty,printascii,max=63,domain"`

	// Description of the service.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Repository's url of the service.
	Repository string `yaml:"repository" json:"repository,omitempty" validate:"omitempty,uri"`

	// Tasks defined by the service.
	Tasks map[string]Task `yaml:"tasks" json:"tasks,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Events defined by the service.
	Events map[string]Event `yaml:"events" json:"events,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Configuration of the the service's container.
	Configuration Dependency `yaml:"configuration" json:"configuration,omitempty"`

	// Dependencies are containers the service depends on. Dependencies will be started and stopped alongside the service.
	Dependencies map[string]Dependency `yaml:"dependencies" json:"dependencies,omitempty" validate:"dive,keys,printascii,ne=service,endkeys,required"`
}

// Event describes a service event.
type Event struct {
	// Name of the event.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description of the event.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Data holds the definition of the event.
	Data map[string]Parameter `yaml:"data" json:"data,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Dependency describes a container configuration.
type Dependency struct {
	// Image of the container to use.
	Image string `yaml:"image" json:"image,omitempty" validate:"printascii"`

	// Volumes required by the container.
	Volumes []string `yaml:"volumes" json:"volumes,omitempty" validate:"unique,dive,printascii"`

	// VolumesFrom indicates to also mount other dependencies' volumes.
	VolumesFrom []string `yaml:"volumesfrom" json:"volumesFrom,omitempty" validate:"unique,dive,printascii"`

	// Ports to publish on the public network.
	Ports []string `yaml:"ports" json:"ports,omitempty" validate:"unique,dive,portmap"`

	// Command to execute when container starts.
	Command string `yaml:"command" json:"command,omitempty" validate:"printascii"`

	// Args to pass to the container.
	Args []string `yaml:"args" json:"args,omitempty" validate:"dive,printascii"`

	// Env is the environment variables in key=value format to pass to the container.
	Env []string `yaml:"env" json:"env,omitempty" validate:"unique,dive,printascii,env"`
}

// Task describes a service task.
type Task struct {
	// Name of the task.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description of the task.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Inputs is the definition of the task's inputs.
	Inputs map[string]Parameter `yaml:"inputs" json:"inputs,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Outputs is the definition of the task's outputs.
	Outputs map[string]Output `yaml:"outputs" json:"outputs,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Output describes task output.
type Output struct {
	// Name of the output.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description of the output.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Data describes the parameters of the output.
	Data map[string]Parameter `yaml:"data" json:"data,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Parameter describes the task's inputs, the task's outputs, and the event's data.
type Parameter struct {
	// Name of the parameter.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description of the parameter.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Type of the parameter's data.
	Type string `yaml:"type" json:"type,omitempty" validate:"required,printascii,oneof=String Number Boolean Object Any"`

	// Optional indicates the parameter is optional.
	Optional bool `yaml:"optional" json:"optional,omitempty"`

	// Repeated indicates the parameter is an array.
	Repeated bool `yaml:"repeated" json:"repeated,omitempty"`

	// Definition of the structure of the object when the type is object.
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
		Args:        dep.Args,
		Env:         dep.Env,
	}
}

func toServiceDependencies(deps map[string]Dependency) []*service.Dependency {
	ds := make([]*service.Dependency, 0, len(deps))
	for key, dep := range deps {
		ds = append(ds, toServiceDependency(key, dep))
	}
	return ds
}
