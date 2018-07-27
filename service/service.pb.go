// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/mesg-foundation/core/service/service.proto

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	github.com/mesg-foundation/core/service/service.proto

It has these top-level messages:
	Service
	Task
	Event
	Output
	Parameter
	Dependency
*/
package service

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

type Service struct {
	Name          string                 `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	Tasks         map[string]*Task       `protobuf:"bytes,5,rep,name=tasks" json:"tasks,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Events        map[string]*Event      `protobuf:"bytes,6,rep,name=events" json:"events,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Dependencies  map[string]*Dependency `protobuf:"bytes,7,rep,name=dependencies" json:"dependencies,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Configuration *Dependency            `protobuf:"bytes,8,opt,name=configuration" json:"configuration,omitempty"`
	Repository    string                 `protobuf:"bytes,9,opt,name=repository" json:"repository,omitempty"`
}

func (m *Service) Reset()                    { *m = Service{} }
func (m *Service) String() string            { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()               {}
func (*Service) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

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

func (m *Service) GetTasks() map[string]*Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

func (m *Service) GetEvents() map[string]*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *Service) GetDependencies() map[string]*Dependency {
	if m != nil {
		return m.Dependencies
	}
	return nil
}

func (m *Service) GetConfiguration() *Dependency {
	if m != nil {
		return m.Configuration
	}
	return nil
}

func (m *Service) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

type Task struct {
	Name        string                `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Description string                `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	Inputs      map[string]*Parameter `protobuf:"bytes,6,rep,name=inputs" json:"inputs,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Outputs     map[string]*Output    `protobuf:"bytes,7,rep,name=outputs" json:"outputs,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Task) Reset()                    { *m = Task{} }
func (m *Task) String() string            { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()               {}
func (*Task) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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

func (m *Task) GetInputs() map[string]*Parameter {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *Task) GetOutputs() map[string]*Output {
	if m != nil {
		return m.Outputs
	}
	return nil
}

type Event struct {
	Name        string                `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Description string                `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	Data        map[string]*Parameter `protobuf:"bytes,3,rep,name=data" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

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

func (m *Event) GetData() map[string]*Parameter {
	if m != nil {
		return m.Data
	}
	return nil
}

type Output struct {
	Name        string                `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Description string                `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	Data        map[string]*Parameter `protobuf:"bytes,3,rep,name=data" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Output) Reset()                    { *m = Output{} }
func (m *Output) String() string            { return proto.CompactTextString(m) }
func (*Output) ProtoMessage()               {}
func (*Output) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Output) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Output) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Output) GetData() map[string]*Parameter {
	if m != nil {
		return m.Data
	}
	return nil
}

type Parameter struct {
	Name        string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	Type        string `protobuf:"bytes,3,opt,name=type" json:"type,omitempty"`
	Optional    bool   `protobuf:"varint,4,opt,name=optional" json:"optional,omitempty"`
}

