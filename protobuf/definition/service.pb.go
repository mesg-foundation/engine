// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/definition/service.proto

package definition // import "github.com/mesg-foundation/core/protobuf/definition"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Service represents the service's definition.
type Service struct {
	Hash                 string         `protobuf:"bytes,10,opt,name=hash,proto3" json:"hash,omitempty"`
	Sid                  string         `protobuf:"bytes,12,opt,name=sid,proto3" json:"sid,omitempty"`
	Name                 string         `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description          string         `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Configuration        *Configuration `protobuf:"bytes,8,opt,name=configuration,proto3" json:"configuration,omitempty"`
	Tasks                []*Task        `protobuf:"bytes,5,rep,name=tasks,proto3" json:"tasks,omitempty"`
	Events               []*Event       `protobuf:"bytes,6,rep,name=events,proto3" json:"events,omitempty"`
	Dependencies         []*Dependency  `protobuf:"bytes,7,rep,name=dependencies,proto3" json:"dependencies,omitempty"`
	Repository           string         `protobuf:"bytes,9,opt,name=repository,proto3" json:"repository,omitempty"`
	Source               string         `protobuf:"bytes,13,opt,name=source,proto3" json:"source,omitempty"`
	Workflows            []*Workflow    `protobuf:"bytes,14,rep,name=workflows,proto3" json:"workflows,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Service) Reset()         { *m = Service{} }
func (m *Service) String() string { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()    {}
func (*Service) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{0}
}
func (m *Service) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service.Unmarshal(m, b)
}
func (m *Service) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service.Marshal(b, m, deterministic)
}
func (dst *Service) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service.Merge(dst, src)
}
func (m *Service) XXX_Size() int {
	return xxx_messageInfo_Service.Size(m)
}
func (m *Service) XXX_DiscardUnknown() {
	xxx_messageInfo_Service.DiscardUnknown(m)
}

var xxx_messageInfo_Service proto.InternalMessageInfo

func (m *Service) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *Service) GetSid() string {
	if m != nil {
		return m.Sid
	}
	return ""
}

func (m *Service) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Service) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Service) GetConfiguration() *Configuration {
	if m != nil {
		return m.Configuration
	}
	return nil
}

func (m *Service) GetTasks() []*Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

func (m *Service) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *Service) GetDependencies() []*Dependency {
	if m != nil {
		return m.Dependencies
	}
	return nil
}

func (m *Service) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func (m *Service) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Service) GetWorkflows() []*Workflow {
	if m != nil {
		return m.Workflows
	}
	return nil
}

// Events are emitted by the service whenever the service wants.
type Event struct {
	Key                  string       `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
	Name                 string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description          string       `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Data                 []*Parameter `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{1}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (dst *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(dst, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Event) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Event) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Event) GetData() []*Parameter {
	if m != nil {
		return m.Data
	}
	return nil
}

// Task is a function that requires inputs and returns output.
type Task struct {
	Key                  string       `protobuf:"bytes,8,opt,name=key,proto3" json:"key,omitempty"`
	Name                 string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description          string       `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Inputs               []*Parameter `protobuf:"bytes,6,rep,name=inputs,proto3" json:"inputs,omitempty"`
	Outputs              []*Parameter `protobuf:"bytes,7,rep,name=outputs,proto3" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{2}
}
func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (dst *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(dst, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Task) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Task) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Task) GetInputs() []*Parameter {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *Task) GetOutputs() []*Parameter {
	if m != nil {
		return m.Outputs
	}
	return nil
}

