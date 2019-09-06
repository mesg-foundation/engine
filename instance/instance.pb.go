// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: instance.proto

package instance

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

// Instance represents service's instance.
type Instance struct {
	Hash                 github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,1,opt,name=hash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"hash"`
	ServiceHash          github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,2,opt,name=serviceHash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"serviceHash"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *Instance) Reset()         { *m = Instance{} }
func (m *Instance) String() string { return proto.CompactTextString(m) }
func (*Instance) ProtoMessage()    {}
func (*Instance) Descriptor() ([]byte, []int) {
	return fileDescriptor_fd22322185b2070b, []int{0}
}
func (m *Instance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Instance.Unmarshal(m, b)
}
func (m *Instance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Instance.Marshal(b, m, deterministic)
}
func (m *Instance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Instance.Merge(m, src)
}
func (m *Instance) XXX_Size() int {
	return xxx_messageInfo_Instance.Size(m)
}
func (m *Instance) XXX_DiscardUnknown() {
	xxx_messageInfo_Instance.DiscardUnknown(m)
}

var xxx_messageInfo_Instance proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Instance)(nil), "types.Instance")
}

func init() { proto.RegisterFile("instance.proto", fileDescriptor_fd22322185b2070b) }

var fileDescriptor_fd22322185b2070b = []byte{
	// 182 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0xcc, 0x2b, 0x2e,
	0x49, 0xcc, 0x4b, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48,
	0x2d, 0x96, 0x52, 0x4a, 0xcf, 0x4f, 0xcf, 0xd7, 0x07, 0x0b, 0x25, 0x95, 0xa6, 0xe9, 0x83, 0x78,
	0x60, 0x0e, 0x98, 0x05, 0x51, 0xaa, 0xb4, 0x8a, 0x91, 0x8b, 0xc3, 0x13, 0xaa, 0x5b, 0xc8, 0x9d,
	0x8b, 0x25, 0x23, 0xb1, 0x38, 0x43, 0x82, 0x51, 0x81, 0x51, 0x83, 0xc7, 0xc9, 0xf8, 0xc4, 0x3d,
	0x79, 0x86, 0x5b, 0xf7, 0xe4, 0xb5, 0xd3, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73,
	0xf5, 0x73, 0x53, 0x8b, 0xd3, 0x75, 0xd3, 0xf2, 0x4b, 0xf3, 0x52, 0x12, 0x4b, 0x32, 0xf3, 0xf3,
	0xf4, 0x53, 0xf3, 0xd2, 0x33, 0xf3, 0x52, 0xf5, 0x41, 0xba, 0xf4, 0x3c, 0x12, 0x8b, 0x33, 0x82,
	0xc0, 0x06, 0x08, 0x85, 0x72, 0x71, 0x17, 0xa7, 0x16, 0x95, 0x65, 0x26, 0xa7, 0x82, 0x04, 0x25,
	0x98, 0xc8, 0x37, 0x0f, 0xd9, 0x1c, 0x27, 0x83, 0x13, 0x0f, 0xe5, 0x18, 0xa2, 0xb4, 0x08, 0xeb,
	0x87, 0x85, 0x47, 0x12, 0x1b, 0xd8, 0x97, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x13, 0xee,
	0xc2, 0xda, 0x22, 0x01, 0x00, 0x00,
}
