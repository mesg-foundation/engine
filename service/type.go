package service

// Visibility is the tags to set is the service is visible for whom
type Visibility string

// List of visibilities flags
const (
	vALL     Visibility = "ALL"
	vUSERS   Visibility = "USERS"
	vWORKERS Visibility = "WORKERS"
	vNONE    Visibility = "NONE"
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

// Publication let you configure the part of your service you want to publish
type Publication string

// List of all publications flags
const (
	AllPublication       Publication = "ALL"
	SourcePublication    Publication = "SOURCE"
	ContainerPublication Publication = "CONTAINER"
	NonePublication      Publication = "NONE"
)

// Service is a definition for a service to run
type Service struct {
	Version      string       `yaml:"version"`
	Name         string       `yaml:"name"`
	Description  string       `yaml:"description"`
	Requirements Requirements `yaml:"requirements"`
	Visibility   Visibility   `yaml:"visibility"`
	Publication  Publication  `yaml:"publication"`
	Tasks        Tasks        `yaml:"tasks"`
	Events       Events       `yaml:"events"`
	Dependencies Dependencies `yaml:"services"`
}

// Requirements is a list of requirements to run the service
type Requirements struct {
	memory uint `yaml:"memory"`
	cpu    uint `yaml:"cpu"`
	disk   uint `yaml:"disk"`
}

// Task is a definition of a Task from a service
type Task struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Verifiable  bool       `yaml:"verifiable"`
	Payable     bool       `yaml:"payable"`
	Fees        Fees       `yaml:"fees"`
	Inputs      Parameters `yaml:"inputs"`
	Outputs     Events     `yaml:"outputs"`
}

// Fees is the different fees to apply
type Fees struct {
	developer FeeActor `yaml:"developer"`
	validator FeeActor `yaml:"validator"`
	executor  FeeActor `yaml:"executor"`
	emitters  FeeActor `yaml:"emitters"`
}

// Tasks is a list of Tasks
type Tasks map[string]Task

// Event is the definition of an event emitted from a service
type Event struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Data        Parameters `yaml:"data"`
}

// Events is a list of Events
type Events map[string]Event

// Parameter is the definition of a parameter for a Task
type Parameter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Optional    bool   `yaml:"optional"`
}

// Parameters is a list of Parameters
type Parameters []Parameter

// Dependency is the docker informations about the Dependency
type Dependency struct {
	Image   string   `yaml:"image"`
	Volumes []string `yaml:"volumes"`
	Ports   []string `yaml:"ports"`
	Command string   `yaml:"command"`
}

// Dependencies is a list of Dependencies
type Dependencies map[string]Dependency