// Parameter describes the task's inputs, the task's outputs, and the event's data.
type Parameter struct {
	Key                  string       `protobuf:"bytes,8,opt,name=key,proto3" json:"key,omitempty"`
	Name                 string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description          string       `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Type                 string       `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Optional             bool         `protobuf:"varint,4,opt,name=optional,proto3" json:"optional,omitempty"`
	Repeated             bool         `protobuf:"varint,9,opt,name=repeated,proto3" json:"repeated,omitempty"`
	Object               []*Parameter `protobuf:"bytes,10,rep,name=object,proto3" json:"object,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Parameter) Reset()         { *m = Parameter{} }
func (m *Parameter) String() string { return proto.CompactTextString(m) }
func (*Parameter) ProtoMessage()    {}
func (*Parameter) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{3}
}
func (m *Parameter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Parameter.Unmarshal(m, b)
}
func (m *Parameter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Parameter.Marshal(b, m, deterministic)
}
func (dst *Parameter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Parameter.Merge(dst, src)
}
func (m *Parameter) XXX_Size() int {
	return xxx_messageInfo_Parameter.Size(m)
}
func (m *Parameter) XXX_DiscardUnknown() {
	xxx_messageInfo_Parameter.DiscardUnknown(m)
}

var xxx_messageInfo_Parameter proto.InternalMessageInfo

func (m *Parameter) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Parameter) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Parameter) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Parameter) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Parameter) GetOptional() bool {
	if m != nil {
		return m.Optional
	}
	return false
}

func (m *Parameter) GetRepeated() bool {
	if m != nil {
		return m.Repeated
	}
	return false
}

func (m *Parameter) GetObject() []*Parameter {
	if m != nil {
		return m.Object
	}
	return nil
}

type Configuration struct {
	Volumes              []string `protobuf:"bytes,1,rep,name=volumes,proto3" json:"volumes,omitempty"`
	VolumesFrom          []string `protobuf:"bytes,2,rep,name=volumesFrom,proto3" json:"volumesFrom,omitempty"`
	Ports                []string `protobuf:"bytes,3,rep,name=ports,proto3" json:"ports,omitempty"`
	Args                 []string `protobuf:"bytes,4,rep,name=args,proto3" json:"args,omitempty"`
	Command              string   `protobuf:"bytes,5,opt,name=command,proto3" json:"command,omitempty"`
	Env                  []string `protobuf:"bytes,6,rep,name=env,proto3" json:"env,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Configuration) Reset()         { *m = Configuration{} }
func (m *Configuration) String() string { return proto.CompactTextString(m) }
func (*Configuration) ProtoMessage()    {}
func (*Configuration) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{4}
}
func (m *Configuration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Configuration.Unmarshal(m, b)
}
func (m *Configuration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Configuration.Marshal(b, m, deterministic)
}
func (dst *Configuration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Configuration.Merge(dst, src)
}
func (m *Configuration) XXX_Size() int {
	return xxx_messageInfo_Configuration.Size(m)
}
func (m *Configuration) XXX_DiscardUnknown() {
	xxx_messageInfo_Configuration.DiscardUnknown(m)
}

var xxx_messageInfo_Configuration proto.InternalMessageInfo

func (m *Configuration) GetVolumes() []string {
	if m != nil {
		return m.Volumes
	}
	return nil
}

func (m *Configuration) GetVolumesFrom() []string {
	if m != nil {
		return m.VolumesFrom
	}
	return nil
}

func (m *Configuration) GetPorts() []string {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *Configuration) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *Configuration) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func (m *Configuration) GetEnv() []string {
	if m != nil {
		return m.Env
	}
	return nil
}

// A dependency is a configuration of an other container that runs separately from the service.
type Dependency struct {
	Key                  string   `protobuf:"bytes,8,opt,name=key,proto3" json:"key,omitempty"`
	Image                string   `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
	Volumes              []string `protobuf:"bytes,2,rep,name=volumes,proto3" json:"volumes,omitempty"`
	VolumesFrom          []string `protobuf:"bytes,3,rep,name=volumesFrom,proto3" json:"volumesFrom,omitempty"`
	Ports                []string `protobuf:"bytes,4,rep,name=ports,proto3" json:"ports,omitempty"`
	Args                 []string `protobuf:"bytes,6,rep,name=args,proto3" json:"args,omitempty"`
	Command              string   `protobuf:"bytes,5,opt,name=command,proto3" json:"command,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Dependency) Reset()         { *m = Dependency{} }
func (m *Dependency) String() string { return proto.CompactTextString(m) }
func (*Dependency) ProtoMessage()    {}
func (*Dependency) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{5}
}
func (m *Dependency) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Dependency.Unmarshal(m, b)
}
func (m *Dependency) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Dependency.Marshal(b, m, deterministic)
}
func (dst *Dependency) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Dependency.Merge(dst, src)
}
func (m *Dependency) XXX_Size() int {
	return xxx_messageInfo_Dependency.Size(m)
}
func (m *Dependency) XXX_DiscardUnknown() {
	xxx_messageInfo_Dependency.DiscardUnknown(m)
}

