// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: protobuf/api/service.proto

package api

import (
	context "context"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
	Sid string `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty" validate:"printascii,max=63,domain"`
	// Service's name.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" validate:"required,printascii"`
	// Service's description.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty" validate:"printascii"`
	// Configurations related to the service
	Configuration service.Service_Configuration `protobuf:"bytes,4,opt,name=configuration,proto3" json:"configuration" validate:"required"`
	// The list of tasks this service can execute.
	Tasks []*service.Service_Task `protobuf:"bytes,5,rep,name=tasks,proto3" json:"tasks,omitempty" validate:"dive,required"`
	// The list of events this service can emit.
	Events []*service.Service_Event `protobuf:"bytes,6,rep,name=events,proto3" json:"events,omitempty" validate:"dive,required"`
	// The container dependencies this service requires.
	Dependencies []*service.Service_Dependency `protobuf:"bytes,7,rep,name=dependencies,proto3" json:"dependencies,omitempty" validate:"dive,required"`
	// Service's repository url.
	Repository string `protobuf:"bytes,8,opt,name=repository,proto3" json:"repository,omitempty" validate:"omitempty,uri"`
	// The hash id of service's source code on IPFS.
	Source               string   `protobuf:"bytes,9,opt,name=source,proto3" json:"source,omitempty" validate:"required,printascii"`
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
	Hash                 github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=hash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
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

func (m *CreateServiceResponse) GetHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Hash
	}
	return nil
}

// The request's data for the `Get` API.
type GetServiceRequest struct {
	// The service's hash to fetch.
	Hash                 github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=hash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"hash,omitempty" validate:"required,accaddress"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *GetServiceRequest) Reset()         { *m = GetServiceRequest{} }
func (m *GetServiceRequest) String() string { return proto.CompactTextString(m) }
func (*GetServiceRequest) ProtoMessage()    {}
func (*GetServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{2}
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

func (m *GetServiceRequest) GetHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Hash
	}
	return nil
}

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
	return fileDescriptor_0615fe53b372bcb1, []int{3}
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
	return fileDescriptor_0615fe53b372bcb1, []int{4}
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

// The request's data for the `List` API.
type ExistsServiceRequest struct {
	// The service's hash of the existing service. This hash is nil if exists is fals.
	Hash                 github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=hash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"hash,omitempty" validate:"required,accaddress"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *ExistsServiceRequest) Reset()         { *m = ExistsServiceRequest{} }
func (m *ExistsServiceRequest) String() string { return proto.CompactTextString(m) }
func (*ExistsServiceRequest) ProtoMessage()    {}
func (*ExistsServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{5}
}
func (m *ExistsServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExistsServiceRequest.Unmarshal(m, b)
}
func (m *ExistsServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExistsServiceRequest.Marshal(b, m, deterministic)
}
func (m *ExistsServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExistsServiceRequest.Merge(m, src)
}
func (m *ExistsServiceRequest) XXX_Size() int {
	return xxx_messageInfo_ExistsServiceRequest.Size(m)
}
func (m *ExistsServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExistsServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExistsServiceRequest proto.InternalMessageInfo

func (m *ExistsServiceRequest) GetHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Hash
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
	return fileDescriptor_0615fe53b372bcb1, []int{6}
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

// The request's data for the `Hash` API.
type HashServiceResponse struct {
	// Hash of the service calculated.
	Hash                 github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=hash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *HashServiceResponse) Reset()         { *m = HashServiceResponse{} }
func (m *HashServiceResponse) String() string { return proto.CompactTextString(m) }
func (*HashServiceResponse) ProtoMessage()    {}
func (*HashServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0615fe53b372bcb1, []int{7}
}
func (m *HashServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HashServiceResponse.Unmarshal(m, b)
}
func (m *HashServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HashServiceResponse.Marshal(b, m, deterministic)
}
func (m *HashServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HashServiceResponse.Merge(m, src)
}
func (m *HashServiceResponse) XXX_Size() int {
	return xxx_messageInfo_HashServiceResponse.Size(m)
}
func (m *HashServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HashServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HashServiceResponse proto.InternalMessageInfo

func (m *HashServiceResponse) GetHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Hash
	}
	return nil
}

