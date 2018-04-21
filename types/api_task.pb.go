// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api_task.proto

package types

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

type ExecuteTaskRequest struct {
	Service *ProtoService `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
	Task    string        `protobuf:"bytes,2,opt,name=task" json:"task,omitempty"`
	Data    string        `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ExecuteTaskRequest) Reset()                    { *m = ExecuteTaskRequest{} }
func (m *ExecuteTaskRequest) String() string            { return proto.CompactTextString(m) }
func (*ExecuteTaskRequest) ProtoMessage()               {}
func (*ExecuteTaskRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *ExecuteTaskRequest) GetService() *ProtoService {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *ExecuteTaskRequest) GetTask() string {
	if m != nil {
		return m.Task
	}
	return ""
}

func (m *ExecuteTaskRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type ListenTaskRequest struct {
	Service *ProtoService `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *ListenTaskRequest) Reset()                    { *m = ListenTaskRequest{} }
func (m *ListenTaskRequest) String() string            { return proto.CompactTextString(m) }
func (*ListenTaskRequest) ProtoMessage()               {}
func (*ListenTaskRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *ListenTaskRequest) GetService() *ProtoService {
	if m != nil {
		return m.Service
	}
	return nil
}

type TaskReply struct {
	Task string `protobuf:"bytes,1,opt,name=task" json:"task,omitempty"`
	Data string `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *TaskReply) Reset()                    { *m = TaskReply{} }
func (m *TaskReply) String() string            { return proto.CompactTextString(m) }
func (*TaskReply) ProtoMessage()               {}
func (*TaskReply) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *TaskReply) GetTask() string {
	if m != nil {
		return m.Task
	}
	return ""
}

func (m *TaskReply) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*ExecuteTaskRequest)(nil), "types.ExecuteTaskRequest")
	proto.RegisterType((*ListenTaskRequest)(nil), "types.ListenTaskRequest")
	proto.RegisterType((*TaskReply)(nil), "types.TaskReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Task service

type TaskClient interface {
	Execute(ctx context.Context, in *ExecuteTaskRequest, opts ...grpc.CallOption) (*TaskReply, error)
	Listen(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Task_ListenClient, error)
}

type taskClient struct {
	cc *grpc.ClientConn
}

func NewTaskClient(cc *grpc.ClientConn) TaskClient {
	return &taskClient{cc}
}

func (c *taskClient) Execute(ctx context.Context, in *ExecuteTaskRequest, opts ...grpc.CallOption) (*TaskReply, error) {
	out := new(TaskReply)
	err := grpc.Invoke(ctx, "/types.Task/Execute", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskClient) Listen(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Task_ListenClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Task_serviceDesc.Streams[0], c.cc, "/types.Task/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &taskListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Task_ListenClient interface {
	Recv() (*TaskReply, error)
	grpc.ClientStream
}

type taskListenClient struct {
	grpc.ClientStream
}

func (x *taskListenClient) Recv() (*TaskReply, error) {
	m := new(TaskReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Task service

type TaskServer interface {
	Execute(context.Context, *ExecuteTaskRequest) (*TaskReply, error)
	Listen(*ListenTaskRequest, Task_ListenServer) error
}

func RegisterTaskServer(s *grpc.Server, srv TaskServer) {
	s.RegisterService(&_Task_serviceDesc, srv)
}

func _Task_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/types.Task/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServer).Execute(ctx, req.(*ExecuteTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Task_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenTaskRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TaskServer).Listen(m, &taskListenServer{stream})
}

type Task_ListenServer interface {
	Send(*TaskReply) error
	grpc.ServerStream
}

type taskListenServer struct {
	grpc.ServerStream
}

func (x *taskListenServer) Send(m *TaskReply) error {
	return x.ServerStream.SendMsg(m)
}

var _Task_serviceDesc = grpc.ServiceDesc{
	ServiceName: "types.Task",
	HandlerType: (*TaskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Execute",
			Handler:    _Task_Execute_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _Task_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api_task.proto",
}

func init() { proto.RegisterFile("api_task.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0x2c, 0xc8, 0x8c,
	0x2f, 0x49, 0x2c, 0xce, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48,
	0x2d, 0x96, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x85, 0x88, 0x2a, 0x65, 0x73, 0x09,
	0xb9, 0x56, 0xa4, 0x26, 0x97, 0x96, 0xa4, 0x86, 0x24, 0x16, 0x67, 0x07, 0xa5, 0x16, 0x96, 0xa6,
	0x16, 0x97, 0x08, 0xe9, 0x72, 0xb1, 0x43, 0x95, 0x49, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b, 0x09,
	0xeb, 0x81, 0x75, 0xeb, 0x05, 0x80, 0x34, 0x05, 0x43, 0xa4, 0x82, 0x60, 0x6a, 0x84, 0x84, 0xb8,
	0x58, 0x40, 0x16, 0x49, 0x30, 0x29, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x20, 0xb1, 0x94, 0xc4,
	0x92, 0x44, 0x09, 0x66, 0x88, 0x18, 0x88, 0xad, 0xe4, 0xc4, 0x25, 0xe8, 0x93, 0x59, 0x5c, 0x92,
	0x9a, 0x47, 0xbe, 0x5d, 0x4a, 0xc6, 0x5c, 0x9c, 0x10, 0xdd, 0x05, 0x39, 0x95, 0x70, 0x8b, 0x19,
	0xb1, 0x58, 0xcc, 0x84, 0xb0, 0xd8, 0xa8, 0x8a, 0x8b, 0x05, 0xa4, 0x49, 0xc8, 0x82, 0x8b, 0x1d,
	0xea, 0x5b, 0x21, 0x49, 0xa8, 0x2d, 0x98, 0xbe, 0x97, 0x12, 0x80, 0x4a, 0xc1, 0xed, 0x51, 0x62,
	0x10, 0xb2, 0xe0, 0x62, 0x83, 0x38, 0x5d, 0x48, 0x02, 0x2a, 0x8b, 0xe1, 0x13, 0x6c, 0xfa, 0x0c,
	0x18, 0x93, 0xd8, 0xc0, 0x01, 0x6d, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x0b, 0x0e, 0x82, 0xd6,
	0x90, 0x01, 0x00, 0x00,
}