func (m *Parameter) Reset()                    { *m = Parameter{} }
func (m *Parameter) String() string            { return proto.CompactTextString(m) }
func (*Parameter) ProtoMessage()               {}
func (*Parameter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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

type Dependency struct {
	Image       string   `protobuf:"bytes,1,opt,name=image" json:"image,omitempty"`
	Volumes     []string `protobuf:"bytes,2,rep,name=volumes" json:"volumes,omitempty"`
	Volumesfrom []string `protobuf:"bytes,3,rep,name=volumesfrom" json:"volumesfrom,omitempty"`
	Ports       []string `protobuf:"bytes,4,rep,name=ports" json:"ports,omitempty"`
	Command     string   `protobuf:"bytes,5,opt,name=command" json:"command,omitempty"`
}

func (m *Dependency) Reset()                    { *m = Dependency{} }
func (m *Dependency) String() string            { return proto.CompactTextString(m) }
func (*Dependency) ProtoMessage()               {}
func (*Dependency) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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

func (m *Dependency) GetVolumesfrom() []string {
	if m != nil {
		return m.Volumesfrom
	}
	return nil
}

func (m *Dependency) GetPorts() []string {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *Dependency) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func init() {
	proto.RegisterType((*Service)(nil), "service.Service")
	proto.RegisterType((*Task)(nil), "service.Task")
	proto.RegisterType((*Event)(nil), "service.Event")
	proto.RegisterType((*Output)(nil), "service.Output")
	proto.RegisterType((*Parameter)(nil), "service.Parameter")
	proto.RegisterType((*Dependency)(nil), "service.Dependency")
}

func init() {
	proto.RegisterFile("github.com/mesg-foundation/core/service/service.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 567 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x95, 0xd1, 0x6e, 0xd3, 0x3c,
	0x14, 0xc7, 0x95, 0x36, 0x69, 0x9b, 0x93, 0xed, 0xfb, 0xc0, 0x70, 0xe1, 0x05, 0x84, 0xa2, 0x02,
	0x52, 0x90, 0x58, 0x2a, 0xca, 0x90, 0x80, 0xeb, 0x0d, 0x34, 0x4d, 0x08, 0x14, 0xf6, 0x02, 0x5e,
	0xe2, 0x96, 0x68, 0x4d, 0x1c, 0x1c, 0xa7, 0x52, 0x5f, 0x82, 0xb7, 0x41, 0x42, 0xdc, 0xf1, 0x66,
	0x28, 0x76, 0xea, 0xc6, 0x34, 0x57, 0x45, 0x5c, 0xd5, 0xc7, 0xe7, 0xfc, 0x7f, 0xe7, 0xf8, 0x6f,
	0x57, 0x81, 0x57, 0xcb, 0x4c, 0x7c, 0xa9, 0x6f, 0xa2, 0x84, 0xe5, 0xb3, 0x9c, 0x56, 0xcb, 0xd3,
	0x05, 0xab, 0x8b, 0x94, 0x88, 0x8c, 0x15, 0xb3, 0x84, 0x71, 0x3a, 0xab, 0x28, 0x5f, 0x67, 0x89,
	0xfe, 0x8d, 0x4a, 0xce, 0x04, 0x43, 0xe3, 0x36, 0x9c, 0x7e, 0xb7, 0x61, 0xfc, 0x59, 0xad, 0x11,
	0x02, 0xbb, 0x20, 0x39, 0xc5, 0x56, 0x60, 0x85, 0x6e, 0x2c, 0xd7, 0x28, 0x00, 0x2f, 0xa5, 0x55,
	0xc2, 0xb3, 0xb2, 0x41, 0xe2, 0x81, 0x4c, 0x75, 0xb7, 0xd0, 0x0b, 0x70, 0x04, 0xa9, 0x6e, 0x2b,
	0xec, 0x04, 0xc3, 0xd0, 0x9b, 0x3f, 0x88, 0xb6, 0x9d, 0x5a, 0x6c, 0x74, 0xdd, 0x64, 0x2f, 0x0a,
	0xc1, 0x37, 0xb1, 0xaa, 0x44, 0x67, 0x30, 0xa2, 0x6b, 0x5a, 0x88, 0x0a, 0x8f, 0xa4, 0xe6, 0xe1,
	0x9e, 0xe6, 0x42, 0xa6, 0x95, 0xa8, 0xad, 0x45, 0xef, 0xe0, 0x28, 0xa5, 0x25, 0x2d, 0x52, 0x5a,
	0x24, 0x19, 0xad, 0xf0, 0x58, 0x6a, 0xa7, 0x7b, 0xda, 0xf3, 0x4e, 0x91, 0x22, 0x18, 0x3a, 0xf4,
	0x06, 0x8e, 0x13, 0x56, 0x2c, 0xb2, 0x65, 0xcd, 0xa5, 0x4f, 0x78, 0x12, 0x58, 0xa1, 0x37, 0xbf,
	0xa7, 0x41, 0x1a, 0xb0, 0x89, 0xcd, 0x4a, 0xf4, 0x08, 0x80, 0xd3, 0x92, 0x55, 0x99, 0x60, 0x7c,
	0x83, 0x5d, 0x69, 0x46, 0x67, 0xc7, 0x7f, 0x0f, 0xb0, 0x3b, 0x2d, 0xba, 0x03, 0xc3, 0x5b, 0xba,
	0x69, 0xed, 0x6c, 0x96, 0xe8, 0x31, 0x38, 0x6b, 0xb2, 0xaa, 0xa9, 0xf4, 0xd1, 0x9b, 0x1f, 0xeb,
	0x96, 0x8d, 0x2a, 0x56, 0xb9, 0xb7, 0x83, 0xd7, 0x96, 0x7f, 0x09, 0x5e, 0xc7, 0x82, 0x1e, 0xd2,
	0x13, 0x93, 0xf4, 0x9f, 0x26, 0x49, 0x59, 0x17, 0x75, 0x0d, 0x77, 0xf7, 0x1c, 0xe9, 0x01, 0x3e,
	0x33, 0x81, 0xbd, 0x6e, 0xec, 0xa8, 0xd3, 0x5f, 0x03, 0xb0, 0x9b, 0xa1, 0x0f, 0x7e, 0x34, 0xa3,
	0xac, 0x28, 0x6b, 0xfd, 0x02, 0x4e, 0x0c, 0x27, 0xa2, 0x4b, 0x99, 0x6b, 0xaf, 0x5f, 0x15, 0xa2,
	0x33, 0x18, 0xb3, 0x5a, 0x48, 0x8d, 0xba, 0x79, 0xdf, 0xd4, 0x7c, 0x54, 0x49, 0x25, 0xda, 0x96,
	0xfa, 0x1f, 0xc0, 0xeb, 0xc0, 0x7a, 0xce, 0x1d, 0x9a, 0xe7, 0x46, 0x1a, 0xfa, 0x89, 0x70, 0x92,
	0x53, 0x41, 0x79, 0xd7, 0xcc, 0x2b, 0x38, 0xea, 0xf6, 0xe9, 0xe1, 0x3d, 0x35, 0x79, 0xff, 0x6b,
	0x9e, 0xd2, 0x75, 0x3d, 0xfc, 0x61, 0x81, 0x23, 0xaf, 0xeb, 0x40, 0x13, 0x9f, 0x83, 0x9d, 0x12,
	0x41, 0xf0, 0x50, 0xda, 0x81, 0xcd, 0x27, 0x10, 0x9d, 0x13, 0x41, 0x94, 0x19, 0xb2, 0xca, 0xbf,
	0x02, 0x57, 0x6f, 0xfd, 0xad, 0x0f, 0xd3, 0x9f, 0x16, 0x8c, 0xd4, 0x81, 0x0e, 0x9c, 0xfd, 0xd4,
	0x98, 0xfd, 0xe4, 0x0f, 0x97, 0xfe, 0xed, 0xf0, 0x5f, 0xc1, 0xd5, 0xfb, 0x07, 0x8e, 0x8f, 0xc0,
	0x16, 0x9b, 0x92, 0xe2, 0xa1, 0x52, 0x35, 0x6b, 0xe4, 0xc3, 0x84, 0xc9, 0x2c, 0x59, 0x61, 0x3b,
	0xb0, 0xc2, 0x49, 0xac, 0xe3, 0xe9, 0x37, 0x0b, 0x60, 0xf7, 0x47, 0x42, 0xf7, 0xc1, 0xc9, 0x72,
	0xb2, 0xdc, 0x76, 0x55, 0x01, 0xc2, 0x30, 0x5e, 0xb3, 0x55, 0x9d, 0xd3, 0x0a, 0x0f, 0x82, 0x61,
	0xe8, 0xc6, 0xdb, 0xb0, 0x19, 0xa8, 0x5d, 0x2e, 0x38, 0xcb, 0xa5, 0x69, 0x6e, 0xdc, 0xdd, 0x6a,
	0x88, 0x25, 0xe3, 0xa2, 0xc2, 0xb6, 0xcc, 0xa9, 0xa0, 0x21, 0x26, 0x2c, 0xcf, 0x49, 0x91, 0x62,
	0x47, 0x76, 0xda, 0x86, 0x37, 0x23, 0xf9, 0x1d, 0x78, 0xf9, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xf2,
	0x50, 0xc4, 0xa2, 0x40, 0x06, 0x00, 0x00,
}
