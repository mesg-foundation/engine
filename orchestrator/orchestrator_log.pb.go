// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: orchestrator_log.proto

package orchestrator

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_mesg_foundation_engine_hash "github.com/mesg-foundation/engine/hash"
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

// Severity levels
type OrchestratorLog_Severity int32

const (
	OrchestratorLog_Unknown OrchestratorLog_Severity = 0
	OrchestratorLog_Debug   OrchestratorLog_Severity = 1
	OrchestratorLog_Info    OrchestratorLog_Severity = 2
	OrchestratorLog_Error   OrchestratorLog_Severity = 3
)

var OrchestratorLog_Severity_name = map[int32]string{
	0: "Unknown",
	1: "Debug",
	2: "Info",
	3: "Error",
}

var OrchestratorLog_Severity_value = map[string]int32{
	"Unknown": 0,
	"Debug":   1,
	"Info":    2,
	"Error":   3,
}

func (x OrchestratorLog_Severity) String() string {
	return proto.EnumName(OrchestratorLog_Severity_name, int32(x))
}

func (OrchestratorLog_Severity) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6332e0853eb3f185, []int{0, 0}
}

// OrchestratorLogsResponse is the message send on the Logs stream.
type OrchestratorLog struct {
	// Severity of the log
	Severity OrchestratorLog_Severity `protobuf:"varint,1,opt,name=severity,proto3,enum=mesg.types.OrchestratorLog_Severity" json:"severity,omitempty"`
	// Message of the log
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	// Data of the log
	Data                 *OrchestratorLog_Data `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *OrchestratorLog) Reset()         { *m = OrchestratorLog{} }
func (m *OrchestratorLog) String() string { return proto.CompactTextString(m) }
func (*OrchestratorLog) ProtoMessage()    {}
func (*OrchestratorLog) Descriptor() ([]byte, []int) {
	return fileDescriptor_6332e0853eb3f185, []int{0}
}
func (m *OrchestratorLog) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrchestratorLog.Unmarshal(m, b)
}
func (m *OrchestratorLog) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrchestratorLog.Marshal(b, m, deterministic)
}
func (m *OrchestratorLog) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrchestratorLog.Merge(m, src)
}
func (m *OrchestratorLog) XXX_Size() int {
	return xxx_messageInfo_OrchestratorLog.Size(m)
}
func (m *OrchestratorLog) XXX_DiscardUnknown() {
	xxx_messageInfo_OrchestratorLog.DiscardUnknown(m)
}

var xxx_messageInfo_OrchestratorLog proto.InternalMessageInfo

func (m *OrchestratorLog) GetSeverity() OrchestratorLog_Severity {
	if m != nil {
		return m.Severity
	}
	return OrchestratorLog_Unknown
}

func (m *OrchestratorLog) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *OrchestratorLog) GetData() *OrchestratorLog_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

type OrchestratorLog_Data struct {
	// Hash of the process
	ProcessHash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,1,opt,name=processHash,proto3,casttype=github.com/mesg-foundation/engine/hash.Hash" json:"processHash,omitempty"`
	// Key of the node being executed.
	NodeKey string `protobuf:"bytes,2,opt,name=nodeKey,proto3" json:"nodeKey,omitempty"`
	// Type of the node being executed.
	NodeType string `protobuf:"bytes,3,opt,name=nodeType,proto3" json:"nodeType,omitempty"`
	// Hash of the event that trigger this node. Can be empty.
	EventHash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,4,opt,name=eventHash,proto3,casttype=github.com/mesg-foundation/engine/hash.Hash" json:"eventHash,omitempty"`
	// Hash of the parent execution that trigger this node. Can be empty.
	ParentHash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,5,opt,name=parentHash,proto3,casttype=github.com/mesg-foundation/engine/hash.Hash" json:"parentHash,omitempty"`
	// Hash of the execution created by this process. Can be empty.
	ExecutionHash        github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,6,opt,name=executionHash,proto3,casttype=github.com/mesg-foundation/engine/hash.Hash" json:"executionHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *OrchestratorLog_Data) Reset()         { *m = OrchestratorLog_Data{} }