var xxx_messageInfo_Dependency proto.InternalMessageInfo

func (m *Dependency) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Dependency) GetImage() string {
	if m != nil {
		return m.Image
	}
	return ""
}

func (m *Dependency) GetVolumes() []string {
	if m != nil {
		return m.Volumes
	}
	return nil
}

func (m *Dependency) GetVolumesFrom() []string {
	if m != nil {
		return m.VolumesFrom
	}
	return nil
}

func (m *Dependency) GetPorts() []string {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *Dependency) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *Dependency) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

// Workflow represents the workflow's definition.
type Workflow struct {
	Key                  string            `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Trigger              *Workflow_Trigger `protobuf:"bytes,2,opt,name=trigger,proto3" json:"trigger,omitempty"`
	Tasks                []*Workflow_Task  `protobuf:"bytes,3,rep,name=tasks,proto3" json:"tasks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Workflow) Reset()         { *m = Workflow{} }
func (m *Workflow) String() string { return proto.CompactTextString(m) }
func (*Workflow) ProtoMessage()    {}
func (*Workflow) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{6}
}
func (m *Workflow) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow.Unmarshal(m, b)
}
func (m *Workflow) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow.Marshal(b, m, deterministic)
}
func (dst *Workflow) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow.Merge(dst, src)
}
func (m *Workflow) XXX_Size() int {
	return xxx_messageInfo_Workflow.Size(m)
}
func (m *Workflow) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow proto.InternalMessageInfo

func (m *Workflow) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Workflow) GetTrigger() *Workflow_Trigger {
	if m != nil {
		return m.Trigger
	}
	return nil
}

func (m *Workflow) GetTasks() []*Workflow_Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

