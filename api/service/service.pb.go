// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	service.proto

It has these top-level messages:
	EmitRequest
	EmitReply
	SubscribeRequest
	SubscribeReply
*/
package service

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

type EmitRequest_Type int32

const (
	EmitRequest_SYSTEM      EmitRequest_Type = 0
	EmitRequest_EVENT       EmitRequest_Type = 1
	EmitRequest_TASK_OUTPUT EmitRequest_Type = 2
)

var EmitRequest_Type_name = map[int32]string{
	0: "SYSTEM",
	1: "EVENT",
	2: "TASK_OUTPUT",
}
var EmitRequest_Type_value = map[string]int32{
	"SYSTEM":      0,
	"EVENT":       1,
	"TASK_OUTPUT": 2,
}

func (x EmitRequest_Type) String() string {
	return proto.EnumName(EmitRequest_Type_name, int32(x))
}
func (EmitRequest_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type SubscribeReply_Type int32

const (
	SubscribeReply_SYSTEM SubscribeReply_Type = 0
	SubscribeReply_TASK   SubscribeReply_Type = 1
)

var SubscribeReply_Type_name = map[int32]string{
	0: "SYSTEM",
	1: "TASK",
}
var SubscribeReply_Type_value = map[string]int32{
	"SYSTEM": 0,
	"TASK":   1,
}

func (x SubscribeReply_Type) String() string {
	return proto.EnumName(SubscribeReply_Type_name, int32(x))
}
func (SubscribeReply_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

// Emit Request
type EmitRequest struct {
	Type EmitRequest_Type `protobuf:"varint,1,opt,name=type,enum=service.EmitRequest_Type" json:"type,omitempty"`
	Id   string           `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	Data string           `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *EmitRequest) Reset()                    { *m = EmitRequest{} }
func (m *EmitRequest) String() string            { return proto.CompactTextString(m) }
func (*EmitRequest) ProtoMessage()               {}
func (*EmitRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EmitRequest) GetType() EmitRequest_Type {
	if m != nil {
		return m.Type
	}
	return EmitRequest_SYSTEM
}

func (m *EmitRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *EmitRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

// Emit Reply
type EmitReply struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
}

func (m *EmitReply) Reset()                    { *m = EmitReply{} }
func (m *EmitReply) String() string            { return proto.CompactTextString(m) }
func (*EmitReply) ProtoMessage()               {}
func (*EmitReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *EmitReply) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

// Subscribe Request
type SubscribeRequest struct {
}

func (m *SubscribeRequest) Reset()                    { *m = SubscribeRequest{} }
func (m *SubscribeRequest) String() string            { return proto.CompactTextString(m) }
func (*SubscribeRequest) ProtoMessage()               {}
func (*SubscribeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// Subscribe Reply
type SubscribeReply struct {
	Type SubscribeReply_Type `protobuf:"varint,1,opt,name=type,enum=service.SubscribeReply_Type" json:"type,omitempty"`
	Id   string              `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	Data string              `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *SubscribeReply) Reset()                    { *m = SubscribeReply{} }
func (m *SubscribeReply) String() string            { return proto.CompactTextString(m) }
func (*SubscribeReply) ProtoMessage()               {}
func (*SubscribeReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SubscribeReply) GetType() SubscribeReply_Type {
	if m != nil {
		return m.Type
	}
	return SubscribeReply_SYSTEM
}

func (m *SubscribeReply) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SubscribeReply) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*EmitRequest)(nil), "service.EmitRequest")
	proto.RegisterType((*EmitReply)(nil), "service.EmitReply")
	proto.RegisterType((*SubscribeRequest)(nil), "service.SubscribeRequest")
	proto.RegisterType((*SubscribeReply)(nil), "service.SubscribeReply")
	proto.RegisterEnum("service.EmitRequest_Type", EmitRequest_Type_name, EmitRequest_Type_value)
	proto.RegisterEnum("service.SubscribeReply_Type", SubscribeReply_Type_name, SubscribeReply_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Service service

type ServiceClient interface {
	Emit(ctx context.Context, in *EmitRequest, opts ...grpc.CallOption) (*EmitReply, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Service_SubscribeClient, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Emit(ctx context.Context, in *EmitRequest, opts ...grpc.CallOption) (*EmitReply, error) {
	out := new(EmitReply)
	err := grpc.Invoke(ctx, "/service.Service/Emit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Service_SubscribeClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Service_serviceDesc.Streams[0], c.cc, "/service.Service/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Service_SubscribeClient interface {
	Recv() (*SubscribeReply, error)
	grpc.ClientStream
}

type serviceSubscribeClient struct {
	grpc.ClientStream
}

func (x *serviceSubscribeClient) Recv() (*SubscribeReply, error) {
	m := new(SubscribeReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Service service

type ServiceServer interface {
	Emit(context.Context, *EmitRequest) (*EmitReply, error)
	Subscribe(*SubscribeRequest, Service_SubscribeServer) error
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_Emit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Emit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Service/Emit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Emit(ctx, req.(*EmitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServiceServer).Subscribe(m, &serviceSubscribeServer{stream})
}

type Service_SubscribeServer interface {
	Send(*SubscribeReply) error
	grpc.ServerStream
}

type serviceSubscribeServer struct {
	grpc.ServerStream
}

func (x *serviceSubscribeServer) Send(m *SubscribeReply) error {
	return x.ServerStream.SendMsg(m)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Emit",
			Handler:    _Service_Emit_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Service_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 284 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0xb3, 0x31, 0x6d, 0xcd, 0x14, 0x63, 0x18, 0x0a, 0xa6, 0xd2, 0x43, 0xdd, 0x53, 0x2f,
	0x86, 0x12, 0x3f, 0x81, 0x48, 0x4e, 0xe2, 0x1f, 0xb2, 0x5b, 0xc1, 0x93, 0x34, 0xcd, 0x1e, 0x16,
	0x2a, 0x59, 0x37, 0x5b, 0x21, 0x57, 0xf1, 0x03, 0xf8, 0x91, 0x25, 0x9b, 0x5a, 0x5a, 0x89, 0x07,
	0x6f, 0x3b, 0x8f, 0xf7, 0x98, 0xf7, 0xdb, 0x81, 0x93, 0x4a, 0xe8, 0x77, 0xb9, 0x12, 0xb1, 0xd2,
	0xa5, 0x29, 0x71, 0xb0, 0x1d, 0xe9, 0x17, 0x81, 0x61, 0xfa, 0x2a, 0x4d, 0x26, 0xde, 0x36, 0xa2,
	0x32, 0x78, 0x09, 0x9e, 0xa9, 0x95, 0x88, 0xc8, 0x94, 0xcc, 0x82, 0x64, 0x1c, 0xff, 0xc4, 0xf6,
	0x3c, 0x31, 0xaf, 0x95, 0xc8, 0xac, 0x0d, 0x03, 0x70, 0x65, 0x11, 0xb9, 0x53, 0x32, 0xf3, 0x33,
	0x57, 0x16, 0x88, 0xe0, 0x15, 0x4b, 0xb3, 0x8c, 0x8e, 0xac, 0x62, 0xdf, 0x34, 0x06, 0xaf, 0x49,
	0x20, 0x40, 0x9f, 0x3d, 0x33, 0x9e, 0xde, 0x85, 0x0e, 0xfa, 0xd0, 0x4b, 0x9f, 0xd2, 0x7b, 0x1e,
	0x12, 0x3c, 0x85, 0x21, 0xbf, 0x66, 0xb7, 0x2f, 0x0f, 0x0b, 0xfe, 0xb8, 0xe0, 0xa1, 0x4b, 0x2f,
	0xc0, 0x6f, 0xb7, 0xa9, 0x75, 0x8d, 0x23, 0xe8, 0x09, 0xad, 0x4b, 0x6d, 0x0b, 0xf9, 0x59, 0x3b,
	0x50, 0x84, 0x90, 0x6d, 0xf2, 0x6a, 0xa5, 0x65, 0x2e, 0xb6, 0xad, 0xe8, 0x27, 0x81, 0x60, 0x4f,
	0x6c, 0xc2, 0xf3, 0x03, 0x98, 0xc9, 0x0e, 0xe6, 0xd0, 0xf6, 0x5f, 0x9e, 0x49, 0x07, 0xcf, 0x31,
	0x78, 0x0d, 0x44, 0x48, 0x92, 0x0f, 0x02, 0x03, 0xd6, 0xee, 0xc1, 0x04, 0xbc, 0x86, 0x04, 0x47,
	0x5d, 0xdf, 0x78, 0x8e, 0xbf, 0x54, 0xb5, 0xae, 0xa9, 0x83, 0x37, 0xe0, 0xef, 0xea, 0xe1, 0xb8,
	0xab, 0x72, 0x9b, 0x3e, 0xfb, 0x83, 0x86, 0x3a, 0x73, 0x92, 0xf7, 0xed, 0x95, 0xaf, 0xbe, 0x03,
	0x00, 0x00, 0xff, 0xff, 0xa7, 0xe4, 0xe5, 0x8a, 0xf6, 0x01, 0x00, 0x00,
}
