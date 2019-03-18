package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mesg-foundation/core/container"
)

// MainServiceKey is key for main service.
const MainServiceKey = "service"

// Namespacees used for the docker services.
const (
	eventChannel  = "Event"
	taskChannel   = "Task"
	resultChannel = "Result"
)

// Status of the service.
type Status uint

// Possible statuses for service.
const (
	StatusUnknown Status = iota
	StatusStopped
	StatusStarting
	StatusPartial
	StatusRunning
	StatusDeleted
)

func (s Status) String() string {
	switch s {
	case StatusStopped:
		return "STOPPED"
	case StatusStarting:
		return "STARTING"
	case StatusPartial:
		return "PARTIAL"
	case StatusRunning:
		return "RUNNING"
	case StatusDeleted:
		return "DELETED"
	default:
		return "UNKNOWN"
	}
}

// Service represents MESG services configurations.
type Service struct {
	// Name is the service name.
	Name string

	// Sid is the service id. It must be unique.
	Sid string

	// Description is service description.
	Description string

	// Repository holds the service's repository url if it's living on a git host.
	Repository string

	// Tasks are the list of tasks that service can execute.
	Tasks []*Task

	// Events are the list of events that service can emit.
	Events []*Event

	// Configuration is the Docker container that service runs inside.
	Configuration *Dependency

	// Dependencies are the Docker containers that service can depend on.
	Dependencies []*Dependency

	// Hash is calculated from the combination of service's source and mesg.yml.
	Hash string

	// Status is service status.
	Status Status

	// DeployedAt holds the creation time of service.
	DeployedAt time.Time
}

// Dependency returns dependency for given key or error if not found.
func (s *Service) Dependency(key string) (*Dependency, error) {
	if key == MainServiceKey {
		return s.Configuration, nil
	}

	for i := range s.Dependencies {
		if s.Dependencies[i].Key == key {
			return s.Dependencies[i], nil
		}
	}
	return nil, fmt.Errorf("service %q - dependency %q not found", s.Name, key)
}

// Event returns event for given key or error if not found.
func (s *Service) Event(key string) (*Event, error) {
	for i := range s.Events {
		if s.Events[i].Key == key {
			return s.Events[i], nil
		}
	}
	return nil, fmt.Errorf("service %q - event %q not found", s.Name, key)
}

// Task returns task for given key or error if not found.
func (s *Service) Task(key string) (*Task, error) {
	for i := range s.Tasks {
		if s.Tasks[i].Key == key {
			return s.Tasks[i], nil
		}
	}
	return nil, fmt.Errorf("service %q - task %q not found", s.Name, key)
}

// ValidateEventData checks if event data is valid for given event key.
func (s *Service) ValidateEventData(eventKey string, eventData map[string]interface{}) error {
	event, err := s.Event(eventKey)
	if err != nil {
		return err
	}
	return validateParametersSchema(event.Data, eventData)
}

// ValidateTaskInputs checks if task inputs is valid for given task key.
func (s *Service) ValidateTaskInputs(taskKey string, inputs map[string]interface{}) error {
	task, err := s.Task(taskKey)
	if err != nil {
		return err
	}
	return validateParametersSchema(task.Inputs, inputs)
}

// ValidateTaskOutput checks if task output is valid for given task and output key.
func (s *Service) ValidateTaskOutput(taskKey, outputKey string, outputData map[string]interface{}) error {
	task, err := s.Task(taskKey)
	if err != nil {
		return err
	}
	output, err := task.Output(outputKey)
	if err != nil {
		return err
	}
	return validateParametersSchema(output.Data, outputData)
}

// validateConfigurationEnv checks presence of env variables in mesg.yml under env section.
func (s *Service) validateConfigurationEnv(env map[string]string) error {
	var notenv []string
	for key := range env {
		exists := false
		// check if "key=" exists in configuration
		for _, env := range s.Configuration.Env {
			if strings.HasPrefix(env, key+"=") {
				exists = true
				break
			}
		}
		if !exists {
			notenv = append(notenv, key)
		}
	}
	if len(notenv) > 0 {
		return fmt.Errorf("environment variable(s) %q not defined in mesg.yml (under configuration.env key)",
			strings.Join(notenv, ", "))
	}
	return nil
}

// namespace returns the namespace of the service.
func (s *Service) namespace() []string {
	sum := sha1.Sum([]byte(s.Hash))
	return []string{hex.EncodeToString(sum[:])}
}

// EventSubscriptionChannel returns the channel to listen for events from this service.
func (s *Service) EventSubscriptionChannel() string {
	return calculate(append(s.namespace(), eventChannel))
}

// TaskSubscriptionChannel returns the channel to listen for tasks from this service.
func (s *Service) TaskSubscriptionChannel() string {
	return calculate(append(s.namespace(), taskChannel))
}

// ResultSubscriptionChannel returns the channel to listen for tasks from this service.
func (s *Service) ResultSubscriptionChannel() string {
	return calculate(append(s.namespace(), resultChannel))
}

