// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: process.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Type of condition available to compare the values.
type Process_Node_Filter_Condition_Predicate int32

const (
	Process_Node_Filter_Condition_Unknown Process_Node_Filter_Condition_Predicate = 0
	Process_Node_Filter_Condition_EQ      Process_Node_Filter_Condition_Predicate = 1
)

var Process_Node_Filter_Condition_Predicate_name = map[int32]string{
	0: "Unknown",
	1: "EQ",
}

var Process_Node_Filter_Condition_Predicate_value = map[string]int32{
	"Unknown": 0,
	"EQ":      1,
}

func (x Process_Node_Filter_Condition_Predicate) String() string {
	return proto.EnumName(Process_Node_Filter_Condition_Predicate_name, int32(x))
}

func (Process_Node_Filter_Condition_Predicate) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0, 4, 0, 0}
}

// A process is a configuration to trigger a specific task when certains conditions of a trigger are valid.
type Process struct {
	Hash                 []byte          `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Key                  string          `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Nodes                []*Process_Node `protobuf:"bytes,4,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Edges                []*Process_Edge `protobuf:"bytes,5,rep,name=edges,proto3" json:"edges,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Process) Reset()         { *m = Process{} }
func (m *Process) String() string { return proto.CompactTextString(m) }
func (*Process) ProtoMessage()    {}
func (*Process) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0}
}
func (m *Process) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process.Unmarshal(m, b)
}
func (m *Process) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process.Marshal(b, m, deterministic)
}
func (m *Process) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process.Merge(m, src)
}
func (m *Process) XXX_Size() int {
	return xxx_messageInfo_Process.Size(m)
}
func (m *Process) XXX_DiscardUnknown() {
	xxx_messageInfo_Process.DiscardUnknown(m)
}

var xxx_messageInfo_Process proto.InternalMessageInfo

func (m *Process) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Process) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Process) GetNodes() []*Process_Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func (m *Process) GetEdges() []*Process_Edge {
	if m != nil {
		return m.Edges
	}
	return nil
}

