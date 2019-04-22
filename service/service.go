package service

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/mesg-foundation/core/service/importer"
)

// WARNING about hash tags on Service type and its inner types:
// * never change the name attr of hash tag. use an incremented value for
// name attr when a new configuration field added to Service.
// * don't increment the value of name attr if corresponding field's name
// changed but its behavior remains the same.
// * this is required for not breaking Service IDs unless there is a behavioral
// change.

// MainServiceKey is key for main service.
const MainServiceKey = importer.ConfigurationDependencyKey

// Service represents a MESG service.
type Service struct {
	// Hash is calculated from the combination of service's source and mesg.yml.
	// It represents the service uniquely.
	Hash string `hash:"-"`

	// HashVersion is the version of the algorithm used to calculate the of the hash
	HashVersion string `hash:"-"`

	// Sid is the service id.
	// It needs to be unique and can be used to access to service.
	Sid string `hash:"name:1"`

	// Name is the service name.
	Name string `hash:"name:2"`

	// Description is service description.
	Description string `hash:"name:3"`

	// Tasks are the list of tasks that service can execute.
	Tasks []*Task `hash:"name:4"`

	// Events are the list of events that service can emit.
	Events []*Event `hash:"name:5"`

	// Dependencies are the Docker containers that service can depend on.
	Dependencies []*Dependency `hash:"name:6"`

	// Configuration of the service
	Configuration *Dependency `hash:"name:8"`

	// Repository holds the service's repository url if it's living on
	// a Git host.
	Repository string `hash:"name:7"`

	// DeployedAt holds the creation time of service.
	DeployedAt time.Time `hash:"-"`
}

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

// Log holds log streams of dependency.
type Log struct {
	Dependency string
	Standard   io.ReadCloser
	Error      io.ReadCloser
}

// getDependency returns dependency dependencyKey or a not found error.
func (s *Service) getDependency(dependencyKey string) (*Dependency, error) {
	for _, dep := range s.Dependencies {
		if dep.Key == dependencyKey {
			return dep, nil
		}
	}
	return nil, fmt.Errorf("dependency %s do not exist", dependencyKey)
}

// ValidateConfigurationEnv checks presence of env variables in mesg.yml under env section.
func (s *Service) ValidateConfigurationEnv(env map[string]string) error {
	var nonDefined []string
	for key := range env {
		exists := false
		// check if "key=" exists in configuration
		for _, env := range s.Configuration.Env {
			if strings.HasPrefix(env, key+"=") {
				exists = true
			}
		}
		if !exists {
			nonDefined = append(nonDefined, key)
		}
	}
	if len(nonDefined) > 0 {
		sort.Strings(nonDefined)
		return ErrNotDefinedEnv{nonDefined}
	}
	return nil
}

// ErrNotDefinedEnv error returned when optionally given env variables
// are not defined in the mesg.yml file.
type ErrNotDefinedEnv struct {
	Env []string
}

func (e ErrNotDefinedEnv) Error() string {
	return fmt.Sprintf("environment variable(s) %q not defined in mesg.yml (under configuration.env key)",
		strings.Join(e.Env, ", "))
}
