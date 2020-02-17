// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: service.proto

package service

import (
	bytes "bytes"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_mesg_foundation_engine_hash "github.com/mesg-foundation/engine/hash"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Service represents the service's type.
type Service struct {
	// Service's hash.
	Hash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,10,opt,name=hash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"hash" hash:"-" validate:"required"`
	// Service's sid.
	Sid string `protobuf:"bytes,12,opt,name=sid,proto3" json:"sid,omitempty" hash:"name:12" validate:"required,printascii,max=63,domain"`
	// Service's name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" hash:"name:1" validate:"required,printascii"`
	// Service's description.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty" hash:"name:2" validate:"printascii"`
	// Configurations related to the service
	Configuration Service_Configuration `protobuf:"bytes,8,opt,name=configuration,proto3" json:"configuration" hash:"name:8" validate:"required"`
	// The list of tasks this service can execute.
	Tasks []*Service_Task `protobuf:"bytes,5,rep,name=tasks,proto3" json:"tasks,omitempty" hash:"name:5" validate:"dive,required"`
	// The list of events this service can emit.
	Events []*Service_Event `protobuf:"bytes,6,rep,name=events,proto3" json:"events,omitempty" hash:"name:6" validate:"dive,required"`
	// The container dependencies this service requires.
	Dependencies []*Service_Dependency `protobuf:"bytes,7,rep,name=dependencies,proto3" json:"dependencies,omitempty" hash:"name:7" validate:"dive,required"`
	// Service's repository url.
	Repository string `protobuf:"bytes,9,opt,name=repository,proto3" json:"repository,omitempty" hash:"name:9" validate:"omitempty,uri"`
	// The hash id of service's source code on IPFS.
	Source               string   `protobuf:"bytes,13,opt,name=source,proto3" json:"source,omitempty" hash:"name:13" validate:"required,printascii"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Service) Reset()         { *m = Service{} }
func (m *Service) String() string { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()    {}
func (*Service) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}
func (m *Service) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service.Unmarshal(m, b)
}
func (m *Service) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service.Marshal(b, m, deterministic)
}
func (m *Service) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service.Merge(m, src)
}
func (m *Service) XXX_Size() int {
	return xxx_messageInfo_Service.Size(m)
}
func (m *Service) XXX_DiscardUnknown() {
	xxx_messageInfo_Service.DiscardUnknown(m)
}

var xxx_messageInfo_Service proto.InternalMessageInfo

// Events are emitted by the service whenever the service wants.
type Service_Event struct {
	// Event's key.
	Key string `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty" hash:"name:4" validate:"printascii"`
	// Event's name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" hash:"name:1" validate:"printascii"`
	// Event's description.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty" hash:"name:2" validate:"printascii"`
	// List of data of this event.
	Data                 []*Service_Parameter `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty" hash:"name:3" validate:"dive,required"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Service_Event) Reset()         { *m = Service_Event{} }
func (m *Service_Event) String() string { return proto.CompactTextString(m) }
func (*Service_Event) ProtoMessage()    {}
func (*Service_Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0, 0}
}
func (m *Service_Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service_Event.Unmarshal(m, b)
}
func (m *Service_Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service_Event.Marshal(b, m, deterministic)
}
func (m *Service_Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service_Event.Merge(m, src)
}
func (m *Service_Event) XXX_Size() int {
	return xxx_messageInfo_Service_Event.Size(m)
}
func (m *Service_Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Service_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Service_Event proto.InternalMessageInfo

// Task is a function that requires inputs and returns output.
type Service_Task struct {
	// Task's key.
	Key string `protobuf:"bytes,8,opt,name=key,proto3" json:"key,omitempty" hash:"name:8" validate:"printascii"`
	// Task's name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" hash:"name:1" validate:"printascii"`
	// Task's description.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty" hash:"name:2" validate:"printascii"`
	// List inputs of this task.
	Inputs []*Service_Parameter `protobuf:"bytes,6,rep,name=inputs,proto3" json:"inputs,omitempty" hash:"name:6" validate:"dive,required"`
	// List of tasks outputs.
	Outputs              []*Service_Parameter `protobuf:"bytes,7,rep,name=outputs,proto3" json:"outputs,omitempty" hash:"name:7" validate:"dive,required"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Service_Task) Reset()         { *m = Service_Task{} }
func (m *Service_Task) String() string { return proto.CompactTextString(m) }
func (*Service_Task) ProtoMessage()    {}
func (*Service_Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0, 1}
}
func (m *Service_Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service_Task.Unmarshal(m, b)
}
func (m *Service_Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service_Task.Marshal(b, m, deterministic)
}
func (m *Service_Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service_Task.Merge(m, src)
}
func (m *Service_Task) XXX_Size() int {
	return xxx_messageInfo_Service_Task.Size(m)
}
func (m *Service_Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Service_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Service_Task proto.InternalMessageInfo

// Parameter describes the task's inputs, the task's outputs, and the event's data.
type Service_Parameter struct {
	// Parameter's key.
	Key string `protobuf:"bytes,8,opt,name=key,proto3" json:"key,omitempty" hash:"name:8" validate:"printascii"`
	// Parameter's name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" hash:"name:1" validate:"printascii"`
	// Parameter's description.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty" hash:"name:2" validate:"printascii"`
	// Parameter's type: `String`, `Number`, `Boolean`, `Object` or `Any`.
	Type string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty" hash:"name:3" validate:"required,printascii,oneof=String Number Boolean Object Any"`
	// Set the parameter as optional.
	Optional bool `protobuf:"varint,4,opt,name=optional,proto3" json:"optional,omitempty" hash:"name:4"`
	// Mark a parameter as an array of the defined type.
	Repeated bool `protobuf:"varint,9,opt,name=repeated,proto3" json:"repeated,omitempty" hash:"name:9"`
	// Optional object structure type when type is set to `Object`.
	Object               []*Service_Parameter `protobuf:"bytes,10,rep,name=object,proto3" json:"object,omitempty" hash:"name:10" validate:"unique,dive,required"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Service_Parameter) Reset()         { *m = Service_Parameter{} }
func (m *Service_Parameter) String() string { return proto.CompactTextString(m) }
func (*Service_Parameter) ProtoMessage()    {}
func (*Service_Parameter) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0, 2}
}
func (m *Service_Parameter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service_Parameter.Unmarshal(m, b)
}
func (m *Service_Parameter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service_Parameter.Marshal(b, m, deterministic)
}
func (m *Service_Parameter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service_Parameter.Merge(m, src)
}
func (m *Service_Parameter) XXX_Size() int {
	return xxx_messageInfo_Service_Parameter.Size(m)
}
func (m *Service_Parameter) XXX_DiscardUnknown() {
	xxx_messageInfo_Service_Parameter.DiscardUnknown(m)
}

var xxx_messageInfo_Service_Parameter proto.InternalMessageInfo

// A configuration is the configuration of the main container of the service's instance.
type Service_Configuration struct {
	// List of volumes.
	Volumes []string `protobuf:"bytes,1,rep,name=volumes,proto3" json:"volumes,omitempty" hash:"name:1" validate:"unique,dive,printascii"`
	// List of volumes mounted from other dependencies.
	VolumesFrom []string `protobuf:"bytes,2,rep,name=volumesFrom,proto3" json:"volumesFrom,omitempty" hash:"name:2" validate:"unique,dive,printascii"`
	// List of ports the container exposes.
	Ports []string `protobuf:"bytes,3,rep,name=ports,proto3" json:"ports,omitempty" hash:"name:3" validate:"unique,dive,portmap"`
	// Args to pass to the container.
	Args []string `protobuf:"bytes,4,rep,name=args,proto3" json:"args,omitempty" hash:"name:5" validate:"dive,printascii"`
	// Command to run the container.
	Command string `protobuf:"bytes,5,opt,name=command,proto3" json:"command,omitempty" hash:"name:4" validate:"printascii"`
	// Default env vars to apply to service's instance on runtime.
	Env                  []string `protobuf:"bytes,6,rep,name=env,proto3" json:"env,omitempty" hash:"name:6" validate:"unique,dive,env"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Service_Configuration) Reset()         { *m = Service_Configuration{} }
func (m *Service_Configuration) String() string { return proto.CompactTextString(m) }
func (*Service_Configuration) ProtoMessage()    {}
func (*Service_Configuration) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0, 3}
}
func (m *Service_Configuration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service_Configuration.Unmarshal(m, b)
}
func (m *Service_Configuration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service_Configuration.Marshal(b, m, deterministic)
}
func (m *Service_Configuration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service_Configuration.Merge(m, src)
}
func (m *Service_Configuration) XXX_Size() int {
	return xxx_messageInfo_Service_Configuration.Size(m)
}
func (m *Service_Configuration) XXX_DiscardUnknown() {
	xxx_messageInfo_Service_Configuration.DiscardUnknown(m)
}