// Node of the process
type Process_Node struct {
	// Types that are valid to be assigned to Type:
	//	*Process_Node_Result_
	//	*Process_Node_Event_
	//	*Process_Node_Task_
	//	*Process_Node_Map_
	//	*Process_Node_Filter_
	Type                 isProcess_Node_Type `protobuf_oneof:"type"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Process_Node) Reset()         { *m = Process_Node{} }
func (m *Process_Node) String() string { return proto.CompactTextString(m) }
func (*Process_Node) ProtoMessage()    {}
func (*Process_Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0}
}
func (m *Process_Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Node.Unmarshal(m, b)
}
func (m *Process_Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Node.Marshal(b, m, deterministic)
}
func (m *Process_Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Node.Merge(m, src)
}
func (m *Process_Node) XXX_Size() int {
	return xxx_messageInfo_Process_Node.Size(m)
}
func (m *Process_Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Node proto.InternalMessageInfo

type isProcess_Node_Type interface {
	isProcess_Node_Type()
}

type Process_Node_Result_ struct {
	Result *Process_Node_Result `protobuf:"bytes,1,opt,name=result,proto3,oneof"`
}
type Process_Node_Event_ struct {
	Event *Process_Node_Event `protobuf:"bytes,2,opt,name=event,proto3,oneof"`
}
type Process_Node_Task_ struct {
	Task *Process_Node_Task `protobuf:"bytes,3,opt,name=task,proto3,oneof"`
}
type Process_Node_Map_ struct {
	Map *Process_Node_Map `protobuf:"bytes,4,opt,name=map,proto3,oneof"`
}
type Process_Node_Filter_ struct {
	Filter *Process_Node_Filter `protobuf:"bytes,5,opt,name=filter,proto3,oneof"`
}

func (*Process_Node_Result_) isProcess_Node_Type() {}
func (*Process_Node_Event_) isProcess_Node_Type()  {}
func (*Process_Node_Task_) isProcess_Node_Type()   {}
func (*Process_Node_Map_) isProcess_Node_Type()    {}
func (*Process_Node_Filter_) isProcess_Node_Type() {}

func (m *Process_Node) GetType() isProcess_Node_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Process_Node) GetResult() *Process_Node_Result {
	if x, ok := m.GetType().(*Process_Node_Result_); ok {
		return x.Result
	}
	return nil
}

func (m *Process_Node) GetEvent() *Process_Node_Event {
	if x, ok := m.GetType().(*Process_Node_Event_); ok {
		return x.Event
	}
	return nil
}

func (m *Process_Node) GetTask() *Process_Node_Task {
	if x, ok := m.GetType().(*Process_Node_Task_); ok {
		return x.Task
	}
	return nil
}

func (m *Process_Node) GetMap() *Process_Node_Map {
	if x, ok := m.GetType().(*Process_Node_Map_); ok {
		return x.Map
	}
	return nil
}

func (m *Process_Node) GetFilter() *Process_Node_Filter {
	if x, ok := m.GetType().(*Process_Node_Filter_); ok {
		return x.Filter
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Process_Node) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Process_Node_OneofMarshaler, _Process_Node_OneofUnmarshaler, _Process_Node_OneofSizer, []interface{}{
		(*Process_Node_Result_)(nil),
		(*Process_Node_Event_)(nil),
		(*Process_Node_Task_)(nil),
		(*Process_Node_Map_)(nil),
		(*Process_Node_Filter_)(nil),
	}
}

func _Process_Node_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Process_Node)
	// type
	switch x := m.Type.(type) {
	case *Process_Node_Result_:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Result); err != nil {
			return err
		}
	case *Process_Node_Event_:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Event); err != nil {
			return err
		}
	case *Process_Node_Task_:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Task); err != nil {
			return err
		}
	case *Process_Node_Map_:
		_ = b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Map); err != nil {
			return err
		}
	case *Process_Node_Filter_:
		_ = b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Filter); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Process_Node.Type has unexpected type %T", x)
	}
	return nil
}

func _Process_Node_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Process_Node)
	switch tag {
	case 1: // type.result
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Process_Node_Result)
		err := b.DecodeMessage(msg)
		m.Type = &Process_Node_Result_{msg}
		return true, err
	case 2: // type.event
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Process_Node_Event)
		err := b.DecodeMessage(msg)
		m.Type = &Process_Node_Event_{msg}
		return true, err
	case 3: // type.task
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Process_Node_Task)
		err := b.DecodeMessage(msg)
		m.Type = &Process_Node_Task_{msg}
		return true, err
	case 4: // type.map
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Process_Node_Map)
		err := b.DecodeMessage(msg)
		m.Type = &Process_Node_Map_{msg}
		return true, err
	case 5: // type.filter
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Process_Node_Filter)
		err := b.DecodeMessage(msg)
		m.Type = &Process_Node_Filter_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Process_Node_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Process_Node)
	// type
	switch x := m.Type.(type) {
	case *Process_Node_Result_:
		s := proto.Size(x.Result)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Process_Node_Event_:
		s := proto.Size(x.Event)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Process_Node_Task_:
		s := proto.Size(x.Task)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Process_Node_Map_:
		s := proto.Size(x.Map)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Process_Node_Filter_:
		s := proto.Size(x.Filter)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Process_Node_Result struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	InstanceHash         []byte   `protobuf:"bytes,2,opt,name=instanceHash,proto3" json:"instanceHash,omitempty"`
	TaskKey              string   `protobuf:"bytes,3,opt,name=taskKey,proto3" json:"taskKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Process_Node_Result) Reset()         { *m = Process_Node_Result{} }
func (m *Process_Node_Result) String() string { return proto.CompactTextString(m) }
func (*Process_Node_Result) ProtoMessage()    {}
func (*Process_Node_Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0, 0}
}
func (m *Process_Node_Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Node_Result.Unmarshal(m, b)
}
func (m *Process_Node_Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Node_Result.Marshal(b, m, deterministic)
}
func (m *Process_Node_Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Node_Result.Merge(m, src)
}
func (m *Process_Node_Result) XXX_Size() int {
	return xxx_messageInfo_Process_Node_Result.Size(m)
}
func (m *Process_Node_Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Node_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Node_Result proto.InternalMessageInfo

func (m *Process_Node_Result) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Process_Node_Result) GetInstanceHash() []byte {
	if m != nil {
		return m.InstanceHash
	}
	return nil
}

func (m *Process_Node_Result) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

