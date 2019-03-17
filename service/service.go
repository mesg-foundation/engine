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
	Name string `yaml:"name" json:"name,omitempty" validate:"required,printascii,min=1"`

	// Sid is the service id. It must be unique.
	Sid string `yaml:"sid" json:"sid,omitempty" validate:"omitempty,printascii,max=63,domain"`

	// Description is service description.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Repository holds the service's repository url if it's living on a git host.
	Repository string `yaml:"repository" json:"repository,omitempty" validate:"omitempty,uri"`

	// Tasks are the list of tasks that service can execute.
	Tasks map[string]*Task `yaml:"tasks" json:"tasks,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Events are the list of events that service can emit.
	Events map[string]*Event `yaml:"events" json:"events,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Configuration is the Docker container that service runs inside.
	Configuration Dependency `yaml:"configuration" json:"configuration,omitempty"`

	// Dependencies are the Docker containers that service can depend on.
	Dependencies map[string]*Dependency `yaml:"dependencies" json:"dependencies,omitempty" validate:"dive,keys,printascii,ne=service,endkeys,required"`

	// Hash is calculated from the combination of service's source and mesg.yml.
	Hash string `yaml:"-"`

	// Status is service status.
	Status Status `yaml:"-"`

	// DeployedAt holds the creation time of service.
	DeployedAt time.Time `yaml:"-"`
}

// ValidateEventData checks if event data is valid for given event key.
func (s *Service) ValidateEventData(eventKey string, eventData map[string]interface{}) error {
	event, ok := s.Events[eventKey]
	if !ok {
		return fmt.Errorf("service %s - event %q not found", s.Name, eventKey)
	}
	return validateParametersSchema(event.Data, eventData)
}

// ValidateTaskInputs checks if task inputs is valid for given task key.
func (s *Service) ValidateTaskInputs(taskKey string, inputs map[string]interface{}) error {
	task, ok := s.Tasks[taskKey]
	if !ok {
		return fmt.Errorf("task %s not found", taskKey)
	}
	return validateParametersSchema(task.Inputs, inputs)
}

// ValidateTaskOutput checks if task output is valid for given task and output key.
func (s *Service) ValidateTaskOutput(taskKey, outputKey string, outputData map[string]interface{}) error {
	task, ok := s.Tasks[taskKey]
	if !ok {
		return fmt.Errorf("task %s not found", taskKey)
	}
	output, ok := task.Outputs[outputKey]
	if !ok {
		return fmt.Errorf("task %s output not found", taskKey)
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

func (s *Service) ports(depKey string) []container.Port {
	var (
		ports []container.Port
		dep   = &s.Configuration
	)
	if depKey != MainServiceKey {
		dep = s.Dependencies[depKey]
	}
	for _, port := range dep.Ports {
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
	var (
		volumes []container.Mount
		dep     = &s.Configuration
	)
	if depKey != MainServiceKey {
		dep = s.Dependencies[depKey]
	}
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
		depVolumes := s.Configuration.Volumes
		if key != MainServiceKey {
			depVolumes = s.Dependencies[key].Volumes
		}

		for _, volume := range depVolumes {
			volumes = append(volumes, container.Mount{
				Source: volumeKey(service, depKey, volume),
				Target: volume,
			})
		}
	}
	return volumes
}

// Event describes a service task.
type Event struct {
	// Name is the name of event.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of event.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Data holds the input inputs of event.
	Data map[string]*Parameter `yaml:"data" json:"data,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
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
	Env []string `yaml:"env" json:"env,omitempty" validate:"unique,dive,printascii"`
}

// Task describes a service task.
type Task struct {
	// Name is the name of task.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of task.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Parameters are the definition of the execution inputs of task.
	Inputs map[string]*Parameter `yaml:"inputs" json:"inputs,omitempty" validate:"dive,keys,printascii,endkeys,required"`

	// Outputs are the definition of the execution results of task.
	Outputs map[string]*Output `yaml:"outputs" json:"outputs,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
}

// Output describes task output.
type Output struct {
	// Name is the name of task output.
	Name string `yaml:"name" json:"name,omitempty" validate:"printascii"`

	// Description is the description of task output.
	Description string `yaml:"description" json:"description,omitempty" validate:"printascii"`

	// Data holds the output inputs of a task output.
	Data map[string]*Parameter `yaml:"data" json:"data,omitempty" validate:"required,dive,keys,printascii,endkeys,required"`
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
	Object map[string]*Parameter `yaml:"object" json:"object,omitempty" validate:"dive,keys,printascii,endkeys,required"`
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
func validateParametersSchema(parameters map[string]*Parameter, args map[string]interface{}) error {
	for key, param := range parameters {
		if err := param.Validate(args[key]); err != nil {
			return fmt.Errorf("argument %s is %s", key, err)
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
