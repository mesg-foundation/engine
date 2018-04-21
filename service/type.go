package service

import "github.com/mesg-foundation/application/types"

// Visibility is the tags to set is the service is visible for whom
type Visibility string

// List of visibilities flags
const (
	VisibilityAll     Visibility = "ALL"
	VisibilityUsers   Visibility = "USERS"
	VisibilityWorkers Visibility = "WORKERS"
	VisibilityNone    Visibility = "NONE"
)

// Publish let you configure the part of your service you want to publish
type Publish string

// List of all publishes flags
const (
	PublishAll       Publish = "ALL"
	PublishSource    Publish = "SOURCE"
	PublishContainer Publish = "CONTAINER"
	PublishNone      Publish = "NONE"
)

// Service is a definition for a service to run
type Service types.ProtoService

// Task is a definition of a Task from a service
type Task types.ProtoTask

// Fee is the different fees to apply
type Fee types.ProtoFee

// Event is the definition of an event emitted from a service
type Event types.ProtoEvent

// Parameter is the definition of a parameter for a Task
type Parameter types.ProtoParameter

// Dependency is the docker information about the Dependency
type Dependency types.ProtoDependency

// New returns a new service based on the proto type
func New(service *types.ProtoService) (res *Service) {
	res = new(Service)
	res.Name = service.Name
	res.Description = service.Description
	res.Publish = service.Publish
	res.Visibility = service.Visibility
	res.Events = service.Events
	res.Tasks = service.Tasks
	res.Dependencies = service.Dependencies
	return
}

// GetDependencies returns the dependencies according to the service types
func (service *Service) GetDependencies() (dependencies map[string]*Dependency) {
	dependencies = make(map[string]*Dependency)
	for name, dependency := range service.Dependencies {
		dependencies[name] = &Dependency{
			Image:   dependency.Image,
			Ports:   dependency.Ports,
			Command: dependency.Command,
			Volumes: dependency.Volumes,
		}
	}
	return
}