type Process_Node_Event struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	InstanceHash         []byte   `protobuf:"bytes,2,opt,name=instanceHash,proto3" json:"instanceHash,omitempty"`
	EventKey             string   `protobuf:"bytes,3,opt,name=eventKey,proto3" json:"eventKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Process_Node_Event) Reset()         { *m = Process_Node_Event{} }
func (m *Process_Node_Event) String() string { return proto.CompactTextString(m) }
func (*Process_Node_Event) ProtoMessage()    {}
func (*Process_Node_Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0, 1}
}
func (m *Process_Node_Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Node_Event.Unmarshal(m, b)
}
func (m *Process_Node_Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Node_Event.Marshal(b, m, deterministic)
}
func (m *Process_Node_Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Node_Event.Merge(m, src)
}
func (m *Process_Node_Event) XXX_Size() int {
	return xxx_messageInfo_Process_Node_Event.Size(m)
}
func (m *Process_Node_Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Node_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Node_Event proto.InternalMessageInfo

func (m *Process_Node_Event) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Process_Node_Event) GetInstanceHash() []byte {
	if m != nil {
		return m.InstanceHash
	}
	return nil
}

func (m *Process_Node_Event) GetEventKey() string {
	if m != nil {
		return m.EventKey
	}
	return ""
}

type Process_Node_Task struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	InstanceHash         []byte   `protobuf:"bytes,2,opt,name=instanceHash,proto3" json:"instanceHash,omitempty"`
	TaskKey              string   `protobuf:"bytes,3,opt,name=taskKey,proto3" json:"taskKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Process_Node_Task) Reset()         { *m = Process_Node_Task{} }
func (m *Process_Node_Task) String() string { return proto.CompactTextString(m) }
func (*Process_Node_Task) ProtoMessage()    {}
func (*Process_Node_Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0, 2}
}
func (m *Process_Node_Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Node_Task.Unmarshal(m, b)
}
func (m *Process_Node_Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Node_Task.Marshal(b, m, deterministic)
}
func (m *Process_Node_Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Node_Task.Merge(m, src)
}
func (m *Process_Node_Task) XXX_Size() int {
	return xxx_messageInfo_Process_Node_Task.Size(m)
}
func (m *Process_Node_Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Node_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Node_Task proto.InternalMessageInfo

func (m *Process_Node_Task) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Process_Node_Task) GetInstanceHash() []byte {
	if m != nil {
		return m.InstanceHash
	}
	return nil
}

func (m *Process_Node_Task) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

type Process_Node_Map struct {
	Key                  string                     `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Outputs              []*Process_Node_Map_Output `protobuf:"bytes,2,rep,name=outputs,proto3" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *Process_Node_Map) Reset()         { *m = Process_Node_Map{} }
func (m *Process_Node_Map) String() string { return proto.CompactTextString(m) }
func (*Process_Node_Map) ProtoMessage()    {}
func (*Process_Node_Map) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0, 3}
}
func (m *Process_Node_Map) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Node_Map.Unmarshal(m, b)
}
func (m *Process_Node_Map) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Node_Map.Marshal(b, m, deterministic)
}
func (m *Process_Node_Map) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Node_Map.Merge(m, src)
}
func (m *Process_Node_Map) XXX_Size() int {
	return xxx_messageInfo_Process_Node_Map.Size(m)
}
func (m *Process_Node_Map) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Node_Map.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Node_Map proto.InternalMessageInfo

func (m *Process_Node_Map) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Process_Node_Map) GetOutputs() []*Process_Node_Map_Output {
	if m != nil {
		return m.Outputs
	}
	return nil
}

