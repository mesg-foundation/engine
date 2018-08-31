package service

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cnf/structhash"
	"github.com/docker/docker/pkg/archive"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/service/importer"
	uuid "github.com/satori/go.uuid"
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
	// ID is the unique id of service.
	ID string `hash:"-"`

	// Name is the service name.
	Name string `hash:"name:1"`

	// Description is service description.
	Description string `hash:"name:2"`

	// Tasks are the list of tasks that service can execute.
	Tasks []*Task `hash:"name:3"`

	// Events are the list of events that service can emit.
	Events []*Event `hash:"name:4"`

	// Configuration is the Docker container that service runs inside.
	configuration *Dependency `hash:"-"`

	// Dependencies are the Docker containers that service can depend on.
	Dependencies []*Dependency `hash:"name:5"`

	// Repository holds the service's repository url if it's living on
	// a Git host.
	Repository string `hash:"name:6"`

	// DeployedAt holds the creation time of service.
	DeployedAt time.Time `hash:"-"`

	// statuses receives status messages produced during deployment.
	statuses chan DeployStatus `hash:"-"`

	// tempPath is the temporary path that service is hosted in file system.
	tempPath string `hash:"-"`

	// docker is the underlying Docker API.
	docker *container.Container `hash:"-"`
}

// DStatusType indicates the type of status message.
type DStatusType int

const (
	// DRUNNING indicates that status message belongs to an active state.
	DRUNNING DStatusType = iota + 1

	// DDONE indicates that status message belongs to completed state.
	DDONE
)

// DeployStatus represents the deployment status.
type DeployStatus struct {
	Message string
	Type    DStatusType
}

// New creates a new service from a gzipped tarball.
func New(tarball io.Reader, options ...Option) (*Service, error) {
	s := &Service{}
	if err := s.setOptions(options...); err != nil {
		return nil, err
	}
	if err := s.saveContext(tarball); err != nil {
		return nil, err
	}
	if err := s.checkDeprecations(); err != nil {
		return nil, err
	}

	def, err := importer.From(s.tempPath)
	if err != nil {
		os.RemoveAll(s.tempPath)
		return nil, err
	}
	s.injectDefinition(def)

	if err := s.deploy(); err != nil {
		return nil, err
	}
	return s.fromService(), nil
}

// FromService initializes a Service type from s.
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

func (s *Service) fromService() *Service {
	for _, event := range s.Events {
		event.service = s
	}
	for _, task := range s.Tasks {
		task.service = s
		for _, output := range task.Outputs {
			output.task = task
			output.service = s
		}
	}
	for _, dep := range s.Dependencies {
		dep.service = s
	}

	s.ID = s.computeHash()
	return s
}

// Option is the configuration func of Service.
type Option func(*Service)

// ContainerOption returns an option for customized container.
func ContainerOption(container *container.Container) Option {
	return func(s *Service) {
		s.docker = container
	}
}

// DeployStatusOption receives chan statuses to send deploy statuses.
func DeployStatusOption(statuses chan DeployStatus) Option {
	return func(s *Service) {
		s.statuses = statuses
	}
}

// computeHash computes a unique sha1 value for service.
// changes on the names of constant configuration fields of mesg.yml will not
// have effect on computation but extending or removing configurations or changing
// values in mesg.yml will cause computeHash to generate a different value.
func (s *Service) computeHash() string {
	h := sha1.New()
	h.Write(structhash.Dump(s, 1))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Service) saveContext(r io.Reader) error {
	var err error
	s.tempPath, err = s.createTempDir()
	if err != nil {
		return err
	}

	s.sendStatus("Receiving service context...", DRUNNING)
	defer s.sendStatus(fmt.Sprintf("%s Service context received with success.", aurora.Green("✔")), DDONE)

	return archive.Untar(r, s.tempPath, &archive.TarOptions{
		Compression: archive.Gzip,
		NoLchown:    true,
	})
}

func (s *Service) createTempDir() (path string, err error) {
	return ioutil.TempDir("", "mesg-"+uuid.NewV4().String())
}

func (s *Service) injectDefinition(def *importer.ServiceDefinition) {
	s.Name = def.Name
	s.Description = def.Description
	s.Repository = def.Repository
	s.Events = []*Event{}
	s.Tasks = []*Task{}
	s.Dependencies = []*Dependency{}

	for key, dep := range def.Dependencies {
		s.Dependencies = append(s.Dependencies, &Dependency{
			Key:         key,
			Image:       dep.Image,
			Volumes:     dep.Volumes,
			VolumesFrom: dep.VolumesFrom,
			Ports:       dep.Ports,
			Command:     dep.Command,
		})
	}

	for key, event := range def.Events {
		s.Events = append(s.Events, &Event{
			Key:         key,
			Name:        event.Name,
			Description: event.Description,
			Data:        s.defParametersToServiceParameters(event.Data),
		})
	}

	for key, task := range def.Tasks {
		t := &Task{
			Key:         key,
			Name:        def.Name,
			Description: def.Description,
			Inputs:      s.defParametersToServiceParameters(task.Inputs),
			Outputs:     []*Output{},
		}

		for key, output := range task.Outputs {
			t.Outputs = append(t.Outputs, &Output{
				Key:         key,
				Name:        output.Name,
				Description: output.Description,
				Data:        s.defParametersToServiceParameters(output.Data),
			})
		}

		s.Tasks = append(s.Tasks, t)
	}

	s.configuration = &Dependency{}
	if def.Configuration != nil {
		s.configuration.Command = def.Configuration.Command
		s.configuration.Ports = def.Configuration.Ports
		s.configuration.Volumes = def.Configuration.Volumes
		s.configuration.VolumesFrom = def.Configuration.VolumesFrom
	}
}

func (s *Service) defParametersToServiceParameters(params map[string]*importer.Parameter) []*Parameter {
	ps := []*Parameter{}
	for key, param := range params {
		ps = append(ps, &Parameter{
			Key:         key,
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Optional:    param.Optional,
		})
	}
	return ps
}

func (s *Service) deploy() error {
	defer os.RemoveAll(s.tempPath)
	defer s.closeStatusSend()

	s.sendStatus("Building Docker image...", DRUNNING)

	imageHash, err := s.docker.Build(s.tempPath)
	if err != nil {
		return err
	}

	s.sendStatus(fmt.Sprintf("%s Image built with success.", aurora.Green("✔")), DDONE)
	s.sendStatus(fmt.Sprintf("%s Completed.", aurora.Green("✔")), DDONE)

	s.configuration.Key = "service"
	s.configuration.Image = imageHash
	s.Dependencies = append(s.Dependencies, s.configuration)
	return nil
}

func (s *Service) checkDeprecations() error {
	if _, err := os.Stat(filepath.Join(s.tempPath, ".mesgignore")); err == nil {
		// TODO: remove for a future release
		s.sendStatus(fmt.Sprintf("%s [DEPRECATED] Please use .dockerignore instead of .mesgignore",
			aurora.Red("⨯")), DDONE)
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

// closeStatusSend closes status send chan.
func (s *Service) closeStatusSend() {
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
	return nil, fmt.Errorf("Dependency %s do not exist", dependencyKey)
}

// containerNotProvidedError returned when a *container.Container instace not
// provided during Service initialization.
type containerNotProvidedError struct{}

func (e *containerNotProvidedError) Error() string {
	return "container instance not provided"
}
