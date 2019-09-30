// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: protobuf/api/service.proto

package api

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_mesg_foundation_engine_hash "github.com/mesg-foundation/engine/hash"
	service "github.com/mesg-foundation/engine/service"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// The request's data for the `Create` API.
type CreateServiceRequest struct {
	// Service's sid.
	Sid string `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty"`
	// Service's name.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Service's description.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// Configurations related to the service
	Configuration service.Service_Configuration `protobuf:"bytes,4,opt,name=configuration,proto3" json:"configuration"`
	// The list of tasks this service can execute.
	Tasks []*service.Service_Task `protobuf:"bytes,5,rep,name=tasks,proto3" json:"tasks,omitempty"`
	// The list of events this service can emit.
	Events []*service.Service_Event `protobuf:"bytes,6,rep,name=events,proto3" json:"events,omitempty"`
	// The container dependencies this service requires.
	Dependencies []*service.Service_Dependency `protobuf:"bytes,7,rep,name=dependencies,proto3" json:"dependencies,omitempty"`
	// Service's repository url.
	Repository string `protobuf:"bytes,8,opt,name=repository,proto3" json:"repository,omitempty"`
	// The hash id of service's source code on IPFS.
	Source               string   `protobuf:"bytes,9,opt,name=source,proto3" json:"source,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateServiceRequest) Reset()         { *m = CreateServiceRequest{} }
func (m *CreateServiceRequest) String() string { return proto.CompactTextString(m) }
func (*CreateServiceRequest) ProtoMessage()    {}
func (*CreateServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{0}
}
func (m *CreateServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateServiceRequest.Unmarshal(m, b)
}
func (m *CreateServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateServiceRequest.Marshal(b, m, deterministic)
}
func (m *CreateServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateServiceRequest.Merge(m, src)
}
func (m *CreateServiceRequest) XXX_Size() int {
	return xxx_messageInfo_CreateServiceRequest.Size(m)
}
func (m *CreateServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateServiceRequest proto.InternalMessageInfo

func (m *CreateServiceRequest) GetSid() string {
	if m != nil {
		return m.Sid
	}
	return ""
}

func (m *CreateServiceRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateServiceRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateServiceRequest) GetConfiguration() service.Service_Configuration {
	if m != nil {
		return m.Configuration
	}
	return service.Service_Configuration{}
}

func (m *CreateServiceRequest) GetTasks() []*service.Service_Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

func (m *CreateServiceRequest) GetEvents() []*service.Service_Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *CreateServiceRequest) GetDependencies() []*service.Service_Dependency {
	if m != nil {
		return m.Dependencies
	}
	return nil
}

func (m *CreateServiceRequest) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func (m *CreateServiceRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

// The response's data for the `Create` API.
type CreateServiceResponse struct {
	// The service's hash created.
	Hash                 github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,1,opt,name=hash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"hash"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *CreateServiceResponse) Reset()         { *m = CreateServiceResponse{} }
func (m *CreateServiceResponse) String() string { return proto.CompactTextString(m) }
func (*CreateServiceResponse) ProtoMessage()    {}
func (*CreateServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{1}
}
func (m *CreateServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateServiceResponse.Unmarshal(m, b)
}
func (m *CreateServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateServiceResponse.Marshal(b, m, deterministic)
}
func (m *CreateServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateServiceResponse.Merge(m, src)
}
func (m *CreateServiceResponse) XXX_Size() int {
	return xxx_messageInfo_CreateServiceResponse.Size(m)
}
func (m *CreateServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateServiceResponse proto.InternalMessageInfo

// The request's data for the `Delete` API.
type DeleteServiceRequest struct {
	// The service's hash to delete.
	Hash                 github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,1,opt,name=hash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"hash"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *DeleteServiceRequest) Reset()         { *m = DeleteServiceRequest{} }
func (m *DeleteServiceRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteServiceRequest) ProtoMessage()    {}
func (*DeleteServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{2}
}
func (m *DeleteServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteServiceRequest.Unmarshal(m, b)
}
func (m *DeleteServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteServiceRequest.Marshal(b, m, deterministic)
}
func (m *DeleteServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteServiceRequest.Merge(m, src)
}
func (m *DeleteServiceRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteServiceRequest.Size(m)
}
func (m *DeleteServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteServiceRequest proto.InternalMessageInfo

// The response's data for the `Delete` API, doesn't contain anything.
type DeleteServiceResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteServiceResponse) Reset()         { *m = DeleteServiceResponse{} }
func (m *DeleteServiceResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteServiceResponse) ProtoMessage()    {}
func (*DeleteServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{3}
}
func (m *DeleteServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteServiceResponse.Unmarshal(m, b)
}
func (m *DeleteServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteServiceResponse.Marshal(b, m, deterministic)
}
func (m *DeleteServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteServiceResponse.Merge(m, src)
}
func (m *DeleteServiceResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteServiceResponse.Size(m)
}
func (m *DeleteServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteServiceResponse proto.InternalMessageInfo

// The request's data for the `Get` API.
type GetServiceRequest struct {
	// The service's hash to fetch.
	Hash                 github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,1,opt,name=hash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"hash"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *GetServiceRequest) Reset()         { *m = GetServiceRequest{} }
func (m *GetServiceRequest) String() string { return proto.CompactTextString(m) }
func (*GetServiceRequest) ProtoMessage()    {}
func (*GetServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{4}
}
func (m *GetServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetServiceRequest.Unmarshal(m, b)
}
func (m *GetServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetServiceRequest.Marshal(b, m, deterministic)
}
func (m *GetServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetServiceRequest.Merge(m, src)
}
func (m *GetServiceRequest) XXX_Size() int {
	return xxx_messageInfo_GetServiceRequest.Size(m)
}
func (m *GetServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetServiceRequest proto.InternalMessageInfo

// The request's data for the `List` API.
type ListServiceRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListServiceRequest) Reset()         { *m = ListServiceRequest{} }
func (m *ListServiceRequest) String() string { return proto.CompactTextString(m) }
func (*ListServiceRequest) ProtoMessage()    {}
func (*ListServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{5}
}
func (m *ListServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListServiceRequest.Unmarshal(m, b)
}
func (m *ListServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListServiceRequest.Marshal(b, m, deterministic)
}
func (m *ListServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListServiceRequest.Merge(m, src)
}
func (m *ListServiceRequest) XXX_Size() int {
	return xxx_messageInfo_ListServiceRequest.Size(m)
}
func (m *ListServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListServiceRequest proto.InternalMessageInfo

// The response's data for the `List` API.
type ListServiceResponse struct {
	// List of services that match the request's filters.
	Services             []*service.Service `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ListServiceResponse) Reset()         { *m = ListServiceResponse{} }
