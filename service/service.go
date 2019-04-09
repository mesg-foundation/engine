package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/utils/dirhash"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/mr-tron/base58"
)

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
	Hash string `hash:"-"`

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

	// Repository holds the service's repository url if it's living on
	// a Git host.
	Repository string `hash:"name:7"`

	// DeployedAt holds the creation time of service.
	DeployedAt time.Time `hash:"-"`
}

// DStatusType indicates the type of status message.
type DStatusType int

const (
	_ DStatusType = iota // skip zero value.

	// DRunning indicates that status message belongs to a continuous state.
	DRunning

	// DDonePositive indicates that status message belongs to a positive noncontinuous state.
	DDonePositive

	// DDoneNegative indicates that status message belongs to a negative noncontinuous state.
	DDoneNegative
)

// DeployStatus represents the deployment status.
type DeployStatus struct {
	Message string
	Type    DStatusType
}

// New creates a new service from contextDir.
func New(contextDir string, c container.Container, statuses chan DeployStatus, env map[string]string) (*Service, error) {
	var err error
	s := &Service{}
	defer s.closeStatus(statuses)

	def, err := importer.From(contextDir)
	if err != nil {
		return nil, err
	}

	dh := dirhash.New(contextDir)
	envbytes := []byte(xos.EnvMapToString(env))
	hash, err := dh.Sum(envbytes)
	if err != nil {
		return nil, err
	}
	s.Hash = base58.Encode(hash)

	s.injectDefinition(def)

	if err := s.validateConfigurationEnv(env); err != nil {
		return nil, err
	}

	// replace default env with new one.
	defenv := xos.EnvSliceToMap(s.configuration().Env)
	s.configuration().Env = xos.EnvMapToSlice(xos.EnvMergeMaps(defenv, env))

	if err := s.deploy(contextDir, c, statuses); err != nil {
		return nil, err
	}

	return s, nil
}

// deploy deploys service.
func (s *Service) deploy(contextDir string, c container.Container, statuses chan DeployStatus) error {
	s.sendStatus(statuses, "Building Docker image...", DRunning)

	imageHash, err := c.Build(contextDir)
	if err != nil {
		return err
	}

	s.sendStatus(statuses, "Image built with success", DDonePositive)

	s.configuration().Image = imageHash
	// TODO: the following test should be moved in New function
	if s.Sid == "" {
		// make sure that sid doesn't have the same length with id.
		s.Sid = "_" + s.Hash
	}
	return nil
}

// sendStatus sends a status message.
func (s *Service) sendStatus(statuses chan DeployStatus, message string, typ DStatusType) {
	if statuses != nil {
		statuses <- DeployStatus{
			Message: message,
			Type:    typ,
		}
	}
}

// closeStatus closes statuses chan.
func (s *Service) closeStatus(statuses chan DeployStatus) {
	if statuses != nil {
		close(statuses)
	}
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

// validateConfigurationEnv checks presence of env variables in mesg.yml under env section.
func (s *Service) validateConfigurationEnv(env map[string]string) error {
	var nonDefined []string
	for key := range env {
		exists := false
		// check if "key=" exists in configuration
		for _, env := range s.configuration().Env {
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

// helper to return the configuration of the service from the dependencies array
func (s *Service) configuration() *Dependency {
	for _, dep := range s.Dependencies {
		if dep.Key == importer.ConfigurationDependencyKey {
			return dep
		}
	}
	return nil
}

// ErrNotDefinedEnv error returned when optionally given env variables
// are not defined in the mesg.yml file.
type ErrNotDefinedEnv struct {
	env []string
}

func (e ErrNotDefinedEnv) Error() string {
	return fmt.Sprintf("environment variable(s) %q not defined in mesg.yml (under configuration.env key)",
		strings.Join(e.env, ", "))
}
