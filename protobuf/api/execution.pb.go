// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: protobuf/api/execution.proto

package api

import (
	context "context"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	execution "github.com/mesg-foundation/engine/execution"
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

// CreateExecutionRequest defines request to create a single execution.
type CreateExecutionRequest struct {
	// taskKey to filter executions.
	TaskKey string        `protobuf:"bytes,2,opt,name=taskKey,proto3" json:"taskKey,omitempty" validate:"required,printascii"`
	Inputs  *types.Struct `protobuf:"bytes,3,opt,name=inputs,proto3" json:"inputs,omitempty"`
	// tags the execution.
	Tags         []string                                      `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty" validate:"dive,printascii"`
	ParentHash   github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,5,opt,name=parentHash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"parentHash,omitempty" validate:"omitempty,accaddress"`
	EventHash    github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,6,opt,name=eventHash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"eventHash,omitempty" validate:"omitempty,accaddress"`
	ProcessHash  github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,7,opt,name=processHash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"processHash,omitempty" validate:"omitempty,accaddress"`
	NodeKey      string                                        `protobuf:"bytes,8,opt,name=nodeKey,proto3" json:"nodeKey,omitempty"`
	ExecutorHash github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,9,opt,name=executorHash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"executorHash,omitempty" validate:"omitempty,accaddress"`
	// price of running the execution.
	Price                string   `protobuf:"bytes,10,opt,name=price,proto3" json:"price,omitempty" validate:"coins"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateExecutionRequest) Reset()         { *m = CreateExecutionRequest{} }
func (m *CreateExecutionRequest) String() string { return proto.CompactTextString(m) }
func (*CreateExecutionRequest) ProtoMessage()    {}
func (*CreateExecutionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e2c86581f82f05, []int{0}
}
func (m *CreateExecutionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateExecutionRequest.Unmarshal(m, b)
}
func (m *CreateExecutionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateExecutionRequest.Marshal(b, m, deterministic)
}
func (m *CreateExecutionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateExecutionRequest.Merge(m, src)
}
func (m *CreateExecutionRequest) XXX_Size() int {
	return xxx_messageInfo_CreateExecutionRequest.Size(m)
}
func (m *CreateExecutionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateExecutionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateExecutionRequest proto.InternalMessageInfo

func (m *CreateExecutionRequest) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

func (m *CreateExecutionRequest) GetInputs() *types.Struct {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *CreateExecutionRequest) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *CreateExecutionRequest) GetParentHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.ParentHash
	}
	return nil
}

func (m *CreateExecutionRequest) GetEventHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.EventHash
	}
	return nil
}

func (m *CreateExecutionRequest) GetProcessHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.ProcessHash
	}
	return nil
}

func (m *CreateExecutionRequest) GetNodeKey() string {
	if m != nil {
		return m.NodeKey
	}
	return ""
}

func (m *CreateExecutionRequest) GetExecutorHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.ExecutorHash
	}
	return nil
}

func (m *CreateExecutionRequest) GetPrice() string {
	if m != nil {
		return m.Price
	}
	return ""
}

// CreateExecutionResponse defines response for execution creation.
type CreateExecutionResponse struct {
	// Execution's hash.
	Hash                 github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=hash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *CreateExecutionResponse) Reset()         { *m = CreateExecutionResponse{} }
func (m *CreateExecutionResponse) String() string { return proto.CompactTextString(m) }
func (*CreateExecutionResponse) ProtoMessage()    {}
func (*CreateExecutionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e2c86581f82f05, []int{1}
}
func (m *CreateExecutionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateExecutionResponse.Unmarshal(m, b)
}
func (m *CreateExecutionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateExecutionResponse.Marshal(b, m, deterministic)
}
func (m *CreateExecutionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateExecutionResponse.Merge(m, src)
}
func (m *CreateExecutionResponse) XXX_Size() int {
	return xxx_messageInfo_CreateExecutionResponse.Size(m)
}
func (m *CreateExecutionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateExecutionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateExecutionResponse proto.InternalMessageInfo

func (m *CreateExecutionResponse) GetHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Hash
	}
	return nil
}

// GetExecutionRequest defines request to retrieve a single execution.
type GetExecutionRequest struct {
	// Execution's hash to fetch.
	Hash                 github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=hash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"hash,omitempty" validate:"required,accaddress"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *GetExecutionRequest) Reset()         { *m = GetExecutionRequest{} }
func (m *GetExecutionRequest) String() string { return proto.CompactTextString(m) }
func (*GetExecutionRequest) ProtoMessage()    {}
func (*GetExecutionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e2c86581f82f05, []int{2}
}
func (m *GetExecutionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetExecutionRequest.Unmarshal(m, b)
}
func (m *GetExecutionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetExecutionRequest.Marshal(b, m, deterministic)
}
func (m *GetExecutionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetExecutionRequest.Merge(m, src)
}
func (m *GetExecutionRequest) XXX_Size() int {
	return xxx_messageInfo_GetExecutionRequest.Size(m)
}
func (m *GetExecutionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetExecutionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetExecutionRequest proto.InternalMessageInfo

func (m *GetExecutionRequest) GetHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Hash
	}
	return nil
}

// StreamExecutionRequest defines request to retrieve a stream of executions.
type StreamExecutionRequest struct {
	// Filter used to filter a stream of executions.
	Filter               *StreamExecutionRequest_Filter `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *StreamExecutionRequest) Reset()         { *m = StreamExecutionRequest{} }
func (m *StreamExecutionRequest) String() string { return proto.CompactTextString(m) }
func (*StreamExecutionRequest) ProtoMessage()    {}
func (*StreamExecutionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e2c86581f82f05, []int{3}
}
func (m *StreamExecutionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamExecutionRequest.Unmarshal(m, b)
}
func (m *StreamExecutionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamExecutionRequest.Marshal(b, m, deterministic)
}
func (m *StreamExecutionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamExecutionRequest.Merge(m, src)
}
func (m *StreamExecutionRequest) XXX_Size() int {
	return xxx_messageInfo_StreamExecutionRequest.Size(m)
}
func (m *StreamExecutionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamExecutionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StreamExecutionRequest proto.InternalMessageInfo

func (m *StreamExecutionRequest) GetFilter() *StreamExecutionRequest_Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

// Filter contains filtering criteria.
type StreamExecutionRequest_Filter struct {
	// Statuses to filter executions. One status needs to be present in the execution.
	Statuses []execution.Status `protobuf:"varint,1,rep,packed,name=statuses,proto3,enum=mesg.types.Status" json:"statuses,omitempty"`
	// Instance's hash to filter executions.
	InstanceHash github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,2,opt,name=instanceHash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"instanceHash,omitempty" validate:"omitempty,accaddress"`
	// taskKey to filter executions.
	TaskKey string `protobuf:"bytes,3,opt,name=taskKey,proto3" json:"taskKey,omitempty" validate:"printascii"`
	// tags to filter executions. All tags needs to be present in the execution.
	Tags []string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty" validate:"dive,printascii"`
	// Executor's hash to filter executions.
	ExecutorHash         github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,5,opt,name=executorHash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"executorHash,omitempty" validate:"omitempty,accaddress"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *StreamExecutionRequest_Filter) Reset()         { *m = StreamExecutionRequest_Filter{} }
func (m *StreamExecutionRequest_Filter) String() string { return proto.CompactTextString(m) }
func (*StreamExecutionRequest_Filter) ProtoMessage()    {}
func (*StreamExecutionRequest_Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e2c86581f82f05, []int{3, 0}
}
func (m *StreamExecutionRequest_Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamExecutionRequest_Filter.Unmarshal(m, b)
}
func (m *StreamExecutionRequest_Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamExecutionRequest_Filter.Marshal(b, m, deterministic)
}
func (m *StreamExecutionRequest_Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamExecutionRequest_Filter.Merge(m, src)
}
func (m *StreamExecutionRequest_Filter) XXX_Size() int {
	return xxx_messageInfo_StreamExecutionRequest_Filter.Size(m)
}
func (m *StreamExecutionRequest_Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamExecutionRequest_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_StreamExecutionRequest_Filter proto.InternalMessageInfo

func (m *StreamExecutionRequest_Filter) GetStatuses() []execution.Status {
	if m != nil {
		return m.Statuses
	}
	return nil
}

func (m *StreamExecutionRequest_Filter) GetInstanceHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.InstanceHash
	}
	return nil
}

func (m *StreamExecutionRequest_Filter) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

func (m *StreamExecutionRequest_Filter) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *StreamExecutionRequest_Filter) GetExecutorHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.ExecutorHash
	}
	return nil
}