func init() {
	proto.RegisterType((*CreateServiceRequest)(nil), "mesg.api.CreateServiceRequest")
	proto.RegisterType((*CreateServiceResponse)(nil), "mesg.api.CreateServiceResponse")
	proto.RegisterType((*GetServiceRequest)(nil), "mesg.api.GetServiceRequest")
	proto.RegisterType((*ListServiceRequest)(nil), "mesg.api.ListServiceRequest")
	proto.RegisterType((*ListServiceResponse)(nil), "mesg.api.ListServiceResponse")
	proto.RegisterType((*ExistsServiceRequest)(nil), "mesg.api.ExistsServiceRequest")
	proto.RegisterType((*ExistsServiceResponse)(nil), "mesg.api.ExistsServiceResponse")
	proto.RegisterType((*HashServiceResponse)(nil), "mesg.api.HashServiceResponse")
}

func init() { proto.RegisterFile("protobuf/api/service.proto", fileDescriptor_0615fe53b372bcb1) }

var fileDescriptor_0615fe53b372bcb1 = []byte{
	// 664 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x95, 0xdf, 0x4e, 0x13, 0x41,
	0x14, 0xc6, 0x5b, 0xda, 0x2e, 0x65, 0xc0, 0x0b, 0xa7, 0x60, 0x96, 0x82, 0x2c, 0x8e, 0x17, 0x72,
	0x51, 0xb6, 0x11, 0x12, 0x12, 0x51, 0x2f, 0x28, 0x22, 0x26, 0x90, 0x98, 0x2c, 0x5e, 0x19, 0x63,
	0x1c, 0x76, 0x0f, 0xed, 0x04, 0x77, 0x67, 0xd9, 0x33, 0xdb, 0xd0, 0xc4, 0xc4, 0xc7, 0xf1, 0x55,
	0x7c, 0x8a, 0x7d, 0x04, 0x2f, 0x7a, 0xe9, 0x95, 0xd9, 0xd9, 0xb5, 0xff, 0x58, 0x82, 0x31, 0x31,
	0x5e, 0xb5, 0x33, 0xf3, 0x9d, 0xdf, 0x9c, 0x33, 0xdf, 0x99, 0x59, 0xd2, 0x0c, 0x23, 0xa9, 0xe4,
	0x79, 0x7c, 0xd1, 0xe6, 0xa1, 0x68, 0x23, 0x44, 0x7d, 0xe1, 0x82, 0xad, 0x27, 0x69, 0xdd, 0x07,
	0xec, 0xda, 0x3c, 0x14, 0x4d, 0xd6, 0x95, 0x5d, 0xd9, 0x1e, 0x49, 0xd3, 0x91, 0x1e, 0xe8, 0x7f,
	0x99, 0xba, 0xb9, 0x3e, 0x5a, 0x56, 0x83, 0x10, 0x70, 0x9a, 0xc5, 0xbe, 0xd5, 0xc8, 0xf2, 0x61,
	0x04, 0x5c, 0xc1, 0x59, 0x36, 0xef, 0xc0, 0x55, 0x0c, 0xa8, 0xe8, 0x33, 0x52, 0x41, 0xe1, 0x99,
	0xe5, 0xcd, 0xf2, 0xd6, 0x42, 0xe7, 0xc9, 0x30, 0xb1, 0x1e, 0xf7, 0xf9, 0x67, 0xe1, 0x71, 0x05,
	0xfb, 0x2c, 0x8c, 0x44, 0xa0, 0x38, 0xba, 0x42, 0xb4, 0x7c, 0x7e, 0xfd, 0x72, 0x6f, 0xb7, 0xe5,
	0x49, 0x9f, 0x8b, 0x80, 0x39, 0x69, 0x0c, 0xdd, 0x23, 0xd5, 0x80, 0xfb, 0x60, 0xce, 0xe9, 0x58,
	0x36, 0x4c, 0xac, 0x8d, 0x71, 0x6c, 0x04, 0x57, 0xb1, 0x88, 0xc0, 0x6b, 0x8d, 0x21, 0xcc, 0xd1,
	0x7a, 0xfa, 0x9c, 0x2c, 0x7a, 0x80, 0x6e, 0x24, 0x42, 0x25, 0x64, 0x60, 0x56, 0x74, 0xf8, 0xea,
	0x30, 0xb1, 0x56, 0x8a, 0xb6, 0x66, 0xce, 0xa4, 0x9a, 0x7a, 0xe4, 0x9e, 0x2b, 0x83, 0x0b, 0xd1,
	0x8d, 0x23, 0xae, 0xc3, 0xab, 0x9b, 0xe5, 0xad, 0xc5, 0x9d, 0x47, 0xb6, 0x3e, 0x2c, 0x5d, 0xba,
	0x9d, 0x97, 0x68, 0x1f, 0x4e, 0x0a, 0x3b, 0x6b, 0xdf, 0x13, 0xab, 0x34, 0x4c, 0xac, 0xc6, 0xcd,
	0x24, 0x99, 0x33, 0x0d, 0xa5, 0xa7, 0xa4, 0xa6, 0x38, 0x5e, 0xa2, 0x59, 0xdb, 0xac, 0x6c, 0x2d,
	0xee, 0x98, 0x45, 0xf4, 0x77, 0x1c, 0x2f, 0x3b, 0xeb, 0xc3, 0xc4, 0x32, 0xc7, 0x40, 0x4f, 0xf4,
	0xa1, 0x35, 0xa6, 0x66, 0x10, 0xfa, 0x96, 0x18, 0xd0, 0x87, 0x40, 0xa1, 0x69, 0x68, 0xdc, 0x6a,
	0x11, 0xee, 0x28, 0x55, 0xdc, 0xc1, 0xcb, 0x31, 0xf4, 0x13, 0x59, 0xf2, 0x20, 0x84, 0xc0, 0x83,
	0xc0, 0x15, 0x80, 0xe6, 0xbc, 0xc6, 0x6e, 0x14, 0x61, 0x5f, 0xfd, 0xd6, 0x0d, 0xee, 0x60, 0x4f,
	0x11, 0xe9, 0x0b, 0x42, 0x22, 0x08, 0x25, 0x0a, 0x25, 0xa3, 0x81, 0x59, 0xd7, 0x16, 0xcd, 0xc4,
	0x4b, 0x5f, 0x28, 0xf0, 0x43, 0x35, 0x68, 0xc5, 0x91, 0x60, 0xce, 0x84, 0x9e, 0xee, 0x13, 0x03,
	0x65, 0x1c, 0xb9, 0x60, 0x2e, 0xfc, 0x71, 0x6f, 0xe4, 0x11, 0xec, 0x23, 0x59, 0x99, 0x69, 0x54,
	0x0c, 0x65, 0x80, 0x40, 0x8f, 0x48, 0xb5, 0xc7, 0xb1, 0xa7, 0x5b, 0x75, 0xa9, 0xf3, 0xf4, 0x67,
	0x62, 0x6d, 0x77, 0x85, 0xea, 0xc5, 0xe7, 0xb6, 0x2b, 0xfd, 0xb6, 0x2b, 0xd1, 0x97, 0x98, 0xff,
	0x6c, 0xa3, 0x77, 0x99, 0x5d, 0x03, 0xfb, 0xc0, 0x75, 0x0f, 0x3c, 0x2f, 0x02, 0x44, 0x47, 0x87,
	0xb3, 0x2f, 0xe4, 0xfe, 0x31, 0xa8, 0x99, 0x5b, 0xd0, 0x9d, 0x62, 0x9f, 0xdd, 0x92, 0x2e, 0x77,
	0x5d, 0x9e, 0xb1, 0xd8, 0xdf, 0xee, 0xbe, 0x4c, 0xe8, 0xa9, 0xc0, 0x99, 0xed, 0xd9, 0x6b, 0xd2,
	0x98, 0x9a, 0xcd, 0x2b, 0x6e, 0x93, 0x7a, 0x7e, 0x8b, 0xd1, 0x2c, 0x6b, 0x8b, 0x1b, 0x05, 0x16,
	0x3b, 0x23, 0x11, 0xfb, 0x4a, 0x96, 0x8f, 0xae, 0x05, 0x2a, 0xfc, 0x5f, 0xe5, 0xb5, 0xc9, 0xca,
	0x4c, 0x02, 0x79, 0x29, 0x0f, 0x88, 0x01, 0x7a, 0x41, 0xe7, 0x50, 0x77, 0xf2, 0x11, 0xfb, 0x40,
	0x1a, 0x6f, 0x38, 0xf6, 0xfe, 0x8d, 0xd7, 0x3b, 0x3f, 0xe6, 0xc8, 0x7c, 0x8e, 0xa6, 0x27, 0xc4,
	0xc8, 0xfa, 0x8a, 0xe6, 0xf7, 0x84, 0x87, 0xc2, 0x2e, 0x7a, 0x12, 0x9b, 0xd6, 0xad, 0xeb, 0x59,
	0x76, 0xac, 0x94, 0xbe, 0x9a, 0xc7, 0xa0, 0xe8, 0xda, 0x58, 0x79, 0xa3, 0xa7, 0x9a, 0x45, 0x5e,
	0xb1, 0x52, 0x5a, 0x5a, 0xea, 0x35, 0x5d, 0x1f, 0xc7, 0xde, 0xec, 0x88, 0xe6, 0xc3, 0x5b, 0x56,
	0x47, 0x19, 0x9c, 0x10, 0x23, 0x3b, 0xe9, 0xc9, 0x72, 0x8a, 0xcc, 0x9f, 0x2c, 0xa7, 0xd0, 0x1b,
	0x56, 0xa2, 0xc7, 0xa4, 0x9a, 0xba, 0x70, 0xe7, 0xc9, 0x4c, 0x64, 0x55, 0xe0, 0x1a, 0x2b, 0x75,
	0x6a, 0xef, 0x2b, 0x3c, 0x14, 0xe7, 0x86, 0xfe, 0xe8, 0xec, 0xfe, 0x0a, 0x00, 0x00, 0xff, 0xff,
	0x87, 0x3a, 0xf7, 0x67, 0xde, 0x06, 0x00, 0x00,
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
	// Get returns a Service matching the criteria of the request.
	Get(ctx context.Context, in *GetServiceRequest, opts ...grpc.CallOption) (*service.Service, error)
	// List returns services specified in a request.
	List(ctx context.Context, in *ListServiceRequest, opts ...grpc.CallOption) (*ListServiceResponse, error)
	// Exists return if a service already exists.
	Exists(ctx context.Context, in *ExistsServiceRequest, opts ...grpc.CallOption) (*ExistsServiceResponse, error)
	// Hash return the hash of a service
	Hash(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*HashServiceResponse, error)
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

func (c *serviceClient) Exists(ctx context.Context, in *ExistsServiceRequest, opts ...grpc.CallOption) (*ExistsServiceResponse, error) {
	out := new(ExistsServiceResponse)
	err := c.cc.Invoke(ctx, "/mesg.api.Service/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Hash(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*HashServiceResponse, error) {
	out := new(HashServiceResponse)
	err := c.cc.Invoke(ctx, "/mesg.api.Service/Hash", in, out, opts...)
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
	// Get returns a Service matching the criteria of the request.
	Get(context.Context, *GetServiceRequest) (*service.Service, error)
	// List returns services specified in a request.
	List(context.Context, *ListServiceRequest) (*ListServiceResponse, error)
	// Exists return if a service already exists.
	Exists(context.Context, *ExistsServiceRequest) (*ExistsServiceResponse, error)
	// Hash return the hash of a service
	Hash(context.Context, *CreateServiceRequest) (*HashServiceResponse, error)
}

// UnimplementedServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (*UnimplementedServiceServer) Create(ctx context.Context, req *CreateServiceRequest) (*CreateServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedServiceServer) Get(ctx context.Context, req *GetServiceRequest) (*service.Service, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedServiceServer) List(ctx context.Context, req *ListServiceRequest) (*ListServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (*UnimplementedServiceServer) Exists(ctx context.Context, req *ExistsServiceRequest) (*ExistsServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (*UnimplementedServiceServer) Hash(ctx context.Context, req *CreateServiceRequest) (*HashServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hash not implemented")
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
	in := new(ExistsServiceRequest)
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
		return srv.(ServiceServer).Exists(ctx, req.(*ExistsServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Hash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Hash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Service/Hash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Hash(ctx, req.(*CreateServiceRequest))
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
		{
			MethodName: "Hash",
			Handler:    _Service_Hash_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/api/service.proto",
}
