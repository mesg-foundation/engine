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
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Service) Reset()         { *m = Service{} }
func (m *Service) String() string { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()    {}
func (*Service) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_b2f340b0800c6d58, []int{0}
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
	return fileDescriptor_service_b2f340b0800c6d58, []int{1}
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
	return fileDescriptor_service_b2f340b0800c6d58, []int{2}
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
	return fileDescriptor_service_b2f340b0800c6d58, []int{3}
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

// A configuration is the configuration of the main container of the service's instance.
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
	return fileDescriptor_service_b2f340b0800c6d58, []int{4}
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
	Env                  []string `protobuf:"bytes,9,rep,name=env,proto3" json:"env,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Dependency) Reset()         { *m = Dependency{} }
func (m *Dependency) String() string { return proto.CompactTextString(m) }
func (*Dependency) ProtoMessage()    {}
func (*Dependency) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_b2f340b0800c6d58, []int{5}
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

func (m *Dependency) GetEnv() []string {
	if m != nil {
		return m.Env
	}
	return nil
}

func init() {
	proto.RegisterType((*Service)(nil), "definition.Service")
	proto.RegisterType((*Event)(nil), "definition.Event")
	proto.RegisterType((*Task)(nil), "definition.Task")
	proto.RegisterType((*Parameter)(nil), "definition.Parameter")
	proto.RegisterType((*Configuration)(nil), "definition.Configuration")
	proto.RegisterType((*Dependency)(nil), "definition.Dependency")
}

func init() {
	proto.RegisterFile("protobuf/definition/service.proto", fileDescriptor_service_b2f340b0800c6d58)
}

var fileDescriptor_service_b2f340b0800c6d58 = []byte{
	// 546 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x96, 0xe3, 0x9f, 0xc4, 0xd3, 0x56, 0x2a, 0xab, 0x52, 0x2d, 0x1c, 0x50, 0xf0, 0x01, 0xa5,
	0x87, 0xc6, 0x52, 0x2b, 0x2e, 0x5c, 0x90, 0xf8, 0x3b, 0x23, 0xc3, 0x89, 0xdb, 0xc6, 0x9e, 0xa4,
	0x4b, 0x6a, 0xaf, 0xb5, 0xbb, 0x8e, 0x14, 0xde, 0x86, 0x17, 0xe0, 0x31, 0x38, 0xf2, 0x26, 0xbc,
	0x03, 0xda, 0x71, 0x9c, 0x38, 0x52, 0x88, 0xc4, 0xcf, 0x6d, 0x66, 0xbe, 0x6f, 0x3c, 0xf3, 0x7d,
	0xb3, 0x32, 0x3c, 0xad, 0xb5, 0xb2, 0x6a, 0xd6, 0xcc, 0xd3, 0x02, 0xe7, 0xb2, 0x92, 0x56, 0xaa,
	0x2a, 0x35, 0xa8, 0x57, 0x32, 0xc7, 0x29, 0x61, 0x0c, 0x76, 0x48, 0xf2, 0x73, 0x00, 0xc3, 0x0f,
	0x2d, 0xca, 0x18, 0x04, 0x77, 0xc2, 0xdc, 0x71, 0x18, 0x7b, 0x93, 0x38, 0xa3, 0x98, 0x9d, 0x83,
	0x6f, 0x64, 0xc1, 0x4f, 0xa9, 0xe4, 0x42, 0xc7, 0xaa, 0x44, 0x89, 0xdc, 0x6b, 0x59, 0x2e, 0x66,
	0x63, 0x38, 0x29, 0xd0, 0xe4, 0x5a, 0xd6, 0xee, 0xa3, 0x7c, 0x40, 0x50, 0xbf, 0xc4, 0x5e, 0xc2,
	0x59, 0xae, 0xaa, 0xb9, 0x5c, 0x34, 0x5a, 0x10, 0x67, 0x34, 0xf6, 0x26, 0x27, 0x37, 0x8f, 0xa6,
	0xbb, 0x5d, 0xa6, 0xaf, 0xfb, 0x84, 0x6c, 0x9f, 0xcf, 0x9e, 0x41, 0x68, 0x85, 0x59, 0x1a, 0x1e,
	0x8e, 0xfd, 0xc9, 0xc9, 0xcd, 0x79, 0xbf, 0xf1, 0xa3, 0x30, 0xcb, 0xac, 0x85, 0xd9, 0x15, 0x44,
	0xb8, 0xc2, 0xca, 0x1a, 0x1e, 0x11, 0xf1, 0x41, 0x9f, 0xf8, 0xd6, 0x21, 0xd9, 0x86, 0xc0, 0x5e,
	0xc0, 0x69, 0x81, 0x35, 0x56, 0x05, 0x56, 0xb9, 0x44, 0xc3, 0x87, 0xd4, 0x70, 0xd9, 0x6f, 0x78,
	0xd3, 0xe1, 0xeb, 0x6c, 0x8f, 0xcb, 0x9e, 0x00, 0x68, 0xac, 0x95, 0x91, 0x56, 0xe9, 0x35, 0x8f,
	0x49, 0x70, 0xaf, 0xc2, 0x2e, 0x21, 0x32, 0xaa, 0xd1, 0x39, 0xf2, 0x33, 0xc2, 0x36, 0x59, 0xf2,
	0x05, 0x42, 0x5a, 0xc2, 0x19, 0xbb, 0xc4, 0x35, 0x0f, 0x5a, 0x63, 0x97, 0xb8, 0xfe, 0x4b, 0x63,
	0xaf, 0x20, 0x28, 0x84, 0x15, 0xdc, 0xa7, 0xe5, 0x1f, 0xf6, 0x97, 0x7f, 0x2f, 0xb4, 0x28, 0xd1,
	0xa2, 0xce, 0x88, 0x92, 0x7c, 0xf3, 0x20, 0x70, 0x56, 0x75, 0xb3, 0x47, 0xff, 0x3a, 0xfb, 0x1a,
	0x22, 0x59, 0xd5, 0xcd, 0xd6, 0xeb, 0xdf, 0x4c, 0xdf, 0x90, 0x58, 0x0a, 0x43, 0xd5, 0x58, 0xe2,
	0x0f, 0x8f, 0xf1, 0x3b, 0x56, 0xf2, 0xc3, 0x83, 0x78, 0x5b, 0xfe, 0x6f, 0x5b, 0x33, 0x08, 0xec,
	0xba, 0x46, 0xee, 0xb7, 0x5d, 0x2e, 0x66, 0x8f, 0x61, 0xa4, 0x08, 0x15, 0xf7, 0x74, 0x92, 0x51,
	0xb6, 0xcd, 0x1d, 0xa6, 0xb1, 0x46, 0x61, 0xb1, 0xa0, 0x43, 0x8f, 0xb2, 0x6d, 0xee, 0x1c, 0x50,
	0xb3, 0xcf, 0x98, 0x5b, 0x0e, 0x47, 0x1d, 0x68, 0x49, 0xc9, 0x57, 0x0f, 0xce, 0xf6, 0x5e, 0x39,
	0xe3, 0x30, 0x5c, 0xa9, 0xfb, 0xa6, 0x44, 0xc3, 0xbd, 0xb1, 0x3f, 0x89, 0xb3, 0x2e, 0x75, 0x42,
	0x36, 0xe1, 0x3b, 0xad, 0x4a, 0x3e, 0x20, 0xb4, 0x5f, 0x62, 0x17, 0x10, 0xd6, 0x4a, 0x5b, 0x43,
	0xb7, 0x8f, 0xb3, 0x36, 0x71, 0xf2, 0x84, 0x5e, 0x18, 0x1e, 0x50, 0x91, 0x62, 0x37, 0x25, 0x57,
	0x65, 0x29, 0xaa, 0x82, 0x87, 0xa4, 0xba, 0x4b, 0x9d, 0xa9, 0x58, 0xad, 0xe8, 0x7e, 0x71, 0xe6,
	0xc2, 0xe4, 0xbb, 0x07, 0xb0, 0x7b, 0xf6, 0x07, 0x5c, 0xbf, 0x80, 0x50, 0x96, 0x62, 0xd1, 0xd9,
	0xde, 0x26, 0x7d, 0x21, 0x83, 0xa3, 0x42, 0xfc, 0x23, 0x42, 0x82, 0x43, 0x42, 0xa2, 0x3f, 0x11,
	0x12, 0x6f, 0x85, 0xbc, 0x7a, 0xfe, 0xe9, 0x76, 0x21, 0xed, 0x5d, 0x33, 0x9b, 0xe6, 0xaa, 0x4c,
	0x4b, 0x34, 0x8b, 0xeb, 0xb9, 0x6a, 0xaa, 0x82, 0x8c, 0x4f, 0x73, 0xa5, 0x31, 0x3d, 0xf0, 0xaf,
	0x9c, 0x45, 0x54, 0xbc, 0xfd, 0x15, 0x00, 0x00, 0xff, 0xff, 0xab, 0x02, 0xba, 0x4c, 0x49, 0x05,
	0x00, 0x00,
}
