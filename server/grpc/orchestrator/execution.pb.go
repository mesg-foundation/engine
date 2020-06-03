// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: server/grpc/orchestrator/execution.proto

package orchestrator

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	execution "github.com/mesg-foundation/engine/execution"
	github_com_mesg_foundation_engine_hash "github.com/mesg-foundation/engine/hash"
	types "github.com/mesg-foundation/engine/protobuf/types"
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

// ExecutionCreateRequest is the request of the endpoint ExecutionCreate.
type ExecutionCreateRequest struct {
	TaskKey              string                                      `protobuf:"bytes,1,opt,name=taskKey,proto3" json:"taskKey,omitempty" validate:"required,printascii"`
	Inputs               *types.Struct                               `protobuf:"bytes,2,opt,name=inputs,proto3" json:"inputs,omitempty"`
	Tags                 []string                                    `protobuf:"bytes,3,rep,name=tags,proto3" json:"tags,omitempty" validate:"dive,printascii"`
	ExecutorHash         github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,4,opt,name=executorHash,proto3,casttype=github.com/mesg-foundation/engine/hash.Hash" json:"executorHash,omitempty" validate:"omitempty,hash"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *ExecutionCreateRequest) Reset()         { *m = ExecutionCreateRequest{} }
func (m *ExecutionCreateRequest) String() string { return proto.CompactTextString(m) }
func (*ExecutionCreateRequest) ProtoMessage()    {}
func (*ExecutionCreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9756b90bba6be877, []int{0}
}
func (m *ExecutionCreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecutionCreateRequest.Unmarshal(m, b)
}
func (m *ExecutionCreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecutionCreateRequest.Marshal(b, m, deterministic)
}
func (m *ExecutionCreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecutionCreateRequest.Merge(m, src)
}
func (m *ExecutionCreateRequest) XXX_Size() int {
	return xxx_messageInfo_ExecutionCreateRequest.Size(m)
}
func (m *ExecutionCreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecutionCreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExecutionCreateRequest proto.InternalMessageInfo

func (m *ExecutionCreateRequest) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

func (m *ExecutionCreateRequest) GetInputs() *types.Struct {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *ExecutionCreateRequest) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *ExecutionCreateRequest) GetExecutorHash() github_com_mesg_foundation_engine_hash.Hash {
	if m != nil {
		return m.ExecutorHash
	}
	return nil
}

// ExecutionCreateResponse is the response of the endpoint ExecutionCreate.
type ExecutionCreateResponse struct {
	// Execution's hash.
	Hash                 github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,1,opt,name=hash,proto3,casttype=github.com/mesg-foundation/engine/hash.Hash" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *ExecutionCreateResponse) Reset()         { *m = ExecutionCreateResponse{} }
func (m *ExecutionCreateResponse) String() string { return proto.CompactTextString(m) }
func (*ExecutionCreateResponse) ProtoMessage()    {}
func (*ExecutionCreateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9756b90bba6be877, []int{1}
}
func (m *ExecutionCreateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecutionCreateResponse.Unmarshal(m, b)
}
func (m *ExecutionCreateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecutionCreateResponse.Marshal(b, m, deterministic)
}
func (m *ExecutionCreateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecutionCreateResponse.Merge(m, src)
}
func (m *ExecutionCreateResponse) XXX_Size() int {
	return xxx_messageInfo_ExecutionCreateResponse.Size(m)
}
func (m *ExecutionCreateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecutionCreateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ExecutionCreateResponse proto.InternalMessageInfo

func (m *ExecutionCreateResponse) GetHash() github_com_mesg_foundation_engine_hash.Hash {
	if m != nil {
		return m.Hash
	}
	return nil
}

// ExecutionStreamRequest defines request to retrieve a stream of executions.
type ExecutionStreamRequest struct {
	// Filter used to filter a stream of executions.
	Filter               *ExecutionStreamRequest_Filter `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *ExecutionStreamRequest) Reset()         { *m = ExecutionStreamRequest{} }
func (m *ExecutionStreamRequest) String() string { return proto.CompactTextString(m) }
func (*ExecutionStreamRequest) ProtoMessage()    {}
func (*ExecutionStreamRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9756b90bba6be877, []int{2}
}
func (m *ExecutionStreamRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecutionStreamRequest.Unmarshal(m, b)
}
func (m *ExecutionStreamRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecutionStreamRequest.Marshal(b, m, deterministic)
}
func (m *ExecutionStreamRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecutionStreamRequest.Merge(m, src)
}
func (m *ExecutionStreamRequest) XXX_Size() int {
	return xxx_messageInfo_ExecutionStreamRequest.Size(m)
}
func (m *ExecutionStreamRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecutionStreamRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExecutionStreamRequest proto.InternalMessageInfo

func (m *ExecutionStreamRequest) GetFilter() *ExecutionStreamRequest_Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

// Filter contains filtering criteria.
type ExecutionStreamRequest_Filter struct {
	// Statuses to filter executions. One status needs to be present in the execution.
	Statuses []execution.Status `protobuf:"varint,1,rep,packed,name=statuses,proto3,enum=mesg.types.Status" json:"statuses,omitempty"`
	// Instance's hash to filter executions.
	InstanceHash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,2,opt,name=instanceHash,proto3,casttype=github.com/mesg-foundation/engine/hash.Hash" json:"instanceHash,omitempty" validate:"omitempty,hash"`
	// taskKey to filter executions.
	TaskKey string `protobuf:"bytes,3,opt,name=taskKey,proto3" json:"taskKey,omitempty" validate:"printascii"`
	// tags to filter executions. All tags needs to be present in the execution.
	Tags []string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty" validate:"dive,printascii"`
	// Executor's hash to filter executions.
	ExecutorHash         github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,5,opt,name=executorHash,proto3,casttype=github.com/mesg-foundation/engine/hash.Hash" json:"executorHash,omitempty" validate:"omitempty,hash"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *ExecutionStreamRequest_Filter) Reset()         { *m = ExecutionStreamRequest_Filter{} }
func (m *ExecutionStreamRequest_Filter) String() string { return proto.CompactTextString(m) }
func (*ExecutionStreamRequest_Filter) ProtoMessage()    {}
func (*ExecutionStreamRequest_Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_9756b90bba6be877, []int{2, 0}
}
func (m *ExecutionStreamRequest_Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecutionStreamRequest_Filter.Unmarshal(m, b)
}
func (m *ExecutionStreamRequest_Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecutionStreamRequest_Filter.Marshal(b, m, deterministic)
}
func (m *ExecutionStreamRequest_Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecutionStreamRequest_Filter.Merge(m, src)
}
func (m *ExecutionStreamRequest_Filter) XXX_Size() int {
	return xxx_messageInfo_ExecutionStreamRequest_Filter.Size(m)
}
func (m *ExecutionStreamRequest_Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecutionStreamRequest_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_ExecutionStreamRequest_Filter proto.InternalMessageInfo

func (m *ExecutionStreamRequest_Filter) GetStatuses() []execution.Status {
	if m != nil {
		return m.Statuses
	}
	return nil
}

func (m *ExecutionStreamRequest_Filter) GetInstanceHash() github_com_mesg_foundation_engine_hash.Hash {
	if m != nil {
		return m.InstanceHash
	}
	return nil
}

func (m *ExecutionStreamRequest_Filter) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

func (m *ExecutionStreamRequest_Filter) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *ExecutionStreamRequest_Filter) GetExecutorHash() github_com_mesg_foundation_engine_hash.Hash {
	if m != nil {
		return m.ExecutorHash
	}
	return nil
}

func init() {
	proto.RegisterType((*ExecutionCreateRequest)(nil), "mesg.grpc.orchestrator.ExecutionCreateRequest")
	proto.RegisterType((*ExecutionCreateResponse)(nil), "mesg.grpc.orchestrator.ExecutionCreateResponse")
	proto.RegisterType((*ExecutionStreamRequest)(nil), "mesg.grpc.orchestrator.ExecutionStreamRequest")
	proto.RegisterType((*ExecutionStreamRequest_Filter)(nil), "mesg.grpc.orchestrator.ExecutionStreamRequest.Filter")
}

func init() {
	proto.RegisterFile("server/grpc/orchestrator/execution.proto", fileDescriptor_9756b90bba6be877)
}

var fileDescriptor_9756b90bba6be877 = []byte{
	// 534 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0xcb, 0x6e, 0xd3, 0x40,
	0x14, 0xad, 0xe3, 0x60, 0xc8, 0x34, 0x62, 0x31, 0x52, 0x4a, 0x6a, 0xa4, 0x24, 0xf2, 0xca, 0x12,
	0xed, 0x0c, 0x4a, 0xc5, 0x86, 0xc7, 0x26, 0x15, 0x08, 0x09, 0xd8, 0x38, 0x3b, 0x16, 0x48, 0x13,
	0xe7, 0xc6, 0x19, 0xb5, 0xf1, 0xb8, 0x33, 0xd7, 0x11, 0xf9, 0x00, 0xbe, 0x86, 0x6f, 0xe1, 0x17,
	0xf2, 0x0b, 0x48, 0x59, 0xb2, 0x42, 0x1e, 0xd7, 0x79, 0xb4, 0x91, 0x5a, 0x84, 0xba, 0xcb, 0xe4,
	0x9e, 0x73, 0xcf, 0xcc, 0xb9, 0xc7, 0x97, 0x84, 0x06, 0xf4, 0x1c, 0x34, 0x4f, 0x74, 0x16, 0x73,
	0xa5, 0xe3, 0x29, 0x18, 0xd4, 0x02, 0x95, 0xe6, 0xf0, 0x1d, 0xe2, 0x1c, 0xa5, 0x4a, 0x59, 0xa6,
	0x15, 0x2a, 0x7a, 0x34, 0x03, 0x93, 0xb0, 0x02, 0xc7, 0xb6, 0x71, 0x7e, 0x90, 0xa8, 0x44, 0x71,
	0x8b, 0x19, 0xe5, 0x13, 0x5e, 0x9c, 0xec, 0xc1, 0xfe, 0x2a, 0xb9, 0x7e, 0x67, 0x5d, 0xc6, 0x45,
	0x06, 0xe6, 0x66, 0x6f, 0xff, 0xf9, 0x8d, 0xba, 0x41, 0x9d, 0xc7, 0x58, 0x16, 0x83, 0x9f, 0x35,
	0x72, 0xf4, 0xbe, 0x22, 0x9c, 0x6b, 0x10, 0x08, 0x11, 0x5c, 0xe5, 0x60, 0x90, 0xbe, 0x25, 0x8f,
	0x51, 0x98, 0x8b, 0x4f, 0xb0, 0x68, 0x3b, 0x3d, 0x27, 0x6c, 0x0c, 0x82, 0xd5, 0xb2, 0xdb, 0x99,
	0x8b, 0x4b, 0x39, 0x16, 0x08, 0xaf, 0x03, 0x0d, 0x57, 0xb9, 0xd4, 0x30, 0x3e, 0xc9, 0xb4, 0x4c,
	0x51, 0x98, 0x58, 0xca, 0x20, 0xaa, 0x28, 0xf4, 0x94, 0x78, 0x32, 0xcd, 0x72, 0x34, 0xed, 0x5a,
	0xcf, 0x09, 0x0f, 0xfb, 0x2d, 0x66, 0x9f, 0x58, 0xdd, 0x85, 0x0d, 0xed, 0x2d, 0xa2, 0x6b, 0x10,
	0xed, 0x93, 0x3a, 0x8a, 0xc4, 0xb4, 0xdd, 0x9e, 0x1b, 0x36, 0x06, 0x9d, 0xd5, 0xb2, 0xeb, 0x6f,
	0x94, 0xc6, 0x72, 0x0e, 0x3b, 0x2a, 0x16, 0x4b, 0x33, 0xd2, 0x2c, 0xdf, 0xaa, 0xf4, 0x47, 0x61,
	0xa6, 0xed, 0x7a, 0xcf, 0x09, 0x9b, 0x83, 0xcf, 0xab, 0x65, 0xf7, 0x78, 0xc3, 0x55, 0x33, 0x89,
	0x30, 0xcb, 0x70, 0x71, 0x32, 0x15, 0x66, 0x1a, 0xfc, 0x59, 0x76, 0x5f, 0x24, 0x12, 0xa7, 0xf9,
	0x88, 0xc5, 0x6a, 0xc6, 0x8b, 0x3b, 0x9d, 0x4e, 0x54, 0x9e, 0x8e, 0x45, 0x61, 0x00, 0x87, 0x34,
	0x91, 0x29, 0xf0, 0x02, 0xca, 0x8a, 0x9e, 0xd1, 0x8e, 0x42, 0xf0, 0x8d, 0x3c, 0xbb, 0x65, 0x96,
	0xc9, 0x54, 0x6a, 0x80, 0x9e, 0x93, 0x7a, 0xc1, 0xb2, 0x56, 0x35, 0x07, 0xfc, 0x5f, 0x75, 0x2c,
	0x39, 0xf8, 0xed, 0x6e, 0x4d, 0x63, 0x88, 0x1a, 0xc4, 0xac, 0x9a, 0xc6, 0x17, 0xe2, 0x4d, 0xe4,
	0x25, 0x82, 0xb6, 0x0a, 0x87, 0xfd, 0x57, 0x6c, 0x7f, 0x64, 0xd8, 0x7e, 0x3e, 0xfb, 0x60, 0xc9,
	0xd1, 0x75, 0x13, 0xff, 0x87, 0x4b, 0xbc, 0xf2, 0x2f, 0xca, 0xc8, 0x13, 0x83, 0x02, 0x73, 0x03,
	0xa6, 0xed, 0xf4, 0xdc, 0xf0, 0x69, 0x9f, 0x96, 0xbd, 0x6d, 0x5c, 0xd8, 0xd0, 0xd6, 0xa2, 0x35,
	0xa6, 0xb0, 0x5d, 0xa6, 0x06, 0x45, 0x1a, 0x83, 0xb5, 0xbd, 0xf6, 0x10, 0xb6, 0x6f, 0x2b, 0xd0,
	0xb3, 0x4d, 0x12, 0x5d, 0x9b, 0xc4, 0xe3, 0xd5, 0xb2, 0xdb, 0xda, 0x88, 0xed, 0x0d, 0x60, 0x95,
	0xa8, 0xfa, 0x7f, 0x24, 0xea, 0xd1, 0x43, 0x27, 0xaa, 0xff, 0xcb, 0x21, 0x8d, 0xf5, 0xc4, 0xe8,
	0x05, 0xf1, 0xca, 0x58, 0x51, 0x76, 0xe7, 0x78, 0x77, 0x3e, 0x56, 0x9f, 0xdf, 0x1b, 0x5f, 0xe6,
	0x35, 0x38, 0xa0, 0x43, 0xe2, 0x95, 0x11, 0xb9, 0x87, 0xd8, 0x4e, 0x96, 0xfc, 0xd6, 0x76, 0x3e,
	0xd6, 0x98, 0xe0, 0xe0, 0xa5, 0x33, 0x78, 0xf7, 0xf5, 0xcd, 0xdd, 0x66, 0x14, 0x6b, 0x51, 0xc6,
	0x70, 0x7b, 0x2f, 0x8e, 0x3c, 0xbb, 0x1f, 0xce, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff, 0xbb, 0xa3,
	0x52, 0xc4, 0x3a, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ExecutionClient is the client API for Execution service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ExecutionClient interface {
	// Create an execution on the blockchain.
	Create(ctx context.Context, in *ExecutionCreateRequest, opts ...grpc.CallOption) (*ExecutionCreateResponse, error)
	// Stream returns a stream of executions that satisfy specified filter.
	Stream(ctx context.Context, in *ExecutionStreamRequest, opts ...grpc.CallOption) (Execution_StreamClient, error)
}

type executionClient struct {
	cc *grpc.ClientConn
}

func NewExecutionClient(cc *grpc.ClientConn) ExecutionClient {
	return &executionClient{cc}
}

func (c *executionClient) Create(ctx context.Context, in *ExecutionCreateRequest, opts ...grpc.CallOption) (*ExecutionCreateResponse, error) {
	out := new(ExecutionCreateResponse)
	err := c.cc.Invoke(ctx, "/mesg.grpc.orchestrator.Execution/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executionClient) Stream(ctx context.Context, in *ExecutionStreamRequest, opts ...grpc.CallOption) (Execution_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Execution_serviceDesc.Streams[0], "/mesg.grpc.orchestrator.Execution/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &executionStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Execution_StreamClient interface {
	Recv() (*execution.Execution, error)
	grpc.ClientStream
}

type executionStreamClient struct {
	grpc.ClientStream
}

func (x *executionStreamClient) Recv() (*execution.Execution, error) {
	m := new(execution.Execution)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExecutionServer is the server API for Execution service.
type ExecutionServer interface {
	// Create an execution on the blockchain.
	Create(context.Context, *ExecutionCreateRequest) (*ExecutionCreateResponse, error)
	// Stream returns a stream of executions that satisfy specified filter.
	Stream(*ExecutionStreamRequest, Execution_StreamServer) error
}

// UnimplementedExecutionServer can be embedded to have forward compatible implementations.
type UnimplementedExecutionServer struct {
}

func (*UnimplementedExecutionServer) Create(ctx context.Context, req *ExecutionCreateRequest) (*ExecutionCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedExecutionServer) Stream(req *ExecutionStreamRequest, srv Execution_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}

func RegisterExecutionServer(s *grpc.Server, srv ExecutionServer) {
	s.RegisterService(&_Execution_serviceDesc, srv)
}

func _Execution_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecutionCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutionServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.grpc.orchestrator.Execution/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutionServer).Create(ctx, req.(*ExecutionCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Execution_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExecutionStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExecutionServer).Stream(m, &executionStreamServer{stream})
}

type Execution_StreamServer interface {
	Send(*execution.Execution) error
	grpc.ServerStream
}

type executionStreamServer struct {
	grpc.ServerStream
}

func (x *executionStreamServer) Send(m *execution.Execution) error {
	return x.ServerStream.SendMsg(m)
}

var _Execution_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mesg.grpc.orchestrator.Execution",
	HandlerType: (*ExecutionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Execution_Create_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Execution_Stream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "server/grpc/orchestrator/execution.proto",
}