type Process_Node_Map_Output struct {
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Types that are valid to be assigned to Value:
	//	*Process_Node_Map_Output_Ref
	Value                isProcess_Node_Map_Output_Value `protobuf_oneof:"value"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *Process_Node_Map_Output) Reset()         { *m = Process_Node_Map_Output{} }
func (m *Process_Node_Map_Output) String() string { return proto.CompactTextString(m) }
func (*Process_Node_Map_Output) ProtoMessage()    {}
func (*Process_Node_Map_Output) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0, 3, 0}
}
func (m *Process_Node_Map_Output) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Node_Map_Output.Unmarshal(m, b)
}
func (m *Process_Node_Map_Output) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Node_Map_Output.Marshal(b, m, deterministic)
}
func (m *Process_Node_Map_Output) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Node_Map_Output.Merge(m, src)
}
func (m *Process_Node_Map_Output) XXX_Size() int {
	return xxx_messageInfo_Process_Node_Map_Output.Size(m)
}
func (m *Process_Node_Map_Output) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Node_Map_Output.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Node_Map_Output proto.InternalMessageInfo

type isProcess_Node_Map_Output_Value interface {
	isProcess_Node_Map_Output_Value()
}

type Process_Node_Map_Output_Ref struct {
	Ref *Process_Node_Map_Output_Reference `protobuf:"bytes,2,opt,name=ref,proto3,oneof"`
}

func (*Process_Node_Map_Output_Ref) isProcess_Node_Map_Output_Value() {}

func (m *Process_Node_Map_Output) GetValue() isProcess_Node_Map_Output_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Process_Node_Map_Output) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Process_Node_Map_Output) GetRef() *Process_Node_Map_Output_Reference {
	if x, ok := m.GetValue().(*Process_Node_Map_Output_Ref); ok {
		return x.Ref
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Process_Node_Map_Output) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Process_Node_Map_Output_OneofMarshaler, _Process_Node_Map_Output_OneofUnmarshaler, _Process_Node_Map_Output_OneofSizer, []interface{}{
		(*Process_Node_Map_Output_Ref)(nil),
	}
}

func _Process_Node_Map_Output_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Process_Node_Map_Output)
	// value
	switch x := m.Value.(type) {
	case *Process_Node_Map_Output_Ref:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Ref); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Process_Node_Map_Output.Value has unexpected type %T", x)
	}
	return nil
}

func _Process_Node_Map_Output_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Process_Node_Map_Output)
	switch tag {
	case 2: // value.ref
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Process_Node_Map_Output_Reference)
		err := b.DecodeMessage(msg)
		m.Value = &Process_Node_Map_Output_Ref{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Process_Node_Map_Output_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Process_Node_Map_Output)
	// value
	switch x := m.Value.(type) {
	case *Process_Node_Map_Output_Ref:
		s := proto.Size(x.Ref)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Process_Node_Map_Output_Reference struct {
	NodeKey              string   `protobuf:"bytes,1,opt,name=nodeKey,proto3" json:"nodeKey,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Process_Node_Map_Output_Reference) Reset()         { *m = Process_Node_Map_Output_Reference{} }
func (m *Process_Node_Map_Output_Reference) String() string { return proto.CompactTextString(m) }
func (*Process_Node_Map_Output_Reference) ProtoMessage()    {}
func (*Process_Node_Map_Output_Reference) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0, 3, 0, 0}
}
func (m *Process_Node_Map_Output_Reference) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Node_Map_Output_Reference.Unmarshal(m, b)
}
func (m *Process_Node_Map_Output_Reference) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Node_Map_Output_Reference.Marshal(b, m, deterministic)
}
func (m *Process_Node_Map_Output_Reference) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Node_Map_Output_Reference.Merge(m, src)
}
func (m *Process_Node_Map_Output_Reference) XXX_Size() int {
	return xxx_messageInfo_Process_Node_Map_Output_Reference.Size(m)
}
func (m *Process_Node_Map_Output_Reference) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Node_Map_Output_Reference.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Node_Map_Output_Reference proto.InternalMessageInfo

func (m *Process_Node_Map_Output_Reference) GetNodeKey() string {
	if m != nil {
		return m.NodeKey
	}
	return ""
}

func (m *Process_Node_Map_Output_Reference) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type Process_Node_Filter struct {
	Key                  string                           `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Conditions           []*Process_Node_Filter_Condition `protobuf:"bytes,2,rep,name=conditions,proto3" json:"conditions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *Process_Node_Filter) Reset()         { *m = Process_Node_Filter{} }
func (m *Process_Node_Filter) String() string { return proto.CompactTextString(m) }
func (*Process_Node_Filter) ProtoMessage()    {}
func (*Process_Node_Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0, 4}
}
func (m *Process_Node_Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Node_Filter.Unmarshal(m, b)
}
func (m *Process_Node_Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Node_Filter.Marshal(b, m, deterministic)
}
func (m *Process_Node_Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Node_Filter.Merge(m, src)
}
func (m *Process_Node_Filter) XXX_Size() int {
	return xxx_messageInfo_Process_Node_Filter.Size(m)
}
func (m *Process_Node_Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Node_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Node_Filter proto.InternalMessageInfo

func (m *Process_Node_Filter) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Process_Node_Filter) GetConditions() []*Process_Node_Filter_Condition {
	if m != nil {
		return m.Conditions
	}
	return nil
}