func (d *Dependency) ports() []container.Port {
	var ports []container.Port
	for _, port := range d.Ports {
		parts := strings.Split(port, ":")
		published, _ := strconv.ParseUint(parts[0], 10, 64)
		target := published
		if len(parts) > 1 {
			target, _ = strconv.ParseUint(parts[1], 10, 64)
		}
		ports = append(ports, container.Port{
			Target:    uint32(target),
			Published: uint32(published),
		})
	}
	return ports
}

func (s *Service) volumes(depKey string) []container.Mount {
	var volumes []container.Mount
	dep, _ := s.Dependency(depKey)
	service := s.Sid
	if s.Sid == "" {
		service = s.Hash
	}

	for _, volume := range dep.Volumes {
		volumes = append(volumes, container.Mount{
			Source: volumeKey(service, depKey, volume),
			Target: volume,
		})
	}

	for _, key := range dep.VolumesFrom {
		vdep, _ := s.Dependency(key)
		for _, volume := range vdep.Volumes {
			volumes = append(volumes, container.Mount{
				Source: volumeKey(service, key, volume),
				Target: volume,
			})
		}
	}
	return volumes
}

// Event describes a service task.
type Event struct {
	Key string

	// Name is the name of event.
	Name string

	// Description is the description of event.
	Description string

	// Data holds the input inputs of event.
	Data []*Parameter
}

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	Key string

	// Image is the Docker image.
	Image string

	// Volumes are the Docker volumes.
	Volumes []string

	// VolumesFrom are the docker volumes-from from.
	VolumesFrom []string

	// Ports holds ports configuration for container.
	Ports []string

	// Command is the Docker command which will be executed when container started.
	Command string

	// Args hold the args to pass to the Docker container
	Args []string

	// Env is a slice of environment variables in key=value format.
	Env []string
}

// Task describes a service task.
type Task struct {
	Key string

	// Name is the name of task.
	Name string

	// Description is the description of task.
	Description string

	// Parameters are the definition of the execution inputs of task.
	Inputs []*Parameter

	// Outputs are the definition of the execution results of task.
	Outputs []*Output
}

// Output returns output for given key or error if not found.
func (t *Task) Output(key string) (*Output, error) {
	for i := range t.Outputs {
		if t.Outputs[i].Key == key {
			return t.Outputs[i], nil
		}
	}
	return nil, fmt.Errorf("task %q - output %q not found", t.Key, key)
}

// Output describes task output.
type Output struct {
	Key string

	// Name is the name of task output.
	Name string

	// Description is the description of task output.
	Description string

	// Data holds the output inputs of a task output.
	Data []*Parameter
}

// Parameter describes task input inputs, output inputs of a task
// output and input inputs of an event.
type Parameter struct {
	Key string

	// Name is the name of input.
	Name string

	// Description is the description of input.
	Description string

	// Type is the data type of input.
	Type string

	// Optional indicates if input is optional.
	Optional bool

	// Repeated is to have an array of this input
	Repeated bool

	// Definition of the structure of the object when the type is object
	Object []*Parameter
}

// Validate validates arg by comparing to its parameter schema.
func (p *Parameter) Validate(arg interface{}) error {
	if arg == nil {
		if p.Optional {
			return nil
		}
		return errors.New("required")
	}
	if p.Repeated {
		array, ok := arg.([]interface{})
		if !ok {
			return errors.New("not an array")
		}
		for _, x := range array {
			if err := p.validateType(x); err != nil {
				return err
			}
		}
		return nil
	}
	return p.validateType(arg)
}

const (
	paramStringType  = "String"
	paramNumberType  = "Number"
	paramBooleanType = "Boolean"
	paramObjectType  = "Object"
	paramAnyType     = "Any"
)

// validateType checks if arg comforts its expected type.
func (p *Parameter) validateType(arg interface{}) error {
	switch p.Type {
	case paramStringType:
		if _, ok := arg.(string); !ok {
			return errors.New("not a string")
		}
	case paramNumberType:
		if _, ok := arg.(float64); !ok {
			return errors.New("not a number")
		}
	case paramBooleanType:
		if _, ok := arg.(bool); !ok {
			return errors.New("not a boolean")
		}
	case paramObjectType:
		obj, ok := arg.(map[string]interface{})
		if !ok {
			return errors.New("not an object")
		}
		return validateParametersSchema(p.Object, obj)
	case paramAnyType:
		return nil
	default:
		return errors.New("not valid type")
	}
	return nil
}

// calculate returns a hash according to the given data.
func calculate(data []string) string {
	return strings.Join(data, ".")
}

// ValidateParametersSchema validates data to see if it matches with parameters schema.
func validateParametersSchema(parameters []*Parameter, args map[string]interface{}) error {
	for _, param := range parameters {
		if err := param.Validate(args[param.Key]); err != nil {
			return fmt.Errorf("argument %q is %s", param.Key, err)
		}
	}
	return nil
}

// volumeKey creates a key for service's volume based on the service sid or hash.
func volumeKey(service, dependency, volume string) string {
	h := sha1.New()
	h.Write([]byte(service))
	h.Write([]byte(dependency))
	sum := h.Sum([]byte(volume))
	return hex.EncodeToString(sum)
}
