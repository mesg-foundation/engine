// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ownership.proto

package ownership

import (
	bytes "bytes"
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

type Ownership_Resource int32

const (
	Ownership_None    Ownership_Resource = 0
	Ownership_Service Ownership_Resource = 1
	Ownership_Process Ownership_Resource = 2
	Ownership_Runner  Ownership_Resource = 3
)

var Ownership_Resource_name = map[int32]string{
	0: "None",
	1: "Service",
	2: "Process",
	3: "Runner",
}

var Ownership_Resource_value = map[string]int32{
	"None":    0,
	"Service": 1,
	"Process": 2,
	"Runner":  3,
}

func (x Ownership_Resource) String() string {
	return proto.EnumName(Ownership_Resource_name, int32(x))
}

func (Ownership_Resource) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_21ae26e0dccf9d04, []int{0, 0}
}

// Ownership is a ownership relation between one owner and a resource.
type Ownership struct {
	// Service's hash.
	Hash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,1,opt,name=hash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"hash" hash:"-" validate:"required"`
	// The owner of the resource.
	Owner string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty" hash:"name:2" validate:"required"`
	// Resource's hash.
	ResourceHash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,3,opt,name=resourceHash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"resourceHash" hash:"name:3"`
	// Resource's type.
	Resource             Ownership_Resource `protobuf:"varint,4,opt,name=resource,proto3,enum=mesg.types.Ownership_Resource" json:"resource,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Ownership) Reset()         { *m = Ownership{} }
func (m *Ownership) String() string { return proto.CompactTextString(m) }
func (*Ownership) ProtoMessage()    {}
func (*Ownership) Descriptor() ([]byte, []int) {
	return fileDescriptor_21ae26e0dccf9d04, []int{0}
}
func (m *Ownership) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ownership.Unmarshal(m, b)
}
func (m *Ownership) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ownership.Marshal(b, m, deterministic)
}
func (m *Ownership) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ownership.Merge(m, src)
}
func (m *Ownership) XXX_Size() int {
	return xxx_messageInfo_Ownership.Size(m)
}
func (m *Ownership) XXX_DiscardUnknown() {
	xxx_messageInfo_Ownership.DiscardUnknown(m)
}

var xxx_messageInfo_Ownership proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("mesg.types.Ownership_Resource", Ownership_Resource_name, Ownership_Resource_value)
	proto.RegisterType((*Ownership)(nil), "mesg.types.Ownership")
}

func init() { proto.RegisterFile("ownership.proto", fileDescriptor_21ae26e0dccf9d04) }

var fileDescriptor_21ae26e0dccf9d04 = []byte{
	// 334 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0xcf, 0x4e, 0xc2, 0x30,
	0x18, 0x67, 0x80, 0x08, 0x9f, 0xa8, 0x4b, 0x4f, 0x8b, 0x31, 0x6c, 0x36, 0x31, 0x21, 0x31, 0x74,
	0x09, 0xc4, 0xcb, 0xbc, 0x11, 0x0f, 0x9e, 0xd4, 0x8c, 0x9b, 0xb7, 0x31, 0x3e, 0xb6, 0x26, 0xd2,
	0x62, 0xbb, 0x61, 0x7c, 0x0e, 0x5f, 0xc2, 0x47, 0xe1, 0x19, 0x3c, 0x2c, 0xd1, 0x47, 0xe0, 0x09,
	0xcc, 0x8a, 0xa0, 0x26, 0x1e, 0x8c, 0xb7, 0xfe, 0xbe, 0xfc, 0xfe, 0xf5, 0x6b, 0xe1, 0x50, 0x3e,
	0x0a, 0x54, 0x3a, 0xe5, 0x73, 0x36, 0x57, 0x32, 0x93, 0x04, 0x66, 0xa8, 0x13, 0x96, 0x3d, 0xcd,
	0x51, 0x1f, 0xd1, 0x44, 0x26, 0xd2, 0x37, 0xf3, 0x71, 0x3e, 0xf5, 0x4b, 0x64, 0x80, 0x39, 0xad,
	0xf9, 0xf4, 0xb9, 0x06, 0xad, 0x9b, 0x8d, 0x07, 0x49, 0xa0, 0x9e, 0x46, 0x3a, 0x75, 0x2c, 0xcf,
	0xea, 0xb6, 0x87, 0xa3, 0x65, 0xe1, 0x56, 0x5e, 0x0b, 0xf7, 0x2c, 0xe1, 0x59, 0x9a, 0x8f, 0x59,
	0x2c, 0x67, 0x7e, 0x69, 0xdf, 0x9b, 0xca, 0x5c, 0x4c, 0xa2, 0x8c, 0x4b, 0xe1, 0xa3, 0x48, 0xb8,
	0x40, 0xbf, 0x54, 0xb1, 0xab, 0x48, 0xa7, 0xab, 0xc2, 0x3d, 0x2e, 0x41, 0x40, 0x7b, 0xd4, 0x5b,
	0x44, 0xf7, 0x7c, 0x12, 0x65, 0x18, 0x50, 0x85, 0x0f, 0x39, 0x57, 0x38, 0xa1, 0xa1, 0x09, 0x20,
	0x17, 0xb0, 0x63, 0x9a, 0x3b, 0x55, 0xcf, 0xea, 0xb6, 0x86, 0xa7, 0xab, 0xc2, 0x3d, 0x59, 0xcb,
	0x44, 0x34, 0xc3, 0xa0, 0xff, 0xbb, 0x76, 0xad, 0x21, 0x29, 0xb4, 0x15, 0x6a, 0x99, 0xab, 0x18,
	0xcb, 0x48, 0xa7, 0x66, 0xda, 0x5e, 0xfe, 0xaf, 0xed, 0xfe, 0xb7, 0xd8, 0x01, 0x0d, 0x7f, 0x38,
	0x93, 0x00, 0x9a, 0x1b, 0xec, 0xd4, 0x3d, 0xab, 0x7b, 0xd0, 0xef, 0xb0, 0xaf, 0x05, 0xb3, 0xed,
	0xe2, 0x58, 0xf8, 0xc9, 0x0a, 0xb7, 0x7c, 0x1a, 0x40, 0x73, 0x33, 0x25, 0x4d, 0xa8, 0x5f, 0x4b,
	0x81, 0x76, 0x85, 0xec, 0xc1, 0xee, 0x08, 0xd5, 0x82, 0xc7, 0x68, 0x5b, 0x25, 0xb8, 0x55, 0x32,
	0x46, 0xad, 0xed, 0x2a, 0x01, 0x68, 0x84, 0xb9, 0x10, 0xa8, 0xec, 0xda, 0xf0, 0x7c, 0xf9, 0xd6,
	0xa9, 0xbc, 0xbc, 0x77, 0xac, 0xbb, 0x3f, 0xdc, 0x66, 0xfb, 0x05, 0xc6, 0x0d, 0xf3, 0xa6, 0x83,
	0x8f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfb, 0xb7, 0x4e, 0x6e, 0x16, 0x02, 0x00, 0x00,
}

func (this *Ownership) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Ownership)
	if !ok {
		that2, ok := that.(Ownership)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Hash.Equal(that1.Hash) {
		return false
	}
	if this.Owner != that1.Owner {
		return false
	}
	if !this.ResourceHash.Equal(that1.ResourceHash) {
		return false
	}
	if this.Resource != that1.Resource {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
