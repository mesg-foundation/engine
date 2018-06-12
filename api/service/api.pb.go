// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/mesg-foundation/core/api/service/api.proto

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	github.com/mesg-foundation/core/api/service/api.proto

It has these top-level messages:
	EmitEventRequest
	ListenTaskRequest
	SubmitResultRequest
	EmitEventReply
	TaskData
	SubmitResultReply
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

type EmitEventRequest struct {
	ServiceHash string `protobuf:"bytes,1,opt,name=serviceHash" json:"serviceHash,omitempty"`
	EventKey    string `protobuf:"bytes,2,opt,name=eventKey" json:"eventKey,omitempty"`
	EventData   string `protobuf:"bytes,3,opt,name=eventData" json:"eventData,omitempty"`
}

func (m *EmitEventRequest) Reset()                    { *m = EmitEventRequest{} }
func (m *EmitEventRequest) String() string            { return proto.CompactTextString(m) }
func (*EmitEventRequest) ProtoMessage()               {}
func (*EmitEventRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EmitEventRequest) GetServiceHash() string {
	if m != nil {
		return m.ServiceHash
	}
	return ""
}

func (m *EmitEventRequest) GetEventKey() string {
	if m != nil {
		return m.EventKey
	}
	return ""
}

func (m *EmitEventRequest) GetEventData() string {
	if m != nil {
		return m.EventData
	}
	return ""
}

type ListenTaskRequest struct {
	ServiceHash string `protobuf:"bytes,1,opt,name=serviceHash" json:"serviceHash,omitempty"`
}

func (m *ListenTaskRequest) Reset()                    { *m = ListenTaskRequest{} }
func (m *ListenTaskRequest) String() string            { return proto.CompactTextString(m) }
func (*ListenTaskRequest) ProtoMessage()               {}
func (*ListenTaskRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ListenTaskRequest) GetServiceHash() string {
	if m != nil {
		return m.ServiceHash
	}
	return ""
}

type SubmitResultRequest struct {
	ExecutionID string `protobuf:"bytes,1,opt,name=executionID" json:"executionID,omitempty"`
	OutputKey   string `protobuf:"bytes,2,opt,name=outputKey" json:"outputKey,omitempty"`
	OutputData  string `protobuf:"bytes,3,opt,name=outputData" json:"outputData,omitempty"`
}

func (m *SubmitResultRequest) Reset()                    { *m = SubmitResultRequest{} }
func (m *SubmitResultRequest) String() string            { return proto.CompactTextString(m) }
func (*SubmitResultRequest) ProtoMessage()               {}
func (*SubmitResultRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SubmitResultRequest) GetExecutionID() string {
	if m != nil {
		return m.ExecutionID
	}
	return ""
}

func (m *SubmitResultRequest) GetOutputKey() string {
	if m != nil {
		return m.OutputKey
	}
	return ""
}

func (m *SubmitResultRequest) GetOutputData() string {
	if m != nil {
		return m.OutputData
	}
	return ""
}

type EmitEventReply struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
}

