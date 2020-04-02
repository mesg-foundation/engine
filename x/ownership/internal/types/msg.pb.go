// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: x/ownership/internal/types/msg.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

// The message to withdraw coins from an ownership.
type MsgWithdraw struct {
	// The ownership's owner.
	Owner github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=owner,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"owner,omitempty" validate:"required,accaddress"`
	// Hash of the resource to withdraw from.
	ResourceHash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,2,opt,name=resourceHash,proto3,casttype=github.com/mesg-foundation/engine/hash.Hash" json:"resourceHash,omitempty" validate:"required,hash"`
	// amount to withdraw
	Amount               string   `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount,omitempty" validate:"required,coinsPositiveZero"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgWithdraw) Reset()         { *m = MsgWithdraw{} }
func (m *MsgWithdraw) String() string { return proto.CompactTextString(m) }
func (*MsgWithdraw) ProtoMessage()    {}
func (*MsgWithdraw) Descriptor() ([]byte, []int) {
	return fileDescriptor_875fd6c128efe9e4, []int{0}
}
func (m *MsgWithdraw) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgWithdraw.Unmarshal(m, b)
}
func (m *MsgWithdraw) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgWithdraw.Marshal(b, m, deterministic)
}
func (m *MsgWithdraw) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdraw.Merge(m, src)
}
func (m *MsgWithdraw) XXX_Size() int {
	return xxx_messageInfo_MsgWithdraw.Size(m)
}
func (m *MsgWithdraw) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdraw.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdraw proto.InternalMessageInfo

func (m *MsgWithdraw) GetOwner() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *MsgWithdraw) GetResourceHash() github_com_mesg_foundation_engine_hash.Hash {
	if m != nil {
		return m.ResourceHash
	}
	return nil
}

func (m *MsgWithdraw) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func init() {
	proto.RegisterType((*MsgWithdraw)(nil), "mesg.ownership.types.MsgWithdraw")
}

func init() {
	proto.RegisterFile("x/ownership/internal/types/msg.proto", fileDescriptor_875fd6c128efe9e4)
}

var fileDescriptor_875fd6c128efe9e4 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0xd9, 0xc4, 0x81, 0x75, 0xa7, 0xe2, 0xa1, 0x78, 0xb0, 0xa3, 0x28, 0x4c, 0x74, 0xcd,
	0xc1, 0x9b, 0x88, 0xb0, 0x9d, 0x04, 0x11, 0x64, 0x1e, 0x84, 0xdd, 0xb2, 0xe4, 0x2d, 0x0d, 0xae,
	0x79, 0x33, 0x2f, 0xd9, 0xf4, 0x7f, 0xf2, 0x6f, 0xea, 0x1f, 0xd1, 0xa3, 0x27, 0x69, 0x3a, 0x44,
	0x61, 0xe2, 0x29, 0xf9, 0xc8, 0xf7, 0x7d, 0xbf, 0x24, 0x2f, 0x3a, 0x7d, 0x63, 0xb8, 0x31, 0x60,
	0xa9, 0xd0, 0x2b, 0xa6, 0x8d, 0x03, 0x6b, 0xf8, 0x92, 0xb9, 0xf7, 0x15, 0x10, 0x2b, 0x49, 0xe5,
	0x2b, 0x8b, 0x0e, 0xe3, 0xa3, 0x12, 0x48, 0xe5, 0xdf, 0xc6, 0x3c, 0x9c, 0x1f, 0x67, 0x0a, 0x15,
	0xb2, 0xe0, 0x98, 0xfb, 0x05, 0x6b, 0x54, 0x10, 0x61, 0xd7, 0x26, 0xb3, 0x8f, 0x6e, 0x74, 0xf8,
	0x40, 0xea, 0x59, 0xbb, 0x42, 0x5a, 0xbe, 0x89, 0x75, 0xb4, 0x1f, 0x6a, 0x92, 0xce, 0xa0, 0x33,
	0xec, 0x4f, 0x9e, 0xea, 0x2a, 0x3d, 0x59, 0xf3, 0xa5, 0x96, 0xdc, 0xc1, 0x75, 0x66, 0xe1, 0xd5,
	0x6b, 0x0b, 0xf2, 0x92, 0x0b, 0xc1, 0xa5, 0xb4, 0x40, 0x94, 0x7d, 0x56, 0xe9, 0x48, 0x69, 0x57,
	0xf8, 0x79, 0x2e, 0xb0, 0x64, 0x02, 0xa9, 0x44, 0xda, 0x2e, 0x23, 0x92, 0x2f, 0xed, 0x65, 0xf3,
	0xb1, 0x10, 0xe3, 0x36, 0x31, 0x6d, 0x09, 0x31, 0x46, 0x7d, 0x0b, 0x84, 0xde, 0x0a, 0xb8, 0xe3,
	0x54, 0x24, 0xdd, 0x40, 0xbc, 0xaf, 0xab, 0x34, 0xd9, 0x41, 0x2c, 0x38, 0x15, 0x0d, 0xeb, 0xe2,
	0x07, 0xab, 0x79, 0xf5, 0x68, 0x81, 0xde, 0x48, 0xee, 0x34, 0x1a, 0x06, 0x46, 0x69, 0x03, 0xac,
	0xb1, 0xe6, 0x4d, 0xe5, 0xf4, 0x17, 0x20, 0x1e, 0x47, 0x3d, 0x5e, 0xa2, 0x37, 0x2e, 0xd9, 0x1b,
	0x74, 0x86, 0x07, 0x93, 0xf3, 0xba, 0x4a, 0xcf, 0x76, 0xa0, 0x04, 0x6a, 0x43, 0x8f, 0x48, 0xda,
	0xe9, 0x35, 0xcc, 0xc0, 0x62, 0x36, 0xdd, 0x06, 0x27, 0xb7, 0xb3, 0x9b, 0xff, 0xf9, 0x7f, 0x8f,
	0x6c, 0xde, 0x0b, 0xbf, 0x7e, 0xf5, 0x15, 0x00, 0x00, 0xff, 0xff, 0x0d, 0x13, 0xa3, 0xab, 0xd7,
	0x01, 0x00, 0x00,
}