// Trigger keeps the conditions to match in order to start workflow.
type Workflow_Trigger struct {
	InstanceHash         string           `protobuf:"bytes,1,opt,name=instanceHash,proto3" json:"instanceHash,omitempty"`
	EventKey             string           `protobuf:"bytes,2,opt,name=eventKey,proto3" json:"eventKey,omitempty"`
	Filter               *Workflow_Filter `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Workflow_Trigger) Reset()         { *m = Workflow_Trigger{} }
func (m *Workflow_Trigger) String() string { return proto.CompactTextString(m) }
func (*Workflow_Trigger) ProtoMessage()    {}
func (*Workflow_Trigger) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{6, 0}
}
func (m *Workflow_Trigger) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Trigger.Unmarshal(m, b)
}
func (m *Workflow_Trigger) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Trigger.Marshal(b, m, deterministic)
}
func (dst *Workflow_Trigger) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Trigger.Merge(dst, src)
}
func (m *Workflow_Trigger) XXX_Size() int {
	return xxx_messageInfo_Workflow_Trigger.Size(m)
}
func (m *Workflow_Trigger) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Trigger.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Trigger proto.InternalMessageInfo

func (m *Workflow_Trigger) GetInstanceHash() string {
	if m != nil {
		return m.InstanceHash
	}
	return ""
}

func (m *Workflow_Trigger) GetEventKey() string {
	if m != nil {
		return m.EventKey
	}
	return ""
}

func (m *Workflow_Trigger) GetFilter() *Workflow_Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

// Filter keeps additional conditions for trigger.
type Workflow_Filter struct {
	TaskKey              string   `protobuf:"bytes,1,opt,name=taskKey,proto3" json:"taskKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Workflow_Filter) Reset()         { *m = Workflow_Filter{} }
func (m *Workflow_Filter) String() string { return proto.CompactTextString(m) }
func (*Workflow_Filter) ProtoMessage()    {}
func (*Workflow_Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{6, 1}
}
func (m *Workflow_Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Filter.Unmarshal(m, b)
}
func (m *Workflow_Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Filter.Marshal(b, m, deterministic)
}
func (dst *Workflow_Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Filter.Merge(dst, src)
}
func (m *Workflow_Filter) XXX_Size() int {
	return xxx_messageInfo_Workflow_Filter.Size(m)
}
func (m *Workflow_Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Filter proto.InternalMessageInfo

func (m *Workflow_Filter) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

// Task represents a task to execute when conditions are met.
type Workflow_Task struct {
	InstanceHash         string   `protobuf:"bytes,1,opt,name=instanceHash,proto3" json:"instanceHash,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Workflow_Task) Reset()         { *m = Workflow_Task{} }
func (m *Workflow_Task) String() string { return proto.CompactTextString(m) }
func (*Workflow_Task) ProtoMessage()    {}
func (*Workflow_Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_7c01ece46eb6a4a2, []int{6, 2}
}
func (m *Workflow_Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Task.Unmarshal(m, b)
}
func (m *Workflow_Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Task.Marshal(b, m, deterministic)
}
func (dst *Workflow_Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Task.Merge(dst, src)
}
func (m *Workflow_Task) XXX_Size() int {
	return xxx_messageInfo_Workflow_Task.Size(m)
}
func (m *Workflow_Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Task proto.InternalMessageInfo

func (m *Workflow_Task) GetInstanceHash() string {
	if m != nil {
		return m.InstanceHash
	}
	return ""
}

func (m *Workflow_Task) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func init() {
	proto.RegisterType((*Service)(nil), "definition.Service")
	proto.RegisterType((*Event)(nil), "definition.Event")
	proto.RegisterType((*Task)(nil), "definition.Task")
	proto.RegisterType((*Parameter)(nil), "definition.Parameter")
	proto.RegisterType((*Configuration)(nil), "definition.Configuration")
	proto.RegisterType((*Dependency)(nil), "definition.Dependency")
	proto.RegisterType((*Workflow)(nil), "definition.Workflow")
	proto.RegisterType((*Workflow_Trigger)(nil), "definition.Workflow.Trigger")
	proto.RegisterType((*Workflow_Filter)(nil), "definition.Workflow.Filter")
	proto.RegisterType((*Workflow_Task)(nil), "definition.Workflow.Task")
}

func init() {
	proto.RegisterFile("protobuf/definition/service.proto", fileDescriptor_service_7c01ece46eb6a4a2)
}

var fileDescriptor_service_7c01ece46eb6a4a2 = []byte{
	// 699 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xdd, 0x6a, 0xdb, 0x4a,
	0x10, 0x46, 0xfe, 0x91, 0xed, 0x71, 0x72, 0xc8, 0x59, 0x72, 0xc2, 0xc6, 0xa7, 0x14, 0x57, 0x17,
	0xc5, 0xb9, 0x88, 0x0d, 0x0e, 0xed, 0x45, 0x29, 0x14, 0xfa, 0x13, 0x0a, 0xb9, 0x29, 0x6a, 0xa0,
	0xd0, 0x3b, 0x59, 0x1a, 0x3b, 0x5b, 0x5b, 0x5a, 0xb1, 0xbb, 0x72, 0x70, 0x2f, 0xfa, 0x2e, 0x79,
	0x81, 0x5e, 0xf5, 0x39, 0xfa, 0x1a, 0x7d, 0x8d, 0xb2, 0x23, 0xc9, 0x96, 0xc1, 0x35, 0xa1, 0xed,
	0xdd, 0xcc, 0x7c, 0xdf, 0x6a, 0x66, 0xbe, 0x99, 0x5d, 0xc1, 0xa3, 0x54, 0x49, 0x23, 0x27, 0xd9,
	0x74, 0x14, 0xe1, 0x54, 0x24, 0xc2, 0x08, 0x99, 0x8c, 0x34, 0xaa, 0xa5, 0x08, 0x71, 0x48, 0x18,
	0x83, 0x0d, 0xe2, 0xdd, 0xd5, 0xa1, 0xf5, 0x3e, 0x47, 0x19, 0x83, 0xc6, 0x4d, 0xa0, 0x6f, 0x38,
	0xf4, 0x9d, 0x41, 0xc7, 0x27, 0x9b, 0x1d, 0x41, 0x5d, 0x8b, 0x88, 0x1f, 0x50, 0xc8, 0x9a, 0x96,
	0x95, 0x04, 0x31, 0x72, 0x27, 0x67, 0x59, 0x9b, 0xf5, 0xa1, 0x1b, 0xa1, 0x0e, 0x95, 0x48, 0xed,
	0x47, 0x79, 0x8d, 0xa0, 0x6a, 0x88, 0xbd, 0x80, 0xc3, 0x50, 0x26, 0x53, 0x31, 0xcb, 0x54, 0x40,
	0x9c, 0x76, 0xdf, 0x19, 0x74, 0xc7, 0xa7, 0xc3, 0x4d, 0x2d, 0xc3, 0x57, 0x55, 0x82, 0xbf, 0xcd,
	0x67, 0x8f, 0xa1, 0x69, 0x02, 0x3d, 0xd7, 0xbc, 0xd9, 0xaf, 0x0f, 0xba, 0xe3, 0xa3, 0xea, 0xc1,
	0xeb, 0x40, 0xcf, 0xfd, 0x1c, 0x66, 0x67, 0xe0, 0xe2, 0x12, 0x13, 0xa3, 0xb9, 0x4b, 0xc4, 0x7f,
	0xab, 0xc4, 0x37, 0x16, 0xf1, 0x0b, 0x02, 0x7b, 0x06, 0x07, 0x11, 0xa6, 0x98, 0x44, 0x98, 0x84,
	0x02, 0x35, 0x6f, 0xd1, 0x81, 0x93, 0xea, 0x81, 0xd7, 0x25, 0xbe, 0xf2, 0xb7, 0xb8, 0xec, 0x21,
	0x80, 0xc2, 0x54, 0x6a, 0x61, 0xa4, 0x5a, 0xf1, 0x0e, 0x35, 0x5c, 0x89, 0xb0, 0x13, 0x70, 0xb5,
	0xcc, 0x54, 0x88, 0xfc, 0x90, 0xb0, 0xc2, 0x63, 0x63, 0xe8, 0xdc, 0x4a, 0x35, 0x9f, 0x2e, 0xe4,
	0xad, 0xe6, 0xff, 0x50, 0xc2, 0xe3, 0x6a, 0xc2, 0x0f, 0x05, 0xe8, 0x6f, 0x68, 0xde, 0x67, 0x68,
	0x52, 0xe1, 0x76, 0x18, 0x73, 0x5c, 0xf1, 0x46, 0x3e, 0x8c, 0x39, 0xae, 0x7e, 0x73, 0x18, 0x67,
	0xd0, 0x88, 0x02, 0x13, 0xf0, 0x3a, 0xe5, 0xff, 0xaf, 0x9a, 0xff, 0x5d, 0xa0, 0x82, 0x18, 0x0d,
	0x2a, 0x9f, 0x28, 0xde, 0x57, 0x07, 0x1a, 0x56, 0xde, 0x32, 0x77, 0xfb, 0x4f, 0x73, 0x9f, 0x83,
	0x2b, 0x92, 0x34, 0x5b, 0xcf, 0xe7, 0x17, 0xd9, 0x0b, 0x12, 0x1b, 0x41, 0x4b, 0x66, 0x86, 0xf8,
	0xad, 0x7d, 0xfc, 0x92, 0xe5, 0x7d, 0x77, 0xa0, 0xb3, 0x0e, 0xff, 0xb5, 0xaa, 0x19, 0x34, 0xcc,
	0x2a, 0x45, 0x5e, 0xcf, 0x4f, 0x59, 0x9b, 0xf5, 0xa0, 0x2d, 0x09, 0x0d, 0x16, 0x34, 0x92, 0xb6,
	0xbf, 0xf6, 0x2d, 0xa6, 0x30, 0xc5, 0xc0, 0x60, 0x44, 0xcb, 0xd1, 0xf6, 0xd7, 0xbe, 0x55, 0x40,
	0x4e, 0x3e, 0x61, 0x68, 0x38, 0xec, 0x55, 0x20, 0x27, 0x79, 0x77, 0x0e, 0x1c, 0x6e, 0xdd, 0x0c,
	0xc6, 0xa1, 0xb5, 0x94, 0x8b, 0x2c, 0x46, 0xcd, 0x9d, 0x7e, 0x7d, 0xd0, 0xf1, 0x4b, 0xd7, 0x36,
	0x52, 0x98, 0x97, 0x4a, 0xc6, 0xbc, 0x46, 0x68, 0x35, 0xc4, 0x8e, 0xa1, 0x99, 0x4a, 0x65, 0x34,
	0xcd, 0xbe, 0xe3, 0xe7, 0x8e, 0x6d, 0x2f, 0x50, 0x33, 0xcd, 0x1b, 0x14, 0x24, 0xdb, 0x66, 0x09,
	0x65, 0x1c, 0x07, 0x49, 0xc4, 0x9b, 0xd4, 0x75, 0xe9, 0x5a, 0x51, 0x31, 0x59, 0xd2, 0xfc, 0x3a,
	0xbe, 0x35, 0xbd, 0x6f, 0x0e, 0xc0, 0xe6, 0xaa, 0xec, 0x50, 0xfd, 0x18, 0x9a, 0x22, 0x0e, 0x66,
	0xa5, 0xec, 0xb9, 0x53, 0x6d, 0xa4, 0xb6, 0xb7, 0x91, 0xfa, 0x9e, 0x46, 0x1a, 0xbb, 0x1a, 0x71,
	0xef, 0xd3, 0x88, 0xf7, 0xa3, 0x06, 0xed, 0xf2, 0xc2, 0x95, 0x45, 0x3b, 0x9b, 0xa2, 0x9f, 0x42,
	0xcb, 0x28, 0x31, 0x9b, 0xa1, 0xa2, 0x95, 0xe8, 0x8e, 0x1f, 0xec, 0xba, 0xa9, 0xc3, 0xeb, 0x9c,
	0xe3, 0x97, 0x64, 0x36, 0x2a, 0x9f, 0xaa, 0xfc, 0x7e, 0x9d, 0xee, 0x3e, 0xb5, 0x79, 0xb3, 0x7a,
	0x5f, 0xa0, 0x55, 0x7c, 0x84, 0x79, 0x70, 0x20, 0x12, 0x6d, 0x82, 0x24, 0xc4, 0xb7, 0xf6, 0x2d,
	0xce, 0xcb, 0xd9, 0x8a, 0xd9, 0xe5, 0xa2, 0x17, 0xec, 0x0a, 0x57, 0xc5, 0xae, 0xae, 0x7d, 0x76,
	0x01, 0xee, 0x54, 0x2c, 0x0c, 0x2a, 0x5a, 0xd5, 0xee, 0xf8, 0xff, 0x9d, 0xc9, 0x2f, 0x89, 0xe2,
	0x17, 0xd4, 0x9e, 0x07, 0x6e, 0x1e, 0xb1, 0x5a, 0xd9, 0x92, 0xae, 0xd6, 0x42, 0x94, 0x6e, 0xef,
	0x79, 0xf1, 0x0e, 0xdc, 0xa7, 0xc0, 0x42, 0xca, 0xda, 0x5a, 0xca, 0x97, 0x4f, 0x3e, 0x5e, 0xcc,
	0x84, 0xb9, 0xc9, 0x26, 0xc3, 0x50, 0xc6, 0xa3, 0x18, 0xf5, 0xec, 0x7c, 0x2a, 0xb3, 0x24, 0xa2,
	0x85, 0x1e, 0x85, 0x52, 0xe1, 0x68, 0xc7, 0x7f, 0x6b, 0xe2, 0x52, 0xf0, 0xe2, 0x67, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xf5, 0x25, 0xf6, 0x5f, 0xd5, 0x06, 0x00, 0x00,
}
