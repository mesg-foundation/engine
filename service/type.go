package service

import proto "github.com/golang/protobuf/proto"

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

// List of all publishs flags
const (
	PublishAll       Publish = "ALL"
	PublishSource    Publish = "SOURCE"
	PublishContainer Publish = "CONTAINER"
	PublishNone      Publish = "NONE"
)

// Service is a definition for a service to run
type Service struct {
	Name         string                `yaml:"name" json:"name" protobuf:"bytes,1,opt,name=name"`
	Description  string                `yaml:"description" json:"description" protobuf:"bytes,2,opt,name=description"`
	Visibility   Visibility            `yaml:"visibility" json:"visibility" protobuf:"bytes,3,opt,name=visibility"`
	Publish      Publish               `yaml:"publish" json:"publish" protobuf:"bytes,4,opt,name=publish"`
	Tasks        map[string]Task       `yaml:"tasks" json:"tasks" protobuf:"bytes,5,rep,name=tasks" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Events       map[string]Event      `yaml:"events" json:"events" protobuf:"bytes,6,rep,name=events" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Dependencies map[string]Dependency `yaml:"dependencies" json:"dependencies" protobuf:"bytes,7,rep,name=dependencies" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Service) Reset()         { *m = Service{} }
func (m *Service) String() string { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()    {}

// Task is a definition of a Task from a service
type Task struct {
	Name        string               `yaml:"name" json:"name" protobuf:"bytes,1,opt,name=name"`
	Description string               `yaml:"description" json:"description" protobuf:"bytes,2,opt,name=description"`
	Verifiable  bool                 `yaml:"verifiable" json:"verifiable" protobuf:"varint,3,opt,name=verifiable"`
	Payable     bool                 `yaml:"payable" json:"payable" protobuf:"varint,4,opt,name=payable"`
	Fees        Fees                 `yaml:"fees" json:"fees" protobuf:"bytes,5,opt,name=fees"`
	Inputs      map[string]Parameter `yaml:"inputs" json:"inputs" protobuf:"bytes,6,rep,name=inputs" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Outputs     map[string]Event     `yaml:"outputs" json:"outputs" protobuf:"bytes,7,rep,name=outputs" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}

// Fees is the different fees to apply
type Fees struct {
	Developer string `yaml:"developer" json:"developer" protobuf:"bytes,1,opt,name=developer"`
	Validator string `yaml:"validator" json:"validator" protobuf:"bytes,2,opt,name=validator"`
	Executor  string `yaml:"executor" json:"executor" protobuf:"bytes,3,opt,name=executor"`
	Emitters  string `yaml:"emitters" json:"emitters" protobuf:"bytes,4,opt,name=emitters"`
}

func (m *Fees) Reset()         { *m = Fees{} }
func (m *Fees) String() string { return proto.CompactTextString(m) }
func (*Fees) ProtoMessage()    {}

// Event is the definition of an event emitted from a service
type Event struct {
	Name        string               `yaml:"name" json:"name" protobuf:"bytes,1,opt,name=name"`
	Description string               `yaml:"description" json:"description" protobuf:"bytes,2,opt,name=description"`
	Data        map[string]Parameter `yaml:"data" json:"data" protobuf:"bytes,3,rep,name=data" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}

// Parameter is the definition of a parameter for a Task
type Parameter struct {
	Name        string `yaml:"name" json:"name" protobuf:"bytes,1,opt,name=name"`
	Description string `yaml:"description" json:"description" protobuf:"bytes,2,opt,name=description"`
	Type        string `yaml:"type" json:"type" protobuf:"bytes,3,opt,name=type"`
	Optional    bool   `yaml:"optional" json:"optional" protobuf:"varint,4,opt,name=optional"`
}

func (m *Parameter) Reset()         { *m = Parameter{} }
func (m *Parameter) String() string { return proto.CompactTextString(m) }
func (*Parameter) ProtoMessage()    {}

// Dependency is the docker informations about the Dependency
type Dependency struct {
	Image   string   `yaml:"image" json:"image" protobuf:"bytes,1,opt,name=image"`
	Volumes []string `yaml:"volumes" json:"volumes" protobuf:"bytes,2,rep,name=image"`
	Ports   []string `yaml:"ports" json:"ports" protobuf:"bytes,3,rep,name=ports"`
	Command string   `yaml:"command" json:"command" protobuf:"bytes,4,opt,name=command"`
}

func (m *Dependency) Reset()         { *m = Dependency{} }
func (m *Dependency) String() string { return proto.CompactTextString(m) }
func (*Dependency) ProtoMessage()    {}