type Process_Node_Filter_Condition struct {
	Key                  string                                  `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Predicate            Process_Node_Filter_Condition_Predicate `protobuf:"varint,2,opt,name=predicate,proto3,enum=types.Process_Node_Filter_Condition_Predicate" json:"predicate,omitempty"`
	Value                string                                  `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                `json:"-"`
	XXX_unrecognized     []byte                                  `json:"-"`
	XXX_sizecache        int32                                   `json:"-"`
}

func (m *Process_Node_Filter_Condition) Reset()         { *m = Process_Node_Filter_Condition{} }
func (m *Process_Node_Filter_Condition) String() string { return proto.CompactTextString(m) }
func (*Process_Node_Filter_Condition) ProtoMessage()    {}
func (*Process_Node_Filter_Condition) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 0, 4, 0}
}
func (m *Process_Node_Filter_Condition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Node_Filter_Condition.Unmarshal(m, b)
}
func (m *Process_Node_Filter_Condition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Node_Filter_Condition.Marshal(b, m, deterministic)
}
func (m *Process_Node_Filter_Condition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Node_Filter_Condition.Merge(m, src)
}
func (m *Process_Node_Filter_Condition) XXX_Size() int {
	return xxx_messageInfo_Process_Node_Filter_Condition.Size(m)
}
func (m *Process_Node_Filter_Condition) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Node_Filter_Condition.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Node_Filter_Condition proto.InternalMessageInfo

func (m *Process_Node_Filter_Condition) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Process_Node_Filter_Condition) GetPredicate() Process_Node_Filter_Condition_Predicate {
	if m != nil {
		return m.Predicate
	}
	return Process_Node_Filter_Condition_Unknown
}