var xxx_messageInfo_Service_Configuration proto.InternalMessageInfo

// A dependency is a configuration of an other container that runs separately from the service.
type Service_Dependency struct {
	// Dependency's key.
	Key string `protobuf:"bytes,8,opt,name=key,proto3" json:"key,omitempty" hash:"name:8" validate:"printascii"`
	// Image's name of the container.
	Image string `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty" hash:"name:1" validate:"printascii"`
	// List of volumes.
	Volumes []string `protobuf:"bytes,2,rep,name=volumes,proto3" json:"volumes,omitempty" hash:"name:2" validate:"unique,dive,printascii"`
	// List of volumes mounted from other dependencies.
	VolumesFrom []string `protobuf:"bytes,3,rep,name=volumesFrom,proto3" json:"volumesFrom,omitempty" hash:"name:3" validate:"unique,dive,printascii"`
	// List of ports the container exposes.
	Ports []string `protobuf:"bytes,4,rep,name=ports,proto3" json:"ports,omitempty" hash:"name:4" validate:"unique,dive,portmap"`
	// Args to pass to the container.
	Args []string `protobuf:"bytes,6,rep,name=args,proto3" json:"args,omitempty" hash:"name:6" validate:"dive,printascii"`
	// Command to run the container.
	Command string `protobuf:"bytes,5,opt,name=command,proto3" json:"command,omitempty" hash:"name:5" validate:"printascii"`
	// Default env vars to apply to service's instance on runtime.
	Env                  []string `protobuf:"bytes,9,rep,name=env,proto3" json:"env,omitempty" hash:"name:9" validate:"unique,dive,env"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Service_Dependency) Reset()         { *m = Service_Dependency{} }
func (m *Service_Dependency) String() string { return proto.CompactTextString(m) }
func (*Service_Dependency) ProtoMessage()    {}
func (*Service_Dependency) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0, 4}
}
func (m *Service_Dependency) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Service_Dependency.Unmarshal(m, b)
}
func (m *Service_Dependency) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Service_Dependency.Marshal(b, m, deterministic)
}
func (m *Service_Dependency) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Service_Dependency.Merge(m, src)
}
func (m *Service_Dependency) XXX_Size() int {
	return xxx_messageInfo_Service_Dependency.Size(m)
}
func (m *Service_Dependency) XXX_DiscardUnknown() {
	xxx_messageInfo_Service_Dependency.DiscardUnknown(m)
}

var xxx_messageInfo_Service_Dependency proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Service)(nil), "mesg.types.Service")
	proto.RegisterType((*Service_Event)(nil), "mesg.types.Service.Event")
	proto.RegisterType((*Service_Task)(nil), "mesg.types.Service.Task")
	proto.RegisterType((*Service_Parameter)(nil), "mesg.types.Service.Parameter")
	proto.RegisterType((*Service_Configuration)(nil), "mesg.types.Service.Configuration")
	proto.RegisterType((*Service_Dependency)(nil), "mesg.types.Service.Dependency")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 942 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x97, 0xdf, 0x72, 0xdb, 0x44,
	0x14, 0xc6, 0xeb, 0x48, 0xb6, 0xe3, 0x93, 0xfa, 0x82, 0xbd, 0x5a, 0x3c, 0x10, 0xb9, 0x62, 0x06,
	0x5c, 0x70, 0xec, 0xc6, 0x6e, 0x43, 0x93, 0x4e, 0x80, 0x9a, 0xd2, 0xa1, 0xcc, 0xf0, 0x4f, 0x81,
	0x61, 0xa6, 0x33, 0xbd, 0x58, 0x4b, 0x1b, 0x65, 0x49, 0xb4, 0xab, 0xae, 0x56, 0x1e, 0x7c, 0xcf,
	0x0b, 0xf0, 0x16, 0xf0, 0x08, 0xdc, 0x72, 0xd5, 0x07, 0xe0, 0x8a, 0x0b, 0xcd, 0xc0, 0x23, 0xe8,
	0x09, 0x18, 0xad, 0xed, 0x44, 0x8e, 0x65, 0xe7, 0x1f, 0x37, 0xf4, 0xce, 0xca, 0x9c, 0xef, 0xf7,
	0x49, 0xe7, 0x7c, 0x47, 0x1b, 0x41, 0x3d, 0xa2, 0x72, 0xc4, 0x5c, 0xda, 0x09, 0xa5, 0x50, 0x02,
	0x41, 0x40, 0x23, 0xbf, 0xa3, 0xc6, 0x21, 0x8d, 0x1a, 0xb6, 0x2f, 0x7c, 0xd1, 0xd5, 0x7f, 0x1f,
	0xc6, 0x87, 0xdd, 0xec, 0x4a, 0x5f, 0xe8, 0x5f, 0x93, 0x7a, 0xfb, 0x4f, 0x0c, 0xd5, 0x83, 0x09,
	0x01, 0xf9, 0x60, 0x1e, 0x91, 0xe8, 0x08, 0x43, 0xb3, 0xd4, 0xaa, 0x0d, 0x0e, 0x5e, 0x25, 0xd6,
	0xad, 0xbf, 0x12, 0xeb, 0x03, 0x9f, 0xa9, 0xa3, 0x78, 0xd8, 0x71, 0x45, 0xd0, 0xcd, 0xe0, 0x5b,
	0x87, 0x22, 0xe6, 0x1e, 0x51, 0x4c, 0xf0, 0x2e, 0xe5, 0x3e, 0xe3, 0xb4, 0x9b, 0xa9, 0x3a, 0x9f,
	0x93, 0xe8, 0x28, 0x4d, 0xac, 0xb7, 0xb2, 0x8b, 0x3d, 0x7b, 0xcb, 0x6e, 0x8e, 0xc8, 0x09, 0xf3,
	0x88, 0xa2, 0x7b, 0xb6, 0xa4, 0x2f, 0x63, 0x26, 0xa9, 0x67, 0x3b, 0xda, 0x00, 0x7d, 0x0b, 0x46,
	0xc4, 0x3c, 0x7c, 0x5b, 0xfb, 0x7c, 0x9c, 0x26, 0xd6, 0xa3, 0x89, 0x88, 0x93, 0x80, 0xee, 0x6d,
	0xf7, 0x8a, 0xa4, 0xed, 0x50, 0x32, 0xae, 0x48, 0xe4, 0x32, 0xd6, 0x0e, 0xc8, 0x4f, 0xfb, 0x3b,
	0xfd, 0xb6, 0x27, 0x02, 0xc2, 0xb8, 0xed, 0x64, 0x2c, 0xf4, 0x04, 0xcc, 0x4c, 0x8d, 0x4b, 0x9a,
	0x79, 0x2f, 0x4d, 0xac, 0x76, 0x9e, 0x79, 0x01, 0xd2, 0x76, 0xb4, 0x1a, 0x3d, 0x83, 0x0d, 0x8f,
	0x46, 0xae, 0x64, 0x61, 0xf6, 0x78, 0x78, 0x4d, 0xc3, 0xde, 0x4b, 0x13, 0xeb, 0x9d, 0x1c, 0x6c,
	0xee, 0xfe, 0xf2, 0x8c, 0xbc, 0x16, 0x49, 0xa8, 0xbb, 0x82, 0x1f, 0x32, 0x3f, 0x96, 0xba, 0x57,
	0x78, 0xbd, 0x59, 0x6a, 0x6d, 0xf4, 0xee, 0x74, 0xce, 0x06, 0xd4, 0x99, 0x36, 0xbe, 0xf3, 0x69,
	0xbe, 0x70, 0x70, 0x37, 0x6b, 0x7c, 0x9a, 0x58, 0x77, 0x72, 0x9e, 0x0f, 0x8b, 0xdb, 0x39, 0x6f,
	0x81, 0x9e, 0x43, 0x59, 0x91, 0xe8, 0x38, 0xc2, 0xe5, 0xa6, 0xd1, 0xda, 0xe8, 0xe1, 0x22, 0xaf,
	0xef, 0x48, 0x74, 0x3c, 0x78, 0x3f, 0x4d, 0xac, 0x77, 0x73, 0xf8, 0x07, 0x79, 0xbc, 0xc7, 0x46,
	0xb4, 0x7d, 0xe6, 0x31, 0x41, 0xa2, 0x17, 0x50, 0xa1, 0x23, 0xca, 0x55, 0x84, 0x2b, 0x1a, 0xfe,
	0x66, 0x11, 0xfc, 0xb3, 0xac, 0x62, 0x81, 0xbe, 0xb3, 0x82, 0x3e, 0x85, 0x22, 0x0e, 0xb7, 0x3d,
	0x1a, 0x52, 0xee, 0x51, 0xee, 0x32, 0x1a, 0xe1, 0xaa, 0x36, 0xd9, 0x2c, 0x32, 0x79, 0x32, 0xab,
	0x1b, 0x2f, 0x38, 0x7d, 0xb8, 0xc2, 0x69, 0x8e, 0x8f, 0xbe, 0x00, 0x90, 0x34, 0x14, 0x11, 0x53,
	0x42, 0x8e, 0x71, 0x4d, 0x0f, 0xfa, 0x3c, 0x6d, 0x37, 0x4f, 0x13, 0x01, 0x53, 0x34, 0x08, 0xd5,
	0xb8, 0x1d, 0x4b, 0x66, 0x3b, 0x39, 0x35, 0x7a, 0x06, 0x95, 0x48, 0xc4, 0xd2, 0xa5, 0xb8, 0xae,
	0x39, 0xdb, 0x69, 0x62, 0x6d, 0xe5, 0xd3, 0xd7, 0xbf, 0x30, 0x7e, 0x53, 0x40, 0xe3, 0xb7, 0x35,
	0x28, 0xeb, 0x26, 0xa2, 0x5d, 0x30, 0x8e, 0xe9, 0x18, 0x9b, 0x85, 0x11, 0xbc, 0xbf, 0x2c, 0x82,
	0x99, 0x06, 0x3d, 0x9a, 0xdb, 0x85, 0xf3, 0xda, 0xed, 0x65, 0xda, 0xff, 0x7c, 0x05, 0x5e, 0x80,
	0xe9, 0x11, 0x45, 0xb0, 0xa1, 0x67, 0xf9, 0x76, 0xd1, 0x2c, 0xbf, 0x21, 0x92, 0x04, 0x54, 0x51,
	0xb9, 0xd0, 0xfc, 0xfe, 0x8a, 0x51, 0x6a, 0x6c, 0xe3, 0x17, 0x03, 0xcc, 0x2c, 0xcd, 0xb3, 0x56,
	0xad, 0x17, 0xde, 0xea, 0xc3, 0xff, 0x45, 0xab, 0x08, 0x54, 0x18, 0x0f, 0xe3, 0xd3, 0xed, 0xba,
	0x62, 0xb3, 0x56, 0x6e, 0xd8, 0x04, 0x8c, 0x5c, 0xa8, 0x8a, 0x58, 0x69, 0x8f, 0xea, 0x75, 0x3c,
	0x56, 0xed, 0xd6, 0x8c, 0xdc, 0xf8, 0xd9, 0x84, 0xda, 0x29, 0xe2, 0x75, 0x18, 0xcc, 0x31, 0x98,
	0x59, 0x83, 0xb0, 0xa1, 0x19, 0x3f, 0xa4, 0x89, 0x75, 0xb0, 0x2c, 0xa4, 0x45, 0x47, 0x95, 0xe0,
	0x54, 0x1c, 0xee, 0x1f, 0x28, 0xc9, 0xb8, 0xdf, 0xfc, 0x2a, 0x0e, 0x86, 0x54, 0x36, 0x07, 0x42,
	0x9c, 0x50, 0xc2, 0x9b, 0x5f, 0x0f, 0x7f, 0xa4, 0xae, 0x6a, 0x3e, 0xe6, 0x63, 0xdb, 0xd1, 0x26,
	0x68, 0x0b, 0xd6, 0x85, 0xb6, 0x25, 0x27, 0x7a, 0xf1, 0xd7, 0x07, 0x6f, 0xa4, 0x89, 0x55, 0x9f,
	0x5b, 0x7c, 0xe7, 0xb4, 0x24, 0x2b, 0x97, 0x34, 0xa4, 0x44, 0x51, 0x4f, 0xbf, 0xc1, 0x16, 0xcb,
	0x77, 0x6d, 0xe7, 0xb4, 0x04, 0x31, 0xa8, 0x08, 0x6d, 0x89, 0xe1, 0x32, 0xf3, 0xef, 0xa5, 0x89,
	0xd5, 0xc9, 0xf7, 0xfc, 0x5e, 0xfe, 0x61, 0x63, 0xce, 0x5e, 0xc6, 0xb4, 0x7d, 0x3e, 0x6b, 0x13,
	0x83, 0xc6, 0x1f, 0x06, 0xd4, 0xe7, 0x0e, 0x35, 0xf4, 0x25, 0x54, 0x47, 0xe2, 0x24, 0x0e, 0x68,
	0x84, 0x4b, 0x4d, 0xa3, 0x55, 0x1b, 0xf4, 0xd3, 0xc4, 0xea, 0x2e, 0x1b, 0x69, 0x9e, 0x9e, 0x1f,
	0xcd, 0x8c, 0x81, 0xbe, 0x87, 0x8d, 0xe9, 0xcf, 0xa7, 0x52, 0x04, 0x78, 0xad, 0x10, 0xd9, 0xbb,
	0x0c, 0x32, 0xcf, 0x41, 0x4f, 0xa1, 0x1c, 0x0a, 0xa9, 0x22, 0xfd, 0xca, 0x5a, 0xfc, 0x37, 0xa2,
	0xbf, 0x14, 0x28, 0xa4, 0x0a, 0x48, 0x68, 0x3b, 0x13, 0x39, 0xfa, 0x04, 0x4c, 0x22, 0xfd, 0x08,
	0x9b, 0x1a, 0xd3, 0x4e, 0x13, 0xab, 0xb5, 0xf2, 0xb4, 0x9d, 0x8b, 0x70, 0xa6, 0x44, 0x8f, 0xa1,
	0xea, 0x8a, 0x20, 0x20, 0xdc, 0xc3, 0xe5, 0xab, 0x1d, 0x01, 0x33, 0x1d, 0xfa, 0x08, 0x0c, 0xca,
	0x47, 0xfa, 0x85, 0xb2, 0x78, 0x0f, 0x3b, 0xcb, 0x1e, 0x85, 0xf2, 0x91, 0xed, 0x64, 0xc2, 0xc6,
	0xef, 0x26, 0xc0, 0xd9, 0x59, 0x7b, 0x93, 0x65, 0xde, 0x87, 0x32, 0x0b, 0x88, 0x7f, 0xe5, 0x6d,
	0x9e, 0xa8, 0xf2, 0xd9, 0xb9, 0xc1, 0xa0, 0x97, 0x65, 0xc7, 0x28, 0x44, 0xf6, 0xaf, 0x9f, 0x1d,
	0xb3, 0x30, 0x3b, 0xf7, 0xaf, 0x9a, 0x9d, 0x4b, 0xcc, 0xed, 0xba, 0xd9, 0x79, 0x70, 0xd9, 0xec,
	0xd4, 0x0a, 0xef, 0x61, 0xf7, 0xc2, 0xec, 0x0c, 0xfa, 0xaf, 0xfe, 0xde, 0xbc, 0xf5, 0xeb, 0x3f,
	0x9b, 0xa5, 0xe7, 0x77, 0x2f, 0xfe, 0x7c, 0x98, 0x7e, 0xc1, 0x0c, 0x2b, 0xfa, 0x93, 0xa4, 0xff,
	0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7a, 0xb0, 0xb4, 0xf2, 0xd3, 0x0c, 0x00, 0x00,
}

func (this *Service) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Service)
	if !ok {
		that2, ok := that.(Service)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Hash.Equal(that1.Hash) {
		return false
	}
	if this.Sid != that1.Sid {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if !this.Configuration.Equal(&that1.Configuration) {
		return false
	}
	if len(this.Tasks) != len(that1.Tasks) {
		return false
	}
	for i := range this.Tasks {
		if !this.Tasks[i].Equal(that1.Tasks[i]) {
			return false
		}
	}
	if len(this.Events) != len(that1.Events) {
		return false
	}
	for i := range this.Events {
		if !this.Events[i].Equal(that1.Events[i]) {
			return false
		}
	}
	if len(this.Dependencies) != len(that1.Dependencies) {
		return false
	}
	for i := range this.Dependencies {
		if !this.Dependencies[i].Equal(that1.Dependencies[i]) {
			return false
		}
	}
	if this.Repository != that1.Repository {
		return false
	}
	if this.Source != that1.Source {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *Service_Event) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Service_Event)
	if !ok {
		that2, ok := that.(Service_Event)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Key != that1.Key {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if len(this.Data) != len(that1.Data) {
		return false
	}
	for i := range this.Data {
		if !this.Data[i].Equal(that1.Data[i]) {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *Service_Task) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Service_Task)
	if !ok {
		that2, ok := that.(Service_Task)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Key != that1.Key {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if len(this.Inputs) != len(that1.Inputs) {
		return false
	}
	for i := range this.Inputs {
		if !this.Inputs[i].Equal(that1.Inputs[i]) {
			return false
		}
	}
	if len(this.Outputs) != len(that1.Outputs) {
		return false
	}
	for i := range this.Outputs {
		if !this.Outputs[i].Equal(that1.Outputs[i]) {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *Service_Parameter) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Service_Parameter)
	if !ok {
		that2, ok := that.(Service_Parameter)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Key != that1.Key {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if this.Type != that1.Type {
		return false
	}
	if this.Optional != that1.Optional {
		return false
	}
	if this.Repeated != that1.Repeated {
		return false
	}
	if len(this.Object) != len(that1.Object) {
		return false
	}
	for i := range this.Object {
		if !this.Object[i].Equal(that1.Object[i]) {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *Service_Configuration) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Service_Configuration)
	if !ok {
		that2, ok := that.(Service_Configuration)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Volumes) != len(that1.Volumes) {
		return false
	}
	for i := range this.Volumes {
		if this.Volumes[i] != that1.Volumes[i] {
			return false
		}
	}
	if len(this.VolumesFrom) != len(that1.VolumesFrom) {
		return false
	}
	for i := range this.VolumesFrom {
		if this.VolumesFrom[i] != that1.VolumesFrom[i] {
			return false
		}
	}
	if len(this.Ports) != len(that1.Ports) {
		return false
	}
	for i := range this.Ports {
		if this.Ports[i] != that1.Ports[i] {
			return false
		}
	}
	if len(this.Args) != len(that1.Args) {
		return false
	}
	for i := range this.Args {
		if this.Args[i] != that1.Args[i] {
			return false
		}
	}
	if this.Command != that1.Command {
		return false
	}
	if len(this.Env) != len(that1.Env) {
		return false
	}
	for i := range this.Env {
		if this.Env[i] != that1.Env[i] {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *Service_Dependency) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Service_Dependency)
	if !ok {
		that2, ok := that.(Service_Dependency)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Key != that1.Key {
		return false
	}
	if this.Image != that1.Image {
		return false
	}
	if len(this.Volumes) != len(that1.Volumes) {
		return false
	}
	for i := range this.Volumes {
		if this.Volumes[i] != that1.Volumes[i] {
			return false
		}
	}
	if len(this.VolumesFrom) != len(that1.VolumesFrom) {
		return false
	}
	for i := range this.VolumesFrom {
		if this.VolumesFrom[i] != that1.VolumesFrom[i] {
			return false
		}
	}
	if len(this.Ports) != len(that1.Ports) {
		return false
	}
	for i := range this.Ports {
		if this.Ports[i] != that1.Ports[i] {
			return false
		}
	}
	if len(this.Args) != len(that1.Args) {
		return false
	}
	for i := range this.Args {
		if this.Args[i] != that1.Args[i] {
			return false
		}
	}
	if this.Command != that1.Command {
		return false
	}
	if len(this.Env) != len(that1.Env) {
		return false
	}
	for i := range this.Env {
		if this.Env[i] != that1.Env[i] {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