// UpdateExecutionRequest defines request for execution update.
type UpdateExecutionRequest struct {
	// Hash represents execution.
	Hash github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=hash,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"hash,omitempty" validate:"required,accaddress"`
	// result pass to execution
	//
	// Types that are valid to be assigned to Result:
	//	*UpdateExecutionRequest_Outputs
	//	*UpdateExecutionRequest_Error
	Result               isUpdateExecutionRequest_Result `protobuf_oneof:"result"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *UpdateExecutionRequest) Reset()         { *m = UpdateExecutionRequest{} }
func (m *UpdateExecutionRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateExecutionRequest) ProtoMessage()    {}
func (*UpdateExecutionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e2c86581f82f05, []int{4}
}
func (m *UpdateExecutionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateExecutionRequest.Unmarshal(m, b)
}
func (m *UpdateExecutionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateExecutionRequest.Marshal(b, m, deterministic)
}
func (m *UpdateExecutionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateExecutionRequest.Merge(m, src)
}
func (m *UpdateExecutionRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateExecutionRequest.Size(m)
}
func (m *UpdateExecutionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateExecutionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateExecutionRequest proto.InternalMessageInfo

type isUpdateExecutionRequest_Result interface {
	isUpdateExecutionRequest_Result()
}

type UpdateExecutionRequest_Outputs struct {
	Outputs *types.Struct `protobuf:"bytes,2,opt,name=outputs,proto3,oneof" json:"outputs,omitempty"`
}
type UpdateExecutionRequest_Error struct {
	Error string `protobuf:"bytes,3,opt,name=error,proto3,oneof" json:"error,omitempty"`
}

func (*UpdateExecutionRequest_Outputs) isUpdateExecutionRequest_Result() {}
func (*UpdateExecutionRequest_Error) isUpdateExecutionRequest_Result()   {}

func (m *UpdateExecutionRequest) GetResult() isUpdateExecutionRequest_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *UpdateExecutionRequest) GetHash() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *UpdateExecutionRequest) GetOutputs() *types.Struct {
	if x, ok := m.GetResult().(*UpdateExecutionRequest_Outputs); ok {
		return x.Outputs
	}
	return nil
}