func (m *OrchestratorLog_Data) String() string { return proto.CompactTextString(m) }
func (*OrchestratorLog_Data) ProtoMessage()    {}
func (*OrchestratorLog_Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_6332e0853eb3f185, []int{0, 0}
}
func (m *OrchestratorLog_Data) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrchestratorLog_Data.Unmarshal(m, b)
}
func (m *OrchestratorLog_Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrchestratorLog_Data.Marshal(b, m, deterministic)
}
func (m *OrchestratorLog_Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrchestratorLog_Data.Merge(m, src)
}
func (m *OrchestratorLog_Data) XXX_Size() int {
	return xxx_messageInfo_OrchestratorLog_Data.Size(m)
}
func (m *OrchestratorLog_Data) XXX_DiscardUnknown() {
	xxx_messageInfo_OrchestratorLog_Data.DiscardUnknown(m)
}

var xxx_messageInfo_OrchestratorLog_Data proto.InternalMessageInfo

func (m *OrchestratorLog_Data) GetProcessHash() github_com_mesg_foundation_engine_hash.Hash {
	if m != nil {
		return m.ProcessHash
	}
	return nil
}

func (m *OrchestratorLog_Data) GetNodeKey() string {
	if m != nil {
		return m.NodeKey
	}
	return ""
}

func (m *OrchestratorLog_Data) GetNodeType() string {
	if m != nil {
		return m.NodeType
	}
	return ""
}

func (m *OrchestratorLog_Data) GetEventHash() github_com_mesg_foundation_engine_hash.Hash {
	if m != nil {
		return m.EventHash
	}
	return nil
}

func (m *OrchestratorLog_Data) GetParentHash() github_com_mesg_foundation_engine_hash.Hash {
	if m != nil {
		return m.ParentHash
	}
	return nil
}

func (m *OrchestratorLog_Data) GetExecutionHash() github_com_mesg_foundation_engine_hash.Hash {
	if m != nil {
		return m.ExecutionHash
	}
	return nil
}

func init() {
	proto.RegisterEnum("mesg.types.OrchestratorLog_Severity", OrchestratorLog_Severity_name, OrchestratorLog_Severity_value)
	proto.RegisterType((*OrchestratorLog)(nil), "mesg.types.OrchestratorLog")
	proto.RegisterType((*OrchestratorLog_Data)(nil), "mesg.types.OrchestratorLog.Data")
}

func init() { proto.RegisterFile("orchestrator_log.proto", fileDescriptor_6332e0853eb3f185) }

