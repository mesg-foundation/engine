package service

import (
	dockerDependency "github.com/mesg-foundation/core/docker/dependency"
)

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

func (service *Service) GetDependency(name string) dockerDependency.Dependency {
	return dockerDependency.Dependency(service.Dependencies[name])
}

func (service *Service) GetDependenciesKeys() (keys []string) {
	keys = make([]string, len(service.Dependencies))
	for k := range service.Dependencies {
		keys = append(keys, k)
	}
	return
}
