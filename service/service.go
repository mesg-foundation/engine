package service

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xos"
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

	// statuses receives status messages produced during deployment.
	statuses chan DeployStatus `hash:"-"`

	// tempPath is the temporary path that service is hosted in file system.
	tempPath string `hash:"-"`

	// container is the underlying container API.
	container container.Container `hash:"-"`
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

// New creates a new service from a gzipped tarball.
func New(tarball io.Reader, env map[string]string, options ...Option) (*Service, error) {
	s := &Service{}

	defer s.closeStatus()

	if err := s.setOptions(options...); err != nil {
		return nil, err
	}
	if err := s.saveContext(tarball); err != nil {
		return nil, err
	}
	defer os.RemoveAll(s.tempPath)

	// calculate hash
	h := sha1.New()
	if _, err := io.Copy(h, tarball); err != nil {
		return nil, err
	}
	hash := h.Sum(nil)
	s.Hash = fmt.Sprintf("%x", hash)
	fmt.Println("s.Hash", s.Hash)

	def, err := importer.From(s.tempPath)
	if err != nil {
		return nil, err
	}

	s.injectDefinition(def)

	if err := s.validateConfigurationEnv(env); err != nil {
		return nil, err
	}

	// replace default env with new one.
	defenv := xos.EnvSliceToMap(s.configuration().Env)
	s.configuration().Env = xos.EnvMapToSlice(xos.EnvMergeMaps(defenv, env))

	if err := s.deploy(); err != nil {
		return nil, err
	}

	return s.fromService(), nil
}

// FromService upgrades service s by setting its options and private fields.
func FromService(s *Service, options ...Option) (*Service, error) {
	if err := s.setOptions(options...); err != nil {
		return nil, err
	}
	return s.fromService(), nil
}

func (s *Service) setOptions(options ...Option) error {
	for _, option := range options {
		option(s)
	}
	return nil
}

// fromService upgrades service s by setting a calculated ID and cross-referencing its child fields.
// TODO: this function should be deleted.
func (s *Service) fromService() *Service {
	for _, dep := range s.Dependencies {
		dep.service = s
	}

	// s.Hash = s.computeHash()
	return s
}

// Option is the configuration func of Service.
type Option func(*Service)

// ContainerOption returns an option for customized container.
func ContainerOption(container container.Container) Option {
	return func(s *Service) {
		s.container = container
	}
}

// DeployStatusOption receives chan statuses to send deploy statuses.
func DeployStatusOption(statuses chan DeployStatus) Option {
	return func(s *Service) {
		s.statuses = statuses
	}
}

// saveContext downloads service context to a temp dir.
func (s *Service) saveContext(r io.Reader) error {
	var err error
	s.tempPath, err = ioutil.TempDir("", "mesg-")
	if err != nil {
		return err
	}

	s.sendStatus("Receiving service context...", DRunning)
	defer s.sendStatus("Service context received with success", DDonePositive)

	if err := archive.Untar(r, s.tempPath, &archive.TarOptions{
		NoLchown: true,
	}); err != nil {
		return err
	}

	// NOTE: this is check for tar repos, if there is only one
	// directory inside untar archive set temp path to it.
	dirs, err := ioutil.ReadDir(s.tempPath)
	if err != nil {
		return err
	}
	if len(dirs) == 1 && dirs[0].IsDir() {
		s.tempPath = filepath.Join(s.tempPath, dirs[0].Name())
	}
	return nil
}

// deploy deploys service.
func (s *Service) deploy() error {
	s.sendStatus("Building Docker image...", DRunning)

	imageHash, err := s.container.Build(s.tempPath)
	if err != nil {
		return err
	}

	s.sendStatus("Image built with success", DDonePositive)

	s.configuration().Image = imageHash
	// TODO: FOR TEST ONLY
	if s.Hash == "" {
		panic("s.Hash is empty")
	}
	if s.Sid == "" {
		// make sure that sid doesn't have the same length with id.
		s.Sid = "_" + s.Hash
	}
	return nil
}

// sendStatus sends a status message.
func (s *Service) sendStatus(message string, typ DStatusType) {
	if s.statuses != nil {
		s.statuses <- DeployStatus{
			Message: message,
			Type:    typ,
		}
	}
}

// closeStatus closes statuses chan.
func (s *Service) closeStatus() {
	if s.statuses != nil {
		close(s.statuses)
	}
}

// getDependency returns dependency dependencyKey or a not found error.
func (s *Service) getDependency(dependencyKey string) (*Dependency, error) {
	for _, dep := range s.Dependencies {
		if dep.Key == dependencyKey {
			dep.service = s
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
