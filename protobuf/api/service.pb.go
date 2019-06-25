// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/api/service.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import definition "github.com/mesg-foundation/core/protobuf/definition"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request's data for the `Create` API.
type CreateServiceRequest struct {
	Definition           *definition.Service `protobuf:"bytes,1,opt,name=definition,proto3" json:"definition,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *CreateServiceRequest) Reset()         { *m = CreateServiceRequest{} }
func (m *CreateServiceRequest) String() string { return proto.CompactTextString(m) }
func (*CreateServiceRequest) ProtoMessage()    {}
func (*CreateServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_2916d70627fa57e8, []int{0}
}
func (m *CreateServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateServiceRequest.Unmarshal(m, b)
}
func (m *CreateServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateServiceRequest.Marshal(b, m, deterministic)
}
func (dst *CreateServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateServiceRequest.Merge(dst, src)
}
func (m *CreateServiceRequest) XXX_Size() int {
	return xxx_messageInfo_CreateServiceRequest.Size(m)
}
func (m *CreateServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateServiceRequest proto.InternalMessageInfo

func (m *CreateServiceRequest) GetDefinition() *definition.Service {
	if m != nil {
		return m.Definition
	}
	return nil
}

// The response's data for the `Create` API.
type CreateServiceResponse struct {
	// The service's hash created.
	Hash string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	// The service's sid.
	Sid                  string   `protobuf:"bytes,2,opt,name=sid,proto3" json:"sid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateServiceResponse) Reset()         { *m = CreateServiceResponse{} }
func (m *CreateServiceResponse) String() string { return proto.CompactTextString(m) }
func (*CreateServiceResponse) ProtoMessage()    {}
func (*CreateServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_2916d70627fa57e8, []int{1}
}
func (m *CreateServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateServiceResponse.Unmarshal(m, b)
}
func (m *CreateServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateServiceResponse.Marshal(b, m, deterministic)
}
func (dst *CreateServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateServiceResponse.Merge(dst, src)
}
func (m *CreateServiceResponse) XXX_Size() int {
	return xxx_messageInfo_CreateServiceResponse.Size(m)
}
func (m *CreateServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateServiceResponse proto.InternalMessageInfo

func (m *CreateServiceResponse) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *CreateServiceResponse) GetSid() string {
	if m != nil {
		return m.Sid
	}
	return ""
}

// The request's data for the `Delete` API.
type DeleteServiceRequest struct {
	// The service's hash to delete.
	Hash                 string   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteServiceRequest) Reset()         { *m = DeleteServiceRequest{} }
func (m *DeleteServiceRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteServiceRequest) ProtoMessage()    {}
func (*DeleteServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_2916d70627fa57e8, []int{2}
}
func (m *DeleteServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteServiceRequest.Unmarshal(m, b)
}
func (m *DeleteServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteServiceRequest.Marshal(b, m, deterministic)
}
func (dst *DeleteServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteServiceRequest.Merge(dst, src)
}
func (m *DeleteServiceRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteServiceRequest.Size(m)
}
func (m *DeleteServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteServiceRequest proto.InternalMessageInfo

func (m *DeleteServiceRequest) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

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
	return fileDescriptor_service_2916d70627fa57e8, []int{3}
}
func (m *DeleteServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteServiceResponse.Unmarshal(m, b)
}
func (m *DeleteServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteServiceResponse.Marshal(b, m, deterministic)
}
func (dst *DeleteServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteServiceResponse.Merge(dst, src)
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
	Hash                 string   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetServiceRequest) Reset()         { *m = GetServiceRequest{} }
func (m *GetServiceRequest) String() string { return proto.CompactTextString(m) }
func (*GetServiceRequest) ProtoMessage()    {}
func (*GetServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_2916d70627fa57e8, []int{4}
}
func (m *GetServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetServiceRequest.Unmarshal(m, b)
}
func (m *GetServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetServiceRequest.Marshal(b, m, deterministic)
}
func (dst *GetServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetServiceRequest.Merge(dst, src)
}
func (m *GetServiceRequest) XXX_Size() int {
	return xxx_messageInfo_GetServiceRequest.Size(m)
}
func (m *GetServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetServiceRequest proto.InternalMessageInfo

func (m *GetServiceRequest) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type ListServiceRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListServiceRequest) Reset()         { *m = ListServiceRequest{} }
func (m *ListServiceRequest) String() string { return proto.CompactTextString(m) }
func (*ListServiceRequest) ProtoMessage()    {}
func (*ListServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_2916d70627fa57e8, []int{5}
}
func (m *ListServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListServiceRequest.Unmarshal(m, b)
}
func (m *ListServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListServiceRequest.Marshal(b, m, deterministic)
}
func (dst *ListServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListServiceRequest.Merge(dst, src)
}
func (m *ListServiceRequest) XXX_Size() int {
	return xxx_messageInfo_ListServiceRequest.Size(m)
}
func (m *ListServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListServiceRequest proto.InternalMessageInfo

type ListServiceResponse struct {
	// List of services that match the request's filters.
	Services             []*definition.Service `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ListServiceResponse) Reset()         { *m = ListServiceResponse{} }
func (m *ListServiceResponse) String() string { return proto.CompactTextString(m) }
func (*ListServiceResponse) ProtoMessage()    {}
func (*ListServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_2916d70627fa57e8, []int{6}
}
func (m *ListServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListServiceResponse.Unmarshal(m, b)
}
func (m *ListServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListServiceResponse.Marshal(b, m, deterministic)
}
func (dst *ListServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListServiceResponse.Merge(dst, src)
}
func (m *ListServiceResponse) XXX_Size() int {
	return xxx_messageInfo_ListServiceResponse.Size(m)
}
func (m *ListServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListServiceResponse proto.InternalMessageInfo

func (m *ListServiceResponse) GetServices() []*definition.Service {
	if m != nil {
		return m.Services
	}
	return nil
}

func init() {
	proto.RegisterType((*CreateServiceRequest)(nil), "api.CreateServiceRequest")
	proto.RegisterType((*CreateServiceResponse)(nil), "api.CreateServiceResponse")
	proto.RegisterType((*DeleteServiceRequest)(nil), "api.DeleteServiceRequest")
	proto.RegisterType((*DeleteServiceResponse)(nil), "api.DeleteServiceResponse")
	proto.RegisterType((*GetServiceRequest)(nil), "api.GetServiceRequest")
	proto.RegisterType((*ListServiceRequest)(nil), "api.ListServiceRequest")
	proto.RegisterType((*ListServiceResponse)(nil), "api.ListServiceResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServiceXClient is the client API for ServiceX service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceXClient interface {
	// Create a Service from a Service Definition.
	// It will return an unique identifier which is used to interact with the Service.
	Create(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*CreateServiceResponse, error)
	// Delete a Service.
	// An error is returned if one or more Instances of the Service are running.
	Delete(ctx context.Context, in *DeleteServiceRequest, opts ...grpc.CallOption) (*DeleteServiceResponse, error)
	// Get returns a Service matching the criteria of the request.
	Get(ctx context.Context, in *GetServiceRequest, opts ...grpc.CallOption) (*definition.Service, error)
	// List returns services specified in a request.
	List(ctx context.Context, in *ListServiceRequest, opts ...grpc.CallOption) (*ListServiceResponse, error)
}

type serviceXClient struct {
	cc *grpc.ClientConn
}

func NewServiceXClient(cc *grpc.ClientConn) ServiceXClient {
	return &serviceXClient{cc}
}

func (c *serviceXClient) Create(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*CreateServiceResponse, error) {
	out := new(CreateServiceResponse)
	err := c.cc.Invoke(ctx, "/api.ServiceX/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceXClient) Delete(ctx context.Context, in *DeleteServiceRequest, opts ...grpc.CallOption) (*DeleteServiceResponse, error) {
	out := new(DeleteServiceResponse)
	err := c.cc.Invoke(ctx, "/api.ServiceX/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceXClient) Get(ctx context.Context, in *GetServiceRequest, opts ...grpc.CallOption) (*definition.Service, error) {
	out := new(definition.Service)
	err := c.cc.Invoke(ctx, "/api.ServiceX/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceXClient) List(ctx context.Context, in *ListServiceRequest, opts ...grpc.CallOption) (*ListServiceResponse, error) {
	out := new(ListServiceResponse)
	err := c.cc.Invoke(ctx, "/api.ServiceX/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceXServer is the server API for ServiceX service.
type ServiceXServer interface {
	// Create a Service from a Service Definition.
	// It will return an unique identifier which is used to interact with the Service.
	Create(context.Context, *CreateServiceRequest) (*CreateServiceResponse, error)
	// Delete a Service.
	// An error is returned if one or more Instances of the Service are running.
	Delete(context.Context, *DeleteServiceRequest) (*DeleteServiceResponse, error)
	// Get returns a Service matching the criteria of the request.
	Get(context.Context, *GetServiceRequest) (*definition.Service, error)
	// List returns services specified in a request.
	List(context.Context, *ListServiceRequest) (*ListServiceResponse, error)
}

func RegisterServiceXServer(s *grpc.Server, srv ServiceXServer) {
	s.RegisterService(&_ServiceX_serviceDesc, srv)
}

func _ServiceX_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceXServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ServiceX/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceXServer).Create(ctx, req.(*CreateServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceX_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceXServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ServiceX/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceXServer).Delete(ctx, req.(*DeleteServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceX_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceXServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ServiceX/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceXServer).Get(ctx, req.(*GetServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceX_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceXServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ServiceX/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceXServer).List(ctx, req.(*ListServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServiceX_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.ServiceX",
	HandlerType: (*ServiceXServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _ServiceX_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ServiceX_Delete_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ServiceX_Get_Handler,
		},
		{
			MethodName: "List",
			Handler:    _ServiceX_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/api/service.proto",
}

func init() { proto.RegisterFile("protobuf/api/service.proto", fileDescriptor_service_2916d70627fa57e8) }

var fileDescriptor_service_2916d70627fa57e8 = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0x93, 0xa6, 0xaa, 0xca, 0xb1, 0xc0, 0x35, 0xa5, 0xc1, 0x53, 0xf1, 0x02, 0x62, 0x48,
	0xa4, 0x96, 0x0d, 0x31, 0x20, 0x10, 0x1d, 0x60, 0x2a, 0x0b, 0xab, 0x4b, 0xaf, 0xaa, 0x25, 0x94,
	0x98, 0xd8, 0xe5, 0x05, 0x78, 0x71, 0x14, 0xc7, 0x6d, 0x69, 0x62, 0x24, 0x36, 0xeb, 0xf7, 0x77,
	0xbf, 0x73, 0x9f, 0x02, 0x4c, 0x95, 0x85, 0x29, 0x16, 0x9b, 0x55, 0x26, 0x94, 0xcc, 0x34, 0x95,
	0x5f, 0xf2, 0x9d, 0x52, 0x1b, 0x62, 0x24, 0x94, 0x64, 0x17, 0x3b, 0x60, 0x49, 0x2b, 0x99, 0x4b,
	0x23, 0x8b, 0xfc, 0x90, 0xe3, 0xcf, 0x10, 0x3f, 0x94, 0x24, 0x0c, 0xbd, 0xd6, 0xf1, 0x9c, 0x3e,
	0x37, 0xa4, 0x0d, 0x4e, 0x01, 0xf6, 0x33, 0x49, 0x38, 0x0e, 0xaf, 0x8e, 0x27, 0x83, 0x74, 0x1f,
	0xa5, 0x5b, 0xfe, 0x17, 0xc6, 0xef, 0x60, 0xd8, 0x28, 0xd3, 0xaa, 0xc8, 0x35, 0x21, 0x42, 0x77,
	0x2d, 0xf4, 0xda, 0xf6, 0x1c, 0xcd, 0xed, 0x19, 0x4f, 0x20, 0xd2, 0x72, 0x99, 0x74, 0x6c, 0x54,
	0x1d, 0xf9, 0x35, 0xc4, 0x8f, 0xf4, 0x41, 0xad, 0x6f, 0xf1, 0x4c, 0xf3, 0x11, 0x0c, 0x1b, 0x6c,
	0xfd, 0x14, 0xbf, 0x84, 0xd3, 0x19, 0x99, 0x7f, 0x34, 0xc4, 0x80, 0x2f, 0x52, 0x37, 0x48, 0xfe,
	0x04, 0x83, 0x83, 0xd4, 0x2d, 0x90, 0x41, 0xdf, 0x79, 0xd3, 0x49, 0x38, 0x8e, 0xfe, 0x92, 0xb1,
	0x83, 0x26, 0xdf, 0x1d, 0xe8, 0xbb, 0xf4, 0x0d, 0xef, 0xa1, 0x57, 0x7b, 0xc1, 0xf3, 0x54, 0x28,
	0x99, 0xfa, 0x8c, 0x33, 0xe6, 0xbb, 0x72, 0x4b, 0x05, 0x55, 0x45, 0xbd, 0xaf, 0xab, 0xf0, 0x89,
	0x72, 0x15, 0x7e, 0x2f, 0x01, 0xde, 0x40, 0x34, 0x23, 0x83, 0x67, 0x16, 0x6a, 0x39, 0x62, 0xbe,
	0x85, 0x78, 0x80, 0xb7, 0xd0, 0xad, 0x84, 0xe0, 0xc8, 0x8e, 0xb5, 0x8d, 0xb1, 0xa4, 0x7d, 0xb1,
	0x7d, 0x72, 0xd1, 0xb3, 0x3f, 0xd9, 0xf4, 0x27, 0x00, 0x00, 0xff, 0xff, 0x20, 0x4a, 0xdb, 0x71,
	0xaa, 0x02, 0x00, 0x00,
}