var fileDescriptor_6332e0853eb3f185 = []byte{
	// 369 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcf, 0x8b, 0xd3, 0x40,
	0x14, 0xc7, 0xcd, 0x36, 0xbb, 0x9b, 0xbc, 0xfa, 0x23, 0xcc, 0x41, 0x42, 0x2f, 0x86, 0xe2, 0xa1,
	0x20, 0x4e, 0x64, 0x15, 0xbc, 0xca, 0xb2, 0x82, 0xa2, 0xb2, 0x18, 0xed, 0xc5, 0x8b, 0x4c, 0xd2,
	0xd7, 0x49, 0xd0, 0xce, 0x0b, 0x33, 0x93, 0x6a, 0xfe, 0x01, 0xff, 0x4c, 0x6f, 0xfe, 0x13, 0x9e,
	0x64, 0xc6, 0xa6, 0x8d, 0x1e, 0x14, 0x7a, 0x9b, 0xef, 0xe3, 0xfb, 0x3e, 0xf9, 0x40, 0x1e, 0xdc,
	0x25, 0x5d, 0xd5, 0x68, 0xac, 0x16, 0x96, 0xf4, 0xc7, 0xcf, 0x24, 0x79, 0xab, 0xc9, 0x12, 0x83,
	0x0d, 0x1a, 0xc9, 0x6d, 0xdf, 0xa2, 0x99, 0xcd, 0x25, 0x49, 0xca, 0xfd, 0xbc, 0xec, 0xd6, 0xb9,
	0x4b, 0x3e, 0xf8, 0xd7, 0xef, 0xfe, 0xfc, 0x47, 0x08, 0x77, 0xae, 0x47, 0xa8, 0xd7, 0x24, 0xd9,
	0x33, 0x88, 0x0c, 0x6e, 0x51, 0x37, 0xb6, 0x4f, 0x83, 0x2c, 0x58, 0xdc, 0xbe, 0xb8, 0xcf, 0x0f,
	0x58, 0xfe, 0x57, 0x9d, 0xbf, 0xdb, 0x75, 0x8b, 0xfd, 0x16, 0x4b, 0xe1, 0x7c, 0x83, 0xc6, 0x08,
	0x89, 0xe9, 0x49, 0x16, 0x2c, 0xe2, 0x62, 0x88, 0xec, 0x09, 0x84, 0x2b, 0x61, 0x45, 0x3a, 0xc9,
	0x82, 0xc5, 0xf4, 0x22, 0xfb, 0x17, 0xf7, 0x4a, 0x58, 0x51, 0xf8, 0xf6, 0xec, 0xdb, 0x04, 0x42,
	0x17, 0xd9, 0x5b, 0x98, 0xb6, 0x9a, 0x2a, 0x34, 0xe6, 0x85, 0x30, 0xb5, 0xb7, 0xbb, 0x79, 0x99,
	0xff, 0xfc, 0x7e, 0xef, 0x81, 0x6c, 0x6c, 0xdd, 0x95, 0xbc, 0xa2, 0x4d, 0xee, 0x98, 0x0f, 0xd7,
	0xd4, 0xa9, 0x95, 0xb0, 0x0d, 0xa9, 0x1c, 0x95, 0x6c, 0x14, 0xe6, 0xb5, 0x30, 0x35, 0x77, 0x6b,
	0xc5, 0x98, 0xe1, 0x5c, 0x15, 0xad, 0xf0, 0x15, 0xf6, 0x83, 0xeb, 0x2e, 0xb2, 0x19, 0x44, 0xee,
	0xf9, 0xbe, 0x6f, 0xd1, 0xfb, 0xc6, 0xc5, 0x3e, 0xb3, 0x37, 0x10, 0xe3, 0x16, 0x95, 0xf5, 0x1a,
	0xe1, 0x71, 0x1a, 0x07, 0x02, 0xbb, 0x06, 0x68, 0x85, 0x1e, 0x78, 0xa7, 0xc7, 0xf1, 0x46, 0x08,
	0xb6, 0x84, 0x5b, 0xf8, 0x15, 0xab, 0xce, 0x35, 0x3d, 0xf3, 0xec, 0x38, 0xe6, 0x9f, 0x94, 0xf9,
	0x53, 0x88, 0x86, 0xdf, 0xcd, 0xa6, 0x70, 0xbe, 0x54, 0x9f, 0x14, 0x7d, 0x51, 0xc9, 0x0d, 0x16,
	0xc3, 0xe9, 0x15, 0x96, 0x9d, 0x4c, 0x02, 0x16, 0x41, 0xf8, 0x52, 0xad, 0x29, 0x39, 0x71, 0xc3,
	0xe7, 0x5a, 0x93, 0x4e, 0x26, 0x97, 0x8f, 0x3e, 0xf0, 0xff, 0x7f, 0x76, 0x7c, 0xd3, 0xe5, 0x99,
	0x3f, 0xd0, 0xc7, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0xcd, 0x48, 0x86, 0xa7, 0xea, 0x02, 0x00,
	0x00,
}