func (m *EmitEventReply) Reset()                    { *m = EmitEventReply{} }
func (m *EmitEventReply) String() string            { return proto.CompactTextString(m) }
func (*EmitEventReply) ProtoMessage()               {}
func (*EmitEventReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *EmitEventReply) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type TaskData struct {
	ExecutionID string `protobuf:"bytes,1,opt,name=executionID" json:"executionID,omitempty"`
	Error       string `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
	TaskKey     string `protobuf:"bytes,3,opt,name=taskKey" json:"taskKey,omitempty"`
	InputData   string `protobuf:"bytes,4,opt,name=inputData" json:"inputData,omitempty"`
}

func (m *TaskData) Reset()                    { *m = TaskData{} }
func (m *TaskData) String() string            { return proto.CompactTextString(m) }
func (*TaskData) ProtoMessage()               {}
func (*TaskData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *TaskData) GetExecutionID() string {
	if m != nil {
		return m.ExecutionID
	}
	return ""
}

func (m *TaskData) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *TaskData) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

func (m *TaskData) GetInputData() string {
	if m != nil {
		return m.InputData
	}
	return ""
}

type SubmitResultReply struct {
	Error       string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	ExecutionID string `protobuf:"bytes,2,opt,name=executionID" json:"executionID,omitempty"`
}

func (m *SubmitResultReply) Reset()                    { *m = SubmitResultReply{} }
func (m *SubmitResultReply) String() string            { return proto.CompactTextString(m) }
func (*SubmitResultReply) ProtoMessage()               {}
func (*SubmitResultReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SubmitResultReply) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *SubmitResultReply) GetExecutionID() string {
	if m != nil {
		return m.ExecutionID
	}
	return ""
}

func init() {
	proto.RegisterType((*EmitEventRequest)(nil), "api.EmitEventRequest")
	proto.RegisterType((*ListenTaskRequest)(nil), "api.ListenTaskRequest")
	proto.RegisterType((*SubmitResultRequest)(nil), "api.SubmitResultRequest")
	proto.RegisterType((*EmitEventReply)(nil), "api.EmitEventReply")
	proto.RegisterType((*TaskData)(nil), "api.TaskData")
	proto.RegisterType((*SubmitResultReply)(nil), "api.SubmitResultReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Service service

type ServiceClient interface {
	EmitEvent(ctx context.Context, in *EmitEventRequest, opts ...grpc.CallOption) (*EmitEventReply, error)
	ListenTask(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Service_ListenTaskClient, error)
	SubmitResult(ctx context.Context, in *SubmitResultRequest, opts ...grpc.CallOption) (*SubmitResultReply, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) EmitEvent(ctx context.Context, in *EmitEventRequest, opts ...grpc.CallOption) (*EmitEventReply, error) {
	out := new(EmitEventReply)
	err := grpc.Invoke(ctx, "/api.Service/EmitEvent", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ListenTask(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Service_ListenTaskClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Service_serviceDesc.Streams[0], c.cc, "/api.Service/ListenTask", opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceListenTaskClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Service_ListenTaskClient interface {
	Recv() (*TaskData, error)
	grpc.ClientStream
}

type serviceListenTaskClient struct {
	grpc.ClientStream
}

func (x *serviceListenTaskClient) Recv() (*TaskData, error) {
	m := new(TaskData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *serviceClient) SubmitResult(ctx context.Context, in *SubmitResultRequest, opts ...grpc.CallOption) (*SubmitResultReply, error) {
	out := new(SubmitResultReply)
	err := grpc.Invoke(ctx, "/api.Service/SubmitResult", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Service service

type ServiceServer interface {
	EmitEvent(context.Context, *EmitEventRequest) (*EmitEventReply, error)
	ListenTask(*ListenTaskRequest, Service_ListenTaskServer) error
	SubmitResult(context.Context, *SubmitResultRequest) (*SubmitResultReply, error)
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_EmitEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmitEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).EmitEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Service/EmitEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).EmitEvent(ctx, req.(*EmitEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ListenTask_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenTaskRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServiceServer).ListenTask(m, &serviceListenTaskServer{stream})
}

type Service_ListenTaskServer interface {
	Send(*TaskData) error
	grpc.ServerStream
}

type serviceListenTaskServer struct {
	grpc.ServerStream
}

func (x *serviceListenTaskServer) Send(m *TaskData) error {
	return x.ServerStream.SendMsg(m)
}

func _Service_SubmitResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).SubmitResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Service/SubmitResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).SubmitResult(ctx, req.(*SubmitResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EmitEvent",
			Handler:    _Service_EmitEvent_Handler,
		},
		{
			MethodName: "SubmitResult",
			Handler:    _Service_SubmitResult_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListenTask",
			Handler:       _Service_ListenTask_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/mesg-foundation/core/api/service/api.proto",
}

func init() {
	proto.RegisterFile("github.com/mesg-foundation/core/api/service/api.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 368 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0x4f, 0x2a, 0x31,
	0x10, 0xc7, 0x59, 0x78, 0xef, 0xc1, 0xce, 0x53, 0x23, 0x45, 0xc9, 0x66, 0x63, 0x0c, 0xe9, 0xc1,
	0x78, 0x91, 0x35, 0x1a, 0x62, 0xbc, 0x12, 0x48, 0x34, 0x78, 0x02, 0x4f, 0xde, 0xca, 0x5a, 0xa1,
	0x81, 0xdd, 0xae, 0xdb, 0x96, 0x88, 0xdf, 0xce, 0x6f, 0x66, 0xda, 0x65, 0x97, 0x0a, 0x98, 0x70,
	0xeb, 0xfc, 0xdb, 0xe9, 0xfc, 0xfe, 0x33, 0x03, 0x9d, 0x09, 0x93, 0x53, 0x35, 0x6e, 0x87, 0x3c,
	0x0a, 0x22, 0x2a, 0x26, 0x57, 0x6f, 0x5c, 0xc5, 0xaf, 0x44, 0x32, 0x1e, 0x07, 0x21, 0x4f, 0x69,
	0x40, 0x12, 0x16, 0x08, 0x9a, 0x2e, 0x58, 0x68, 0xce, 0xed, 0x24, 0xe5, 0x92, 0xa3, 0x0a, 0x49,
	0x18, 0x8e, 0xe1, 0xb8, 0x1f, 0x31, 0xd9, 0x5f, 0xd0, 0x58, 0x0e, 0xe9, 0xbb, 0xa2, 0x42, 0xa2,
	0x16, 0xfc, 0x5f, 0xbd, 0x7e, 0x20, 0x62, 0xea, 0x39, 0x2d, 0xe7, 0xd2, 0x1d, 0xda, 0x12, 0xf2,
	0xa1, 0x46, 0x75, 0xc6, 0x80, 0x2e, 0xbd, 0xb2, 0xb9, 0x2e, 0x62, 0x74, 0x06, 0xae, 0x39, 0xf7,
	0x88, 0x24, 0x5e, 0xc5, 0x5c, 0xae, 0x05, 0xdc, 0x81, 0xfa, 0x13, 0x13, 0x92, 0xc6, 0xcf, 0x44,
	0xcc, 0xf6, 0x2e, 0x88, 0x15, 0x34, 0x46, 0x6a, 0x1c, 0x31, 0x39, 0xa4, 0x42, 0xcd, 0x6d, 0x52,
	0xfa, 0x41, 0x43, 0xa5, 0xad, 0x3e, 0xf6, 0xf2, 0x44, 0x4b, 0xd2, 0x34, 0x5c, 0xc9, 0x44, 0x59,
	0xa8, 0x6b, 0x01, 0x9d, 0x03, 0x64, 0x81, 0x05, 0x6b, 0x29, 0xf8, 0x02, 0x8e, 0xac, 0xee, 0x24,
	0xf3, 0x25, 0x3a, 0x81, 0xbf, 0x34, 0x4d, 0x79, 0xba, 0xaa, 0x95, 0x05, 0xf8, 0x13, 0x6a, 0xda,
	0x8f, 0xce, 0xd9, 0x83, 0xa9, 0xf8, 0xa3, 0x6c, 0xfd, 0x81, 0x3c, 0xa8, 0x4a, 0x22, 0x66, 0x9a,
	0x33, 0x03, 0xc9, 0x43, 0xed, 0x81, 0xc5, 0x39, 0xe4, 0x9f, 0xcc, 0x43, 0x21, 0xe0, 0x01, 0xd4,
	0x7f, 0xb6, 0xe6, 0x57, 0xcc, 0x4d, 0xb4, 0xf2, 0x16, 0xda, 0xcd, 0x97, 0x03, 0xd5, 0x51, 0xd6,
	0x77, 0x74, 0x0f, 0x6e, 0x61, 0x1e, 0x9d, 0xb6, 0xf5, 0xe2, 0x6c, 0xae, 0x8a, 0xdf, 0xd8, 0x94,
	0x93, 0xf9, 0x12, 0x97, 0xd0, 0x1d, 0xc0, 0x7a, 0xca, 0xa8, 0x69, 0x1e, 0x6d, 0x8d, 0xdd, 0x3f,
	0x34, 0x7a, 0xde, 0x38, 0x5c, 0xba, 0x76, 0x50, 0x17, 0x0e, 0x6c, 0x33, 0xc8, 0x33, 0x4f, 0x76,
	0x8c, 0xde, 0x6f, 0xee, 0xb8, 0x31, 0xc5, 0xbb, 0xee, 0x4b, 0x75, 0xb5, 0x3a, 0xe3, 0x7f, 0x66,
	0xd3, 0x6f, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x28, 0xa4, 0xf2, 0x22, 0x03, 0x00, 0x00,
}