func (m *ListServiceResponse) String() string { return proto.CompactTextString(m) }
func (*ListServiceResponse) ProtoMessage()    {}
func (*ListServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{6}
}
func (m *ListServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListServiceResponse.Unmarshal(m, b)
}
func (m *ListServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListServiceResponse.Marshal(b, m, deterministic)
}
func (m *ListServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListServiceResponse.Merge(m, src)
}
func (m *ListServiceResponse) XXX_Size() int {
	return xxx_messageInfo_ListServiceResponse.Size(m)
}
func (m *ListServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListServiceResponse proto.InternalMessageInfo

func (m *ListServiceResponse) GetServices() []*service.Service {
	if m != nil {
		return m.Services
	}
	return nil
}

// The response's data for the `Exists` API.
type ExistsServiceResponse struct {
	// True if a service already exists, false otherwise.
	Exists               bool     `protobuf:"varint,1,opt,name=exists,proto3" json:"exists,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExistsServiceResponse) Reset()         { *m = ExistsServiceResponse{} }
func (m *ExistsServiceResponse) String() string { return proto.CompactTextString(m) }
func (*ExistsServiceResponse) ProtoMessage()    {}
func (*ExistsServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{7}
}
func (m *ExistsServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExistsServiceResponse.Unmarshal(m, b)
}
func (m *ExistsServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExistsServiceResponse.Marshal(b, m, deterministic)
}
func (m *ExistsServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExistsServiceResponse.Merge(m, src)
}
func (m *ExistsServiceResponse) XXX_Size() int {
	return xxx_messageInfo_ExistsServiceResponse.Size(m)
}
func (m *ExistsServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ExistsServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ExistsServiceResponse proto.InternalMessageInfo

func (m *ExistsServiceResponse) GetExists() bool {
	if m != nil {
		return m.Exists
	}
	return false
}

func init() {
	proto.RegisterType((*CreateServiceRequest)(nil), "mesg.api.CreateServiceRequest")
	proto.RegisterType((*CreateServiceResponse)(nil), "mesg.api.CreateServiceResponse")
	proto.RegisterType((*DeleteServiceRequest)(nil), "mesg.api.DeleteServiceRequest")
	proto.RegisterType((*DeleteServiceResponse)(nil), "mesg.api.DeleteServiceResponse")
	proto.RegisterType((*GetServiceRequest)(nil), "mesg.api.GetServiceRequest")
	proto.RegisterType((*ListServiceRequest)(nil), "mesg.api.ListServiceRequest")
	proto.RegisterType((*ListServiceResponse)(nil), "mesg.api.ListServiceResponse")
	proto.RegisterType((*ExistsServiceResponse)(nil), "mesg.api.ExistsServiceResponse")
}

func init() { proto.RegisterFile("protobuf/api/service.proto", fileDescriptor_0615fe53b372bcb1) }

var fileDescriptor_0615fe53b372bcb1 = []byte{
	// 539 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x4d, 0xea, 0xc4, 0x4d, 0x27, 0x45, 0x82, 0x6d, 0x52, 0x8c, 0x29, 0x4d, 0xf0, 0x29, 0x12,
	0xc2, 0x16, 0xe9, 0x89, 0x6b, 0xda, 0x10, 0x24, 0xe0, 0x62, 0x38, 0x21, 0x24, 0xd8, 0x38, 0x13,
	0x67, 0x55, 0xe2, 0x35, 0x9e, 0x75, 0x45, 0xfe, 0x81, 0x0f, 0xe3, 0x1b, 0x38, 0xf4, 0x23, 0xf8,
	0x02, 0xe4, 0xb5, 0x49, 0x9d, 0xc4, 0x3d, 0x20, 0xd1, 0xdb, 0xee, 0xcc, 0x7b, 0x6f, 0x3c, 0xef,
	0x65, 0x03, 0x76, 0x9c, 0x48, 0x25, 0xa7, 0xe9, 0xdc, 0xe3, 0xb1, 0xf0, 0x08, 0x93, 0x2b, 0x11,
	0xa0, 0xab, 0x8b, 0xac, 0xb5, 0x44, 0x0a, 0x5d, 0x1e, 0x0b, 0xdb, 0x09, 0x65, 0x28, 0xbd, 0x35,
	0x34, 0xbb, 0xe9, 0x8b, 0x3e, 0xe5, 0x68, 0xfb, 0x64, 0xdd, 0x56, 0xab, 0x18, 0x69, 0x53, 0xcb,
	0xf9, 0x61, 0x40, 0xe7, 0x3c, 0x41, 0xae, 0xf0, 0x7d, 0x5e, 0xf7, 0xf1, 0x5b, 0x8a, 0xa4, 0xd8,
	0x7d, 0x30, 0x48, 0xcc, 0xac, 0x7a, 0xbf, 0x3e, 0x38, 0xf0, 0xb3, 0x23, 0x63, 0xd0, 0x88, 0xf8,
	0x12, 0xad, 0x3d, 0x5d, 0xd2, 0x67, 0xd6, 0x87, 0xf6, 0x0c, 0x29, 0x48, 0x44, 0xac, 0x84, 0x8c,
	0x2c, 0x43, 0xb7, 0xca, 0x25, 0xf6, 0x0e, 0xee, 0x05, 0x32, 0x9a, 0x8b, 0x30, 0x4d, 0xb8, 0xc6,
	0x34, 0xfa, 0xf5, 0x41, 0x7b, 0xf8, 0xd4, 0xd5, 0x4b, 0xe8, 0x4f, 0x72, 0x8b, 0xd1, 0xee, 0x79,
	0x19, 0x38, 0x6a, 0xfc, 0xbc, 0xee, 0xd5, 0xfc, 0x4d, 0x36, 0x73, 0xa1, 0xa9, 0x38, 0x5d, 0x92,
	0xd5, 0xec, 0x1b, 0x83, 0xf6, 0xd0, 0xaa, 0x92, 0xf9, 0xc0, 0xe9, 0xd2, 0xcf, 0x61, 0xec, 0x05,
	0x98, 0x78, 0x85, 0x91, 0x22, 0xcb, 0xd4, 0x84, 0x47, 0x55, 0x84, 0x71, 0x86, 0xf0, 0x0b, 0x20,
	0x1b, 0xc1, 0xe1, 0x0c, 0x63, 0x8c, 0x66, 0x18, 0x05, 0x02, 0xc9, 0xda, 0xd7, 0xc4, 0xd3, 0x2a,
	0xe2, 0xc5, 0x5f, 0xdc, 0xca, 0xdf, 0xe0, 0xb0, 0x53, 0x80, 0x04, 0x63, 0x49, 0x42, 0xc9, 0x64,
	0x65, 0xb5, 0xb4, 0x2d, 0xa5, 0x0a, 0x3b, 0x06, 0x93, 0x64, 0x9a, 0x04, 0x68, 0x1d, 0xe8, 0x5e,
	0x71, 0x73, 0xbe, 0x40, 0x77, 0x2b, 0x0d, 0x8a, 0x65, 0x44, 0xc8, 0x26, 0xd0, 0x58, 0x70, 0x5a,
	0xe8, 0x3c, 0x0e, 0x47, 0x67, 0x99, 0x35, 0xbf, 0xae, 0x7b, 0xcf, 0x42, 0xa1, 0x16, 0xe9, 0xd4,
	0x0d, 0xe4, 0xd2, 0xcb, 0x3e, 0xef, 0xf9, 0x5c, 0xa6, 0xd1, 0x4c, 0x3b, 0xe5, 0x61, 0x14, 0x8a,
	0x08, 0xbd, 0x8c, 0xe5, 0xbe, 0xe6, 0xb4, 0xf0, 0xb5, 0x80, 0xf3, 0x19, 0x3a, 0x17, 0xf8, 0x15,
	0x77, 0xf2, 0xfe, 0x6f, 0x03, 0x1e, 0x42, 0x77, 0x6b, 0x40, 0xbe, 0x82, 0xf3, 0x09, 0x1e, 0x4c,
	0x50, 0xdd, 0xd5, 0xd8, 0x0e, 0xb0, 0xb7, 0x82, 0xb6, 0xe4, 0x9d, 0x57, 0x70, 0xb4, 0x51, 0x2d,
	0xdc, 0xf4, 0xa0, 0x55, 0x3c, 0x03, 0xb2, 0xea, 0x3a, 0xde, 0xa3, 0x8a, 0x78, 0xfd, 0x35, 0xc8,
	0xf1, 0xa0, 0x3b, 0xfe, 0x2e, 0x48, 0xd1, 0xb6, 0xd2, 0x31, 0x98, 0xa8, 0x1b, 0x7a, 0x83, 0x96,
	0x5f, 0xdc, 0x86, 0xbf, 0xf7, 0x60, 0xbf, 0xc0, 0xb2, 0x37, 0x60, 0xe6, 0xa1, 0xb2, 0xe2, 0x47,
	0xc4, 0x63, 0xe1, 0x56, 0x3d, 0x3a, 0xbb, 0x77, 0x6b, 0xbf, 0xf0, 0xb0, 0x96, 0x89, 0xe5, 0xf6,
	0x96, 0xc5, 0xaa, 0x12, 0x2d, 0x8b, 0x55, 0x07, 0x52, 0x63, 0x2f, 0xc1, 0x98, 0xa0, 0x62, 0x8f,
	0x6f, 0x90, 0x3b, 0x09, 0xd9, 0x55, 0xce, 0x38, 0x35, 0x36, 0x86, 0x46, 0xe6, 0x2c, 0x3b, 0xb9,
	0xe1, 0xee, 0xfa, 0x6f, 0x3f, 0xb9, 0xa5, 0x5b, 0x5e, 0x27, 0x37, 0xf6, 0x5f, 0xbc, 0xa9, 0x8c,
	0xc2, 0xa9, 0x8d, 0x9a, 0x1f, 0x0d, 0x1e, 0x8b, 0xa9, 0xa9, 0xff, 0xda, 0xce, 0xfe, 0x04, 0x00,
	0x00, 0xff, 0xff, 0x17, 0x2b, 0x2c, 0x45, 0x44, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceClient interface {
	// Create a Service from a Service Definition.
	// It will return an unique identifier which is used to interact with the Service.
	Create(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*CreateServiceResponse, error)
	// Delete a Service.
	// An error is returned if one or more Instances of the Service are running.
	Delete(ctx context.Context, in *DeleteServiceRequest, opts ...grpc.CallOption) (*DeleteServiceResponse, error)
	// Get returns a Service matching the criteria of the request.
	Get(ctx context.Context, in *GetServiceRequest, opts ...grpc.CallOption) (*service.Service, error)
	// List returns services specified in a request.
	List(ctx context.Context, in *ListServiceRequest, opts ...grpc.CallOption) (*ListServiceResponse, error)
	// Exists return if a service already exists.
	Exists(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*ExistsServiceResponse, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Create(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*CreateServiceResponse, error) {
	out := new(CreateServiceResponse)
	err := c.cc.Invoke(ctx, "/mesg.api.Service/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Delete(ctx context.Context, in *DeleteServiceRequest, opts ...grpc.CallOption) (*DeleteServiceResponse, error) {
	out := new(DeleteServiceResponse)
	err := c.cc.Invoke(ctx, "/mesg.api.Service/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Get(ctx context.Context, in *GetServiceRequest, opts ...grpc.CallOption) (*service.Service, error) {
	out := new(service.Service)
	err := c.cc.Invoke(ctx, "/mesg.api.Service/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) List(ctx context.Context, in *ListServiceRequest, opts ...grpc.CallOption) (*ListServiceResponse, error) {
	out := new(ListServiceResponse)
	err := c.cc.Invoke(ctx, "/mesg.api.Service/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Exists(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*ExistsServiceResponse, error) {
	out := new(ExistsServiceResponse)
	err := c.cc.Invoke(ctx, "/mesg.api.Service/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
type ServiceServer interface {
	// Create a Service from a Service Definition.
	// It will return an unique identifier which is used to interact with the Service.
	Create(context.Context, *CreateServiceRequest) (*CreateServiceResponse, error)
	// Delete a Service.
	// An error is returned if one or more Instances of the Service are running.
	Delete(context.Context, *DeleteServiceRequest) (*DeleteServiceResponse, error)
	// Get returns a Service matching the criteria of the request.
	Get(context.Context, *GetServiceRequest) (*service.Service, error)
	// List returns services specified in a request.
	List(context.Context, *ListServiceRequest) (*ListServiceResponse, error)
	// Exists return if a service already exists.
	Exists(context.Context, *CreateServiceRequest) (*ExistsServiceResponse, error)
}

// UnimplementedServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (*UnimplementedServiceServer) Create(ctx context.Context, req *CreateServiceRequest) (*CreateServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedServiceServer) Delete(ctx context.Context, req *DeleteServiceRequest) (*DeleteServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (*UnimplementedServiceServer) Get(ctx context.Context, req *GetServiceRequest) (*service.Service, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedServiceServer) List(ctx context.Context, req *ListServiceRequest) (*ListServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (*UnimplementedServiceServer) Exists(ctx context.Context, req *CreateServiceRequest) (*ExistsServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Service/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Create(ctx, req.(*CreateServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Service/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Delete(ctx, req.(*DeleteServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Service/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Get(ctx, req.(*GetServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Service/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).List(ctx, req.(*ListServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Service/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Exists(ctx, req.(*CreateServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mesg.api.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Service_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Service_Delete_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Service_Get_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Service_List_Handler,
		},
		{
			MethodName: "Exists",
			Handler:    _Service_Exists_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/api/service.proto",
}
