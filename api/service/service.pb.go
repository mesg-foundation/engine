// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service/service.proto

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	service/service.proto

It has these top-level messages:
	EmitRequest
	EmitReply
	SubscribeRequest
	SubscribeReply
*/
package apiService

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
	Type EmitRequest_Type `protobuf:"varint,1,opt,name=type,enum=api.EmitRequest_Type" json:"type,omitempty"`
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
	Success bool   `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
	Error   string `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
}

func (m *EmitReply) Reset()                    { *m = EmitReply{} }
func (m *EmitReply) String() string            { return proto.CompactTextString(m) }
func (*EmitReply) ProtoMessage()               {}
func (*EmitReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *EmitReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

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
	Type SubscribeReply_Type `protobuf:"varint,1,opt,name=type,enum=api.SubscribeReply_Type" json:"type,omitempty"`
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
	proto.RegisterType((*EmitRequest)(nil), "api.EmitRequest")
	proto.RegisterType((*EmitReply)(nil), "api.EmitReply")
	proto.RegisterType((*SubscribeRequest)(nil), "api.SubscribeRequest")
	proto.RegisterType((*SubscribeReply)(nil), "api.SubscribeReply")
	proto.RegisterEnum("api.EmitRequest_Type", EmitRequest_Type_name, EmitRequest_Type_value)
	proto.RegisterEnum("api.SubscribeReply_Type", SubscribeReply_Type_name, SubscribeReply_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Service service

type ServiceClient interface {
	// Sends a greeting
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
	err := grpc.Invoke(ctx, "/api.Service/Emit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Service_SubscribeClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Service_serviceDesc.Streams[0], c.cc, "/api.Service/Subscribe", opts...)
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
	// Sends a greeting
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
		FullMethod: "/api.Service/Emit",
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
	ServiceName: "api.Service",
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
	Metadata: "service/service.proto",
}

func init() { proto.RegisterFile("service/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 301 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0x3d, 0x4f, 0xf3, 0x30,
	0x10, 0xc7, 0xeb, 0x34, 0x7d, 0xc9, 0x55, 0xca, 0x63, 0xdd, 0x43, 0xa5, 0x08, 0x31, 0x54, 0x9e,
	0x0a, 0x42, 0x01, 0x95, 0xb1, 0x13, 0x43, 0x26, 0xc4, 0x8b, 0x12, 0x17, 0x89, 0x09, 0xa5, 0x89,
	0x07, 0x4b, 0x45, 0x31, 0x76, 0x82, 0x94, 0x8d, 0x0f, 0xc0, 0x87, 0x46, 0x71, 0x02, 0x4a, 0xab,
	0x2e, 0x4c, 0xf6, 0xfd, 0x7d, 0xe7, 0xfb, 0xff, 0xee, 0x60, 0x6e, 0x84, 0xfe, 0x90, 0x99, 0xb8,
	0xea, 0xce, 0x50, 0xe9, 0xa2, 0x2c, 0x70, 0x98, 0x2a, 0xc9, 0xbe, 0x08, 0xcc, 0xa2, 0x37, 0x59,
	0xc6, 0xe2, 0xbd, 0x12, 0xa6, 0xc4, 0x73, 0x70, 0xcb, 0x5a, 0x89, 0x80, 0x2c, 0xc8, 0xd2, 0x5f,
	0xcd, 0xc3, 0x54, 0xc9, 0xb0, 0xf7, 0x1e, 0xf2, 0x5a, 0x89, 0xd8, 0xa6, 0xa0, 0x0f, 0x8e, 0xcc,
	0x03, 0x67, 0x41, 0x96, 0x5e, 0xec, 0xc8, 0x1c, 0x11, 0xdc, 0x3c, 0x2d, 0xd3, 0x60, 0x68, 0x15,
	0x7b, 0x67, 0x21, 0xb8, 0x4d, 0x05, 0x02, 0x8c, 0x93, 0x97, 0x84, 0x47, 0xf7, 0x74, 0x80, 0x1e,
	0x8c, 0xa2, 0xe7, 0xe8, 0x81, 0x53, 0x82, 0xff, 0x60, 0xc6, 0x6f, 0x93, 0xbb, 0xd7, 0xc7, 0x0d,
	0x7f, 0xda, 0x70, 0xea, 0xb0, 0x35, 0x78, 0x6d, 0x37, 0xb5, 0xab, 0x31, 0x80, 0x89, 0xa9, 0xb2,
	0x4c, 0x18, 0x63, 0xed, 0x4c, 0xe3, 0x9f, 0x10, 0x4f, 0x60, 0x24, 0xb4, 0x2e, 0x74, 0xd7, 0xbd,
	0x0d, 0x18, 0x02, 0x4d, 0xaa, 0xad, 0xc9, 0xb4, 0xdc, 0x8a, 0xce, 0x2f, 0xfb, 0x24, 0xe0, 0xf7,
	0xc4, 0xe6, 0xdb, 0xcb, 0x3d, 0xc4, 0xc0, 0x22, 0xee, 0xa7, 0xfc, 0x95, 0xf2, 0xec, 0x08, 0xe5,
	0x14, 0xdc, 0x06, 0x8d, 0x92, 0x95, 0x86, 0x49, 0xd2, 0x0e, 0x1e, 0x2f, 0xc0, 0x6d, 0xf0, 0x90,
	0x1e, 0xce, 0xf5, 0xd4, 0xef, 0x29, 0x6a, 0x57, 0xb3, 0x01, 0xae, 0xc1, 0xfb, 0x75, 0x85, 0xf3,
	0x43, 0x97, 0x6d, 0xd5, 0xff, 0x23, 0xe6, 0xd9, 0xe0, 0x9a, 0x6c, 0xc7, 0x76, 0xc5, 0x37, 0xdf,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xd7, 0x44, 0x8a, 0xa5, 0xfb, 0x01, 0x00, 0x00,
}
