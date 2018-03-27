package service

// Visibility is the tags to set is the service is visible for whom
type Visibility string

// List of visibilities flags
const (
	VisibilityAll     Visibility = "ALL"
	VisibilityUsers   Visibility = "USERS"
	VisibilityWorkers Visibility = "WORKERS"
	VisibilityNone    Visibility = "NONE"
)

// FeeActor is a type to get the different actors in the fee system
type FeeActor string

// List of the different actors in the fees
const (
	Developer FeeActor = "developer"
	Emitters  FeeActor = "emitters"
	Validator FeeActor = "validator"
	Executor  FeeActor = "executor"
)

// Publish let you configure the part of your service you want to publish
type Publish string

// List of all publishs flags
const (
	PublishAll       Publish = "ALL"
	PublishSource    Publish = "SOURCE"
	PublishContainer Publish = "CONTAINER"
	PublishNone      Publish = "NONE"
)

// Service is a definition for a service to run
type Service struct {
	Name         string       `yaml:"name" json:"name"`
	Description  string       `yaml:"description" json:"description"`
	Visibility   Visibility   `yaml:"visibility" json:"visibility"`
	Publish      Publish      `yaml:"publish" json:"publish"`
	Tasks        Tasks        `yaml:"tasks" json:"tasks"`
	Events       Events       `yaml:"events" json:"events"`
	Dependencies Dependencies `yaml:"dependencies" json:"dependencies"`
}

// Task is a definition of a Task from a service
type Task struct {
	Name        string     `yaml:"name" json:"name"`
	Description string     `yaml:"description" json:"description"`
	Verifiable  bool       `yaml:"verifiable" json:"verifiable"`
	Payable     bool       `yaml:"payable" json:"payable"`
	Fees        Fees       `yaml:"fees" json:"fees"`
	Inputs      Parameters `yaml:"inputs" json:"inputs"`
	Outputs     Events     `yaml:"outputs" json:"outputs"`
}

// Fees is the different fees to apply
type Fees struct {
	Developer FeeActor `yaml:"developer" json:"developer"`
	Validator FeeActor `yaml:"validator" json:"validator"`
	Executor  FeeActor `yaml:"executor" json:"executor"`
	Emitters  FeeActor `yaml:"emitters" json:"emitters"`
}

// Tasks is a list of Tasks
type Tasks map[string]Task

// Event is the definition of an event emitted from a service
type Event struct {
	Name        string     `yaml:"name" json:"name"`
	Description string     `yaml:"description" json:"description"`
	Data        Parameters `yaml:"data" json:"data"`
}

// Events is a list of Events
type Events map[string]Event

// Parameter is the definition of a parameter for a Task
type Parameter struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Type        string `yaml:"type" json:"type"`
	Optional    bool   `yaml:"optional" json:"optional"`
}

// Parameters is a list of Parameters
type Parameters map[string]Parameter

// Dependency is the docker informations about the Dependency
type Dependency struct {
	Image   string   `yaml:"image" json:"image"`
	Volumes []string `yaml:"volumes" json:"volumes"`
	Ports   []string `yaml:"ports" json:"ports"`
	Command string   `yaml:"command" json:"command"`
}

// Dependencies is a list of Dependencies
type Dependencies map[string]Dependency
