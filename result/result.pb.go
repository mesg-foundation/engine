// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: result.proto

package result

import (
	bytes "bytes"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_mesg_foundation_engine_hash "github.com/mesg-foundation/engine/hash"
	types "github.com/mesg-foundation/engine/protobuf/types"
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

// Result represents a single result of an execution.
type Result struct {
	// Hash is a unique hash to identify this execution's result.
	Hash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,1,opt,name=hash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"hash" hash:"-"`
	// requestHash is hash of the associated execution.
	RequestHash github_com_mesg_foundation_engine_hash.Hash `protobuf:"bytes,2,opt,name=requestHash,proto3,customtype=github.com/mesg-foundation/engine/hash.Hash" json:"requestHash" hash:"name:2"`
	// result pass to execution
	//
	// Types that are valid to be assigned to Result:
	//	*Result_Outputs
	//	*Result_Error
	Result               isResult_Result `protobuf_oneof:"result"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_4feee897733d2100, []int{0}
}
func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

type isResult_Result interface {
	isResult_Result()
	Equal(interface{}) bool
}

type Result_Outputs struct {
	Outputs *types.Struct `protobuf:"bytes,3,opt,name=outputs,proto3,oneof" json:"outputs,omitempty" hash:"name:3"`
}
type Result_Error struct {
	Error string `protobuf:"bytes,4,opt,name=error,proto3,oneof" json:"error,omitempty" hash:"name:4"`
}

func (*Result_Outputs) isResult_Result() {}
func (*Result_Error) isResult_Result()   {}

func (m *Result) GetResult() isResult_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *Result) GetOutputs() *types.Struct {
	if x, ok := m.GetResult().(*Result_Outputs); ok {
		return x.Outputs
	}
	return nil
}

func (m *Result) GetError() string {
	if x, ok := m.GetResult().(*Result_Error); ok {
		return x.Error
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Result) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Result_Outputs)(nil),
		(*Result_Error)(nil),
	}
}

func init() {
	proto.RegisterType((*Result)(nil), "mesg.types.Result")
}

func init() { proto.RegisterFile("result.proto", fileDescriptor_4feee897733d2100) }

var fileDescriptor_4feee897733d2100 = []byte{
	// 294 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4a, 0x2d, 0x2e,
	0xcd, 0x29, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xca, 0x4d, 0x2d, 0x4e, 0xd7, 0x2b,
	0xa9, 0x2c, 0x48, 0x2d, 0x96, 0x52, 0x4a, 0xcf, 0x4f, 0xcf, 0xd7, 0x07, 0x8b, 0x27, 0x95, 0xa6,
	0xe9, 0x83, 0x78, 0x60, 0x0e, 0x98, 0x05, 0x51, 0x2f, 0x25, 0x0d, 0x97, 0x06, 0xeb, 0xd1, 0x2f,
	0x2e, 0x29, 0x2a, 0x4d, 0x86, 0x1a, 0xa6, 0x74, 0x80, 0x89, 0x8b, 0x2d, 0x08, 0x6c, 0xba, 0x50,
	0x30, 0x17, 0x4b, 0x46, 0x62, 0x71, 0x86, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x8f, 0x93, 0xfd, 0x89,
	0x7b, 0xf2, 0x0c, 0xb7, 0xee, 0xc9, 0x6b, 0xa7, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7,
	0xe7, 0xea, 0x83, 0x2c, 0xd6, 0x4d, 0xcb, 0x2f, 0xcd, 0x4b, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3,
	0x4f, 0xcd, 0x4b, 0xcf, 0xcc, 0x4b, 0xd5, 0x07, 0xe9, 0xd2, 0xf3, 0x48, 0x2c, 0xce, 0xf8, 0x74,
	0x4f, 0x9e, 0x03, 0xc4, 0xb1, 0x52, 0xd2, 0x55, 0x0a, 0x02, 0x1b, 0x26, 0x94, 0xc6, 0xc5, 0x5d,
	0x94, 0x5a, 0x58, 0x9a, 0x5a, 0x5c, 0x02, 0x52, 0x20, 0xc1, 0x04, 0x36, 0xdb, 0x85, 0x3c, 0xb3,
	0x79, 0x21, 0x66, 0xe7, 0x25, 0xe6, 0xa6, 0x5a, 0x19, 0x29, 0x05, 0x21, 0x1b, 0x2c, 0xe4, 0xc2,
	0xc5, 0x9e, 0x5f, 0x5a, 0x52, 0x50, 0x5a, 0x52, 0x2c, 0xc1, 0xac, 0xc0, 0xa8, 0xc1, 0x6d, 0x24,
	0xaa, 0x07, 0x0e, 0x26, 0x98, 0xdf, 0xf5, 0x82, 0xc1, 0xbe, 0x76, 0x12, 0x44, 0x33, 0xc7, 0x58,
	0xc9, 0x83, 0x21, 0x08, 0xa6, 0x55, 0x48, 0x93, 0x8b, 0x35, 0xb5, 0xa8, 0x28, 0xbf, 0x48, 0x82,
	0x45, 0x81, 0x51, 0x83, 0x13, 0x43, 0xb1, 0x09, 0x48, 0x31, 0x44, 0x85, 0x13, 0x07, 0x17, 0x1b,
	0x24, 0x56, 0x9c, 0x8c, 0x4e, 0x3c, 0x94, 0x63, 0x58, 0xf1, 0x48, 0x8e, 0x31, 0x4a, 0x83, 0xb0,
	0x7f, 0x20, 0x7a, 0x92, 0xd8, 0xc0, 0xee, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x47, 0x2c,
	0x80, 0xf0, 0xda, 0x01, 0x00, 0x00,
}

func (this *Result) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Result)
	if !ok {
		that2, ok := that.(Result)
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
	if !this.RequestHash.Equal(that1.RequestHash) {
		return false
	}
	if that1.Result == nil {
		if this.Result != nil {
			return false
		}
	} else if this.Result == nil {
		return false
	} else if !this.Result.Equal(that1.Result) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *Result_Outputs) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Result_Outputs)
	if !ok {
		that2, ok := that.(Result_Outputs)
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
	if !this.Outputs.Equal(that1.Outputs) {
		return false
	}
	return true
}
func (this *Result_Error) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Result_Error)
	if !ok {
		that2, ok := that.(Result_Error)
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
	if this.Error != that1.Error {
		return false
	}
	return true
}