func (m *UpdateExecutionRequest) GetError() string {
	if x, ok := m.GetResult().(*UpdateExecutionRequest_Error); ok {
		return x.Error
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*UpdateExecutionRequest) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*UpdateExecutionRequest_Outputs)(nil),
		(*UpdateExecutionRequest_Error)(nil),
	}
}

// UpdateExecutionResponse defines response for execution update.
type UpdateExecutionResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateExecutionResponse) Reset()         { *m = UpdateExecutionResponse{} }
func (m *UpdateExecutionResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateExecutionResponse) ProtoMessage()    {}
func (*UpdateExecutionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e2c86581f82f05, []int{5}
}
func (m *UpdateExecutionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateExecutionResponse.Unmarshal(m, b)
}
func (m *UpdateExecutionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateExecutionResponse.Marshal(b, m, deterministic)
}
func (m *UpdateExecutionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateExecutionResponse.Merge(m, src)
}
func (m *UpdateExecutionResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateExecutionResponse.Size(m)
}
func (m *UpdateExecutionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateExecutionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateExecutionResponse proto.InternalMessageInfo

// The request's data for the `List` API.
type ListExecutionRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListExecutionRequest) Reset()         { *m = ListExecutionRequest{} }
func (m *ListExecutionRequest) String() string { return proto.CompactTextString(m) }
func (*ListExecutionRequest) ProtoMessage()    {}
func (*ListExecutionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e2c86581f82f05, []int{6}
}
func (m *ListExecutionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListExecutionRequest.Unmarshal(m, b)
}
func (m *ListExecutionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListExecutionRequest.Marshal(b, m, deterministic)
}
func (m *ListExecutionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListExecutionRequest.Merge(m, src)
}
func (m *ListExecutionRequest) XXX_Size() int {
	return xxx_messageInfo_ListExecutionRequest.Size(m)
}
func (m *ListExecutionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListExecutionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListExecutionRequest proto.InternalMessageInfo

// The response's data for the `List` API.
type ListExecutionResponse struct {
	// List of executions that match the request's filters.
	Executions           []*execution.Execution `protobuf:"bytes,1,rep,name=executions,proto3" json:"executions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *ListExecutionResponse) Reset()         { *m = ListExecutionResponse{} }
func (m *ListExecutionResponse) String() string { return proto.CompactTextString(m) }
func (*ListExecutionResponse) ProtoMessage()    {}
func (*ListExecutionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_96e2c86581f82f05, []int{7}
}
func (m *ListExecutionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListExecutionResponse.Unmarshal(m, b)
}
func (m *ListExecutionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListExecutionResponse.Marshal(b, m, deterministic)
}
func (m *ListExecutionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListExecutionResponse.Merge(m, src)
}
func (m *ListExecutionResponse) XXX_Size() int {
	return xxx_messageInfo_ListExecutionResponse.Size(m)
}
func (m *ListExecutionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListExecutionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListExecutionResponse proto.InternalMessageInfo

func (m *ListExecutionResponse) GetExecutions() []*execution.Execution {
	if m != nil {
		return m.Executions
	}
	return nil
}

func init() {
	proto.RegisterType((*CreateExecutionRequest)(nil), "mesg.api.CreateExecutionRequest")
	proto.RegisterType((*CreateExecutionResponse)(nil), "mesg.api.CreateExecutionResponse")
	proto.RegisterType((*GetExecutionRequest)(nil), "mesg.api.GetExecutionRequest")
	proto.RegisterType((*StreamExecutionRequest)(nil), "mesg.api.StreamExecutionRequest")
	proto.RegisterType((*StreamExecutionRequest_Filter)(nil), "mesg.api.StreamExecutionRequest.Filter")
	proto.RegisterType((*UpdateExecutionRequest)(nil), "mesg.api.UpdateExecutionRequest")
	proto.RegisterType((*UpdateExecutionResponse)(nil), "mesg.api.UpdateExecutionResponse")
	proto.RegisterType((*ListExecutionRequest)(nil), "mesg.api.ListExecutionRequest")
	proto.RegisterType((*ListExecutionResponse)(nil), "mesg.api.ListExecutionResponse")
}

func init() { proto.RegisterFile("protobuf/api/execution.proto", fileDescriptor_96e2c86581f82f05) }

var fileDescriptor_96e2c86581f82f05 = []byte{
	// 745 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0x51, 0x6f, 0x12, 0x4b,
	0x14, 0x06, 0x16, 0x16, 0x38, 0x6d, 0x6e, 0x6e, 0xe6, 0xb6, 0x94, 0x72, 0xef, 0x65, 0x71, 0x5f,
	0x24, 0xa6, 0x5d, 0x2c, 0x8d, 0x2f, 0x46, 0x63, 0x8a, 0xa9, 0xad, 0xd1, 0x68, 0xb2, 0xe8, 0x8b,
	0x4f, 0x4e, 0x77, 0xa7, 0x74, 0xd2, 0xb2, 0xb3, 0x9d, 0x99, 0x25, 0xed, 0x8b, 0xff, 0xc2, 0xbf,
	0xe6, 0x8b, 0x09, 0xff, 0x41, 0x1e, 0x7d, 0x30, 0x86, 0x19, 0x60, 0xb7, 0xb0, 0x34, 0x51, 0x43,
	0x9f, 0x60, 0xe6, 0x9c, 0x33, 0xdf, 0x77, 0x66, 0xbe, 0xf3, 0x2d, 0xfc, 0x17, 0x72, 0x26, 0xd9,
	0x49, 0x74, 0xda, 0xc2, 0x21, 0x6d, 0x91, 0x2b, 0xe2, 0x45, 0x92, 0xb2, 0xc0, 0x51, 0xdb, 0xa8,
	0xd4, 0x27, 0xa2, 0xe7, 0xe0, 0x90, 0xd6, 0xec, 0x1e, 0xeb, 0xb1, 0xd6, 0x2c, 0x79, 0xbc, 0x52,
	0x0b, 0xf5, 0x4f, 0x67, 0xd7, 0xfe, 0x9d, 0x85, 0xe5, 0x75, 0x48, 0x44, 0x4b, 0x48, 0x1e, 0x79,
	0x72, 0x12, 0xac, 0xcf, 0x05, 0xe7, 0xa0, 0xec, 0xaf, 0x05, 0xa8, 0x3c, 0xe7, 0x04, 0x4b, 0x72,
	0x38, 0x8d, 0xb8, 0xe4, 0x32, 0x22, 0x42, 0xa2, 0x27, 0x50, 0x94, 0x58, 0x9c, 0xbf, 0x22, 0xd7,
	0xd5, 0x5c, 0x23, 0xdb, 0x2c, 0x77, 0xec, 0xd1, 0xd0, 0xaa, 0x0f, 0xf0, 0x05, 0xf5, 0xb1, 0x24,
	0x8f, 0x6d, 0x4e, 0x2e, 0x23, 0xca, 0x89, 0xbf, 0x13, 0x72, 0x1a, 0x48, 0x2c, 0x3c, 0x4a, 0x6d,
	0x77, 0x5a, 0x82, 0x76, 0xc1, 0xa4, 0x41, 0x18, 0x49, 0x51, 0x35, 0x1a, 0xd9, 0xe6, 0x5a, 0x7b,
	0xd3, 0x51, 0x4d, 0x4d, 0xe9, 0x38, 0x5d, 0xc5, 0xd2, 0x9d, 0x24, 0xa1, 0x36, 0xe4, 0x25, 0xee,
	0x89, 0x6a, 0xbe, 0x61, 0x34, 0xcb, 0x9d, 0xfa, 0x68, 0x68, 0xd5, 0x62, 0x24, 0x9f, 0x0e, 0xc8,
	0x0d, 0x14, 0x95, 0x8b, 0x24, 0x40, 0x88, 0x39, 0x09, 0xe4, 0x31, 0x16, 0x67, 0xd5, 0x42, 0x23,
	0xdb, 0x5c, 0xef, 0xbc, 0x1b, 0x0d, 0x2d, 0x2b, 0xae, 0x64, 0x7d, 0x2a, 0x49, 0x3f, 0x94, 0xd7,
	0x3b, 0xd8, 0xf3, 0xb0, 0xef, 0x73, 0x22, 0x84, 0xfd, 0x7d, 0x68, 0xed, 0xf6, 0xa8, 0x3c, 0x8b,
	0x4e, 0x1c, 0x8f, 0xf5, 0x5b, 0x1e, 0x13, 0x7d, 0x26, 0x26, 0x3f, 0xbb, 0xc2, 0x3f, 0xd7, 0x57,
	0xe5, 0x1c, 0x78, 0xde, 0x81, 0xae, 0x70, 0x13, 0x38, 0x88, 0x43, 0x99, 0x0c, 0xa6, 0xa0, 0xe6,
	0x0a, 0x41, 0x63, 0x18, 0x34, 0x80, 0xb5, 0x90, 0x33, 0x8f, 0x08, 0xa1, 0x50, 0x8b, 0x2b, 0x44,
	0x4d, 0x02, 0xa1, 0x2a, 0x14, 0x03, 0xe6, 0x93, 0xb1, 0x04, 0x4a, 0x63, 0x09, 0xb8, 0xd3, 0x25,
	0xba, 0x82, 0x75, 0x2d, 0x25, 0xc6, 0x15, 0xa5, 0xf2, 0x0a, 0x29, 0xdd, 0x40, 0x42, 0x0f, 0xa0,
	0x10, 0x72, 0xea, 0x91, 0x2a, 0x28, 0x51, 0x6e, 0x8c, 0x86, 0xd6, 0xdf, 0x31, 0xa4, 0xc7, 0x68,
	0x20, 0x6c, 0x57, 0xa7, 0xd8, 0x1f, 0x61, 0x6b, 0x41, 0xdc, 0x22, 0x64, 0x81, 0x20, 0xe8, 0x10,
	0xf2, 0x67, 0x63, 0xe2, 0x59, 0x45, 0x7c, 0xef, 0xd7, 0x59, 0xa9, 0x72, 0xfb, 0x13, 0xfc, 0x73,
	0x44, 0xe4, 0xc2, 0xec, 0xf4, 0x6e, 0x9c, 0xde, 0x5d, 0x32, 0x38, 0x7f, 0x74, 0x2b, 0x1a, 0xff,
	0x9b, 0x01, 0x95, 0xae, 0xe4, 0x04, 0xf7, 0x17, 0x38, 0x3c, 0x03, 0xf3, 0x94, 0x5e, 0x48, 0xc2,
	0x15, 0x8b, 0xb5, 0xf6, 0x7d, 0x67, 0x6a, 0x2b, 0x4e, 0x7a, 0x85, 0xf3, 0x42, 0xa5, 0xbb, 0x93,
	0xb2, 0xda, 0x67, 0x03, 0x4c, 0xbd, 0x85, 0x1c, 0x28, 0x09, 0x89, 0x65, 0x24, 0x88, 0xa8, 0x66,
	0x1b, 0x46, 0xf3, 0xaf, 0x36, 0xd2, 0xa7, 0x69, 0x5a, 0x5d, 0x15, 0x73, 0x67, 0x39, 0x63, 0x79,
	0xd0, 0x40, 0x48, 0x1c, 0x78, 0x44, 0xc9, 0x23, 0xb7, 0x4a, 0x79, 0x24, 0x91, 0xd0, 0x7e, 0xec,
	0x5a, 0x86, 0x12, 0xc8, 0xf6, 0x68, 0x68, 0x6d, 0xc6, 0xa0, 0xa9, 0x66, 0xf5, 0x3b, 0xee, 0x33,
	0x3f, 0x01, 0x85, 0xbb, 0x9a, 0x00, 0xfb, 0x4b, 0x16, 0x2a, 0xef, 0x43, 0x3f, 0xcd, 0xb3, 0xef,
	0x4a, 0x77, 0x68, 0x0f, 0x8a, 0x2c, 0x92, 0xca, 0xdf, 0x73, 0xb7, 0xf8, 0xfb, 0x71, 0xc6, 0x9d,
	0xe6, 0xa1, 0x0a, 0x14, 0x08, 0xe7, 0x8c, 0xeb, 0x77, 0x39, 0xce, 0xb8, 0x7a, 0xd9, 0x29, 0x81,
	0xc9, 0x89, 0x88, 0x2e, 0xa4, 0xbd, 0x0d, 0x5b, 0x0b, 0x7d, 0xe9, 0x71, 0xb5, 0x2b, 0xb0, 0xf1,
	0x9a, 0x8a, 0x85, 0x41, 0xb3, 0xdf, 0xc0, 0xe6, 0xdc, 0xfe, 0x64, 0xbe, 0x1f, 0x01, 0xcc, 0xbe,
	0x75, 0x5a, 0xb3, 0x33, 0x8e, 0xba, 0xa5, 0xb8, 0x24, 0x91, 0xd8, 0xfe, 0x91, 0x83, 0xf2, 0x2c,
	0x82, 0xde, 0x82, 0xa9, 0xfd, 0x03, 0x35, 0xe2, 0xe1, 0x49, 0xff, 0x5c, 0xd6, 0xee, 0xdd, 0x92,
	0x31, 0x69, 0x22, 0x83, 0x9e, 0x82, 0x71, 0x44, 0x24, 0xfa, 0x3f, 0xce, 0x4d, 0x71, 0x8f, 0x5a,
	0x3a, 0x4f, 0x3b, 0x83, 0x5e, 0x42, 0x7e, 0xdc, 0x2d, 0xaa, 0xc7, 0xf5, 0x69, 0xb7, 0x52, 0xb3,
	0x96, 0xc6, 0x67, 0x4c, 0x0e, 0xc1, 0xd4, 0x2e, 0x90, 0x6c, 0x2d, 0xdd, 0x17, 0x96, 0xf2, 0x79,
	0x98, 0x1d, 0xdf, 0x90, 0x7e, 0xb2, 0xe4, 0x31, 0xe9, 0xe2, 0x4c, 0xde, 0xd0, 0xb2, 0x67, 0xce,
	0x74, 0x0a, 0x1f, 0x0c, 0x1c, 0xd2, 0x13, 0x53, 0x09, 0x69, 0xff, 0x67, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x00, 0x46, 0xdf, 0x57, 0x29, 0x09, 0x00, 0x00,
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
	// Create creates a single Execution specified in a request.
	Create(ctx context.Context, in *CreateExecutionRequest, opts ...grpc.CallOption) (*CreateExecutionResponse, error)
	// Get returns a single Execution specified in a request.
	Get(ctx context.Context, in *GetExecutionRequest, opts ...grpc.CallOption) (*execution.Execution, error)
	// List returns all Executions matching the criteria of the request.
	List(ctx context.Context, in *ListExecutionRequest, opts ...grpc.CallOption) (*ListExecutionResponse, error)
	// Stream returns a stream of executions that satisfy criteria
	// specified in a request.
	Stream(ctx context.Context, in *StreamExecutionRequest, opts ...grpc.CallOption) (Execution_StreamClient, error)
	// Update updates execution with outputs or an error.
	Update(ctx context.Context, in *UpdateExecutionRequest, opts ...grpc.CallOption) (*UpdateExecutionResponse, error)
}

type executionClient struct {
	cc *grpc.ClientConn
}

func NewExecutionClient(cc *grpc.ClientConn) ExecutionClient {
	return &executionClient{cc}
}

func (c *executionClient) Create(ctx context.Context, in *CreateExecutionRequest, opts ...grpc.CallOption) (*CreateExecutionResponse, error) {
	out := new(CreateExecutionResponse)
	err := c.cc.Invoke(ctx, "/mesg.api.Execution/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executionClient) Get(ctx context.Context, in *GetExecutionRequest, opts ...grpc.CallOption) (*execution.Execution, error) {
	out := new(execution.Execution)
	err := c.cc.Invoke(ctx, "/mesg.api.Execution/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executionClient) List(ctx context.Context, in *ListExecutionRequest, opts ...grpc.CallOption) (*ListExecutionResponse, error) {
	out := new(ListExecutionResponse)
	err := c.cc.Invoke(ctx, "/mesg.api.Execution/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executionClient) Stream(ctx context.Context, in *StreamExecutionRequest, opts ...grpc.CallOption) (Execution_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Execution_serviceDesc.Streams[0], "/mesg.api.Execution/Stream", opts...)
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

func (c *executionClient) Update(ctx context.Context, in *UpdateExecutionRequest, opts ...grpc.CallOption) (*UpdateExecutionResponse, error) {
	out := new(UpdateExecutionResponse)
	err := c.cc.Invoke(ctx, "/mesg.api.Execution/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExecutionServer is the server API for Execution service.
type ExecutionServer interface {
	// Create creates a single Execution specified in a request.
	Create(context.Context, *CreateExecutionRequest) (*CreateExecutionResponse, error)
	// Get returns a single Execution specified in a request.
	Get(context.Context, *GetExecutionRequest) (*execution.Execution, error)
	// List returns all Executions matching the criteria of the request.
	List(context.Context, *ListExecutionRequest) (*ListExecutionResponse, error)
	// Stream returns a stream of executions that satisfy criteria
	// specified in a request.
	Stream(*StreamExecutionRequest, Execution_StreamServer) error
	// Update updates execution with outputs or an error.
	Update(context.Context, *UpdateExecutionRequest) (*UpdateExecutionResponse, error)
}

// UnimplementedExecutionServer can be embedded to have forward compatible implementations.
type UnimplementedExecutionServer struct {
}

func (*UnimplementedExecutionServer) Create(ctx context.Context, req *CreateExecutionRequest) (*CreateExecutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedExecutionServer) Get(ctx context.Context, req *GetExecutionRequest) (*execution.Execution, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedExecutionServer) List(ctx context.Context, req *ListExecutionRequest) (*ListExecutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (*UnimplementedExecutionServer) Stream(req *StreamExecutionRequest, srv Execution_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (*UnimplementedExecutionServer) Update(ctx context.Context, req *UpdateExecutionRequest) (*UpdateExecutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func RegisterExecutionServer(s *grpc.Server, srv ExecutionServer) {
	s.RegisterService(&_Execution_serviceDesc, srv)
}

func _Execution_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateExecutionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutionServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Execution/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutionServer).Create(ctx, req.(*CreateExecutionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Execution_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetExecutionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutionServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Execution/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutionServer).Get(ctx, req.(*GetExecutionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Execution_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListExecutionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutionServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Execution/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutionServer).List(ctx, req.(*ListExecutionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Execution_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamExecutionRequest)
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

func _Execution_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateExecutionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutionServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mesg.api.Execution/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutionServer).Update(ctx, req.(*UpdateExecutionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Execution_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mesg.api.Execution",
	HandlerType: (*ExecutionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Execution_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Execution_Get_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Execution_List_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Execution_Update_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Execution_Stream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protobuf/api/execution.proto",
}