func (m *Process_Node_Filter_Condition) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Process_Edge struct {
	Src                  string   `protobuf:"bytes,1,opt,name=src,proto3" json:"src,omitempty"`
	Dst                  string   `protobuf:"bytes,2,opt,name=dst,proto3" json:"dst,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Process_Edge) Reset()         { *m = Process_Edge{} }
func (m *Process_Edge) String() string { return proto.CompactTextString(m) }
func (*Process_Edge) ProtoMessage()    {}
func (*Process_Edge) Descriptor() ([]byte, []int) {
	return fileDescriptor_54c4d0e8c0aaf5c3, []int{0, 1}
}
func (m *Process_Edge) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Process_Edge.Unmarshal(m, b)
}
func (m *Process_Edge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Process_Edge.Marshal(b, m, deterministic)
}
func (m *Process_Edge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Process_Edge.Merge(m, src)
}
func (m *Process_Edge) XXX_Size() int {
	return xxx_messageInfo_Process_Edge.Size(m)
}
func (m *Process_Edge) XXX_DiscardUnknown() {
	xxx_messageInfo_Process_Edge.DiscardUnknown(m)
}

var xxx_messageInfo_Process_Edge proto.InternalMessageInfo

func (m *Process_Edge) GetSrc() string {
	if m != nil {
		return m.Src
	}
	return ""
}

func (m *Process_Edge) GetDst() string {
	if m != nil {
		return m.Dst
	}
	return ""
}

func init() {
	proto.RegisterEnum("types.Process_Node_Filter_Condition_Predicate", Process_Node_Filter_Condition_Predicate_name, Process_Node_Filter_Condition_Predicate_value)
	proto.RegisterType((*Process)(nil), "types.Process")
	proto.RegisterType((*Process_Node)(nil), "types.Process.Node")
	proto.RegisterType((*Process_Node_Result)(nil), "types.Process.Node.Result")
	proto.RegisterType((*Process_Node_Event)(nil), "types.Process.Node.Event")
	proto.RegisterType((*Process_Node_Task)(nil), "types.Process.Node.Task")
	proto.RegisterType((*Process_Node_Map)(nil), "types.Process.Node.Map")
	proto.RegisterType((*Process_Node_Map_Output)(nil), "types.Process.Node.Map.Output")
	proto.RegisterType((*Process_Node_Map_Output_Reference)(nil), "types.Process.Node.Map.Output.Reference")
	proto.RegisterType((*Process_Node_Filter)(nil), "types.Process.Node.Filter")
	proto.RegisterType((*Process_Node_Filter_Condition)(nil), "types.Process.Node.Filter.Condition")
	proto.RegisterType((*Process_Edge)(nil), "types.Process.Edge")
}

func init() { proto.RegisterFile("process.proto", fileDescriptor_54c4d0e8c0aaf5c3) }

var fileDescriptor_54c4d0e8c0aaf5c3 = []byte{
	// 570 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xcf, 0x8e, 0xd3, 0x3c,
	0x14, 0xc5, 0x93, 0xe6, 0x4f, 0x27, 0xb7, 0xf3, 0x7d, 0x1a, 0x19, 0x24, 0x42, 0x16, 0xa8, 0xaa,
	0x58, 0x14, 0x10, 0x2e, 0x14, 0x24, 0x58, 0xb0, 0x2a, 0x14, 0x55, 0x1a, 0x06, 0x06, 0x0b, 0x10,
	0xb0, 0x4b, 0x93, 0xdb, 0x36, 0x6a, 0x6b, 0x47, 0xb1, 0x33, 0xa8, 0x2f, 0xc0, 0x9e, 0x2d, 0x5b,
	0x9e, 0x85, 0xd7, 0xe1, 0x19, 0x90, 0x9d, 0xa6, 0x05, 0x91, 0x22, 0x84, 0xd8, 0xe5, 0x3a, 0xbf,
	0x7b, 0x7c, 0xce, 0xb5, 0x13, 0xf8, 0x2f, 0x2f, 0x44, 0x82, 0x52, 0xd2, 0xbc, 0x10, 0x4a, 0x10,
	0x4f, 0x6d, 0x72, 0x94, 0xbd, 0xaf, 0x01, 0xb4, 0xcf, 0xab, 0x17, 0x84, 0x80, 0xbb, 0x88, 0xe5,
	0x22, 0xb4, 0xbb, 0x76, 0xff, 0x98, 0x99, 0x67, 0x72, 0x02, 0xce, 0x12, 0x37, 0x61, 0xab, 0x6b,
	0xf7, 0x03, 0xa6, 0x1f, 0xc9, 0x0d, 0xf0, 0xb8, 0x48, 0x51, 0x86, 0x6e, 0xd7, 0xe9, 0x77, 0x86,
	0x97, 0xa8, 0x11, 0xa2, 0x5b, 0x11, 0xfa, 0x5c, 0xa4, 0xc8, 0x2a, 0x42, 0xa3, 0x98, 0xce, 0x51,
	0x86, 0x5e, 0x23, 0x3a, 0x4e, 0xe7, 0xc8, 0x2a, 0x22, 0xfa, 0x74, 0x04, 0xae, 0x6e, 0x25, 0xf7,
	0xc1, 0x2f, 0x50, 0x96, 0x2b, 0x65, 0x6c, 0x74, 0x86, 0x51, 0x83, 0x3e, 0x65, 0x86, 0x98, 0x58,
	0x6c, 0xcb, 0x92, 0xbb, 0xe0, 0xe1, 0x05, 0x72, 0x65, 0x8c, 0x76, 0x86, 0x57, 0x9b, 0x9a, 0xc6,
	0x1a, 0x98, 0x58, 0xac, 0x22, 0x09, 0x05, 0x57, 0xc5, 0x72, 0x19, 0x3a, 0xa6, 0x23, 0x6c, 0xea,
	0x78, 0x15, 0xcb, 0xe5, 0xc4, 0x62, 0x86, 0x23, 0xb7, 0xc0, 0x59, 0xc7, 0x79, 0xe8, 0x1a, 0xfc,
	0x4a, 0x13, 0x7e, 0x16, 0xe7, 0x13, 0x8b, 0x69, 0x4a, 0xa7, 0x98, 0x65, 0x2b, 0x85, 0x45, 0xe8,
	0x1d, 0x4e, 0xf1, 0xd4, 0x10, 0x3a, 0x45, 0xc5, 0x46, 0x6f, 0xc1, 0xaf, 0x92, 0xd5, 0x63, 0xb7,
	0xf7, 0x63, 0xef, 0xc1, 0x71, 0xc6, 0xa5, 0x8a, 0x79, 0x82, 0x13, 0x7d, 0x48, 0x2d, 0x73, 0x48,
	0x3f, 0xad, 0x91, 0x10, 0xda, 0xda, 0xea, 0x29, 0x6e, 0x4c, 0xaa, 0x80, 0xd5, 0x65, 0xf4, 0x0e,
	0x3c, 0x13, 0xff, 0x2f, 0x85, 0x23, 0x38, 0x32, 0x43, 0xdb, 0x2b, 0xef, 0xea, 0xe8, 0x0d, 0xb8,
	0x7a, 0x4e, 0xff, 0xdc, 0xf2, 0x37, 0x1b, 0x9c, 0xb3, 0x38, 0x6f, 0xd0, 0x7d, 0x08, 0x6d, 0x51,
	0xaa, 0xbc, 0x54, 0x32, 0x6c, 0x99, 0x8b, 0x75, 0xed, 0xc0, 0x69, 0xd0, 0x17, 0x06, 0x63, 0x35,
	0x1e, 0x7d, 0xb6, 0xc1, 0xaf, 0xd6, 0x1a, 0x64, 0x1f, 0x81, 0x53, 0xe0, 0x6c, 0x7b, 0x83, 0xfa,
	0xbf, 0x97, 0xa4, 0x0c, 0x67, 0x58, 0xa0, 0x4e, 0x61, 0x31, 0xdd, 0x16, 0x3d, 0x80, 0x60, 0xb7,
	0xa6, 0x53, 0xe9, 0x2f, 0xe0, 0x74, 0xb7, 0x41, 0x5d, 0xfe, 0xfa, 0x3d, 0x8d, 0xda, 0xe0, 0x5d,
	0xc4, 0xab, 0x12, 0xa3, 0x8f, 0x2d, 0xf0, 0xab, 0x2b, 0xd1, 0x60, 0xee, 0x09, 0x40, 0x22, 0x78,
	0x9a, 0xa9, 0x4c, 0xf0, 0x3a, 0xf6, 0xf5, 0xc3, 0x97, 0x8a, 0x3e, 0xae, 0x61, 0xf6, 0x43, 0x5f,
	0xf4, 0xc5, 0x86, 0x60, 0xf7, 0xa6, 0x61, 0x97, 0x67, 0x10, 0xe4, 0x05, 0xa6, 0x59, 0x12, 0x2b,
	0x34, 0x1e, 0xff, 0x1f, 0xd2, 0x3f, 0xd9, 0x84, 0x9e, 0xd7, 0x5d, 0x6c, 0x2f, 0x40, 0x2e, 0x6f,
	0x93, 0x6d, 0x4f, 0xb6, 0x2a, 0x7a, 0x5d, 0x08, 0x76, 0x34, 0xe9, 0x40, 0xfb, 0x35, 0x5f, 0x72,
	0xf1, 0x81, 0x9f, 0x58, 0xc4, 0x87, 0xd6, 0xf8, 0xe5, 0x89, 0x3d, 0xf2, 0xc1, 0xd5, 0x7b, 0x46,
	0x37, 0xc1, 0xd5, 0xbf, 0x08, 0xed, 0x53, 0x16, 0x49, 0xed, 0x53, 0x16, 0x89, 0x5e, 0x49, 0xa5,
	0xaa, 0xa7, 0x98, 0x4a, 0x35, 0x1a, 0xbe, 0xbf, 0x33, 0xcf, 0xd4, 0xa2, 0x9c, 0xd2, 0x44, 0xac,
	0x07, 0x6b, 0x94, 0xf3, 0xdb, 0x33, 0x51, 0xf2, 0x34, 0xd6, 0xf6, 0x06, 0xc8, 0xe7, 0x19, 0xc7,
	0x81, 0xf9, 0xf3, 0x4d, 0xcb, 0xd9, 0xc0, 0x44, 0x9a, 0xfa, 0xa6, 0xbe, 0xf7, 0x3d, 0x00, 0x00,
	0xff, 0xff, 0xdc, 0xfc, 0x14, 0xae, 0x1a, 0x05, 0x00, 0x00,
}
