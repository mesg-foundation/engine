// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/types/workflow.proto

package types

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Type of condition available to compare the values.
type Workflow_Node_Filter_Condition_Predicate int32

const (
	Workflow_Node_Filter_Condition_Unknown Workflow_Node_Filter_Condition_Predicate = 0
	Workflow_Node_Filter_Condition_EQ      Workflow_Node_Filter_Condition_Predicate = 1
)

var Workflow_Node_Filter_Condition_Predicate_name = map[int32]string{
	0: "Unknown",
	1: "EQ",
}

var Workflow_Node_Filter_Condition_Predicate_value = map[string]int32{
	"Unknown": 0,
	"EQ":      1,
}

func (x Workflow_Node_Filter_Condition_Predicate) String() string {
	return proto.EnumName(Workflow_Node_Filter_Condition_Predicate_name, int32(x))
}

func (Workflow_Node_Filter_Condition_Predicate) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0, 4, 0, 0}
}

// A workflow is a configuration to trigger a specific task when certains conditions of a trigger are valid.
type Workflow struct {
	Hash                 string           `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Key                  string           `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Nodes                []*Workflow_Node `protobuf:"bytes,4,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Edges                []*Workflow_Edge `protobuf:"bytes,5,rep,name=edges,proto3" json:"edges,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Workflow) Reset()         { *m = Workflow{} }
func (m *Workflow) String() string { return proto.CompactTextString(m) }
func (*Workflow) ProtoMessage()    {}
func (*Workflow) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0}
}

func (m *Workflow) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow.Unmarshal(m, b)
}
func (m *Workflow) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow.Marshal(b, m, deterministic)
}
func (m *Workflow) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow.Merge(m, src)
}
func (m *Workflow) XXX_Size() int {
	return xxx_messageInfo_Workflow.Size(m)
}
func (m *Workflow) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow proto.InternalMessageInfo

func (m *Workflow) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *Workflow) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Workflow) GetNodes() []*Workflow_Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func (m *Workflow) GetEdges() []*Workflow_Edge {
	if m != nil {
		return m.Edges
	}
	return nil
}

// Node of the workflow
type Workflow_Node struct {
	// Types that are valid to be assigned to Type:
	//	*Workflow_Node_Result_
	//	*Workflow_Node_Event_
	//	*Workflow_Node_Task_
	//	*Workflow_Node_Map_
	//	*Workflow_Node_Filter_
	Type                 isWorkflow_Node_Type `protobuf_oneof:"type"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Workflow_Node) Reset()         { *m = Workflow_Node{} }
func (m *Workflow_Node) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node) ProtoMessage()    {}
func (*Workflow_Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0}
}

func (m *Workflow_Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node.Unmarshal(m, b)
}
func (m *Workflow_Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node.Marshal(b, m, deterministic)
}
func (m *Workflow_Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node.Merge(m, src)
}
func (m *Workflow_Node) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node.Size(m)
}
func (m *Workflow_Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node proto.InternalMessageInfo

type isWorkflow_Node_Type interface {
	isWorkflow_Node_Type()
}

type Workflow_Node_Result_ struct {
	Result *Workflow_Node_Result `protobuf:"bytes,1,opt,name=result,proto3,oneof"`
}

type Workflow_Node_Event_ struct {
	Event *Workflow_Node_Event `protobuf:"bytes,2,opt,name=event,proto3,oneof"`
}

type Workflow_Node_Task_ struct {
	Task *Workflow_Node_Task `protobuf:"bytes,3,opt,name=task,proto3,oneof"`
}

type Workflow_Node_Map_ struct {
	Map *Workflow_Node_Map `protobuf:"bytes,4,opt,name=map,proto3,oneof"`
}

type Workflow_Node_Filter_ struct {
	Filter *Workflow_Node_Filter `protobuf:"bytes,5,opt,name=filter,proto3,oneof"`
}

func (*Workflow_Node_Result_) isWorkflow_Node_Type() {}

func (*Workflow_Node_Event_) isWorkflow_Node_Type() {}

func (*Workflow_Node_Task_) isWorkflow_Node_Type() {}

func (*Workflow_Node_Map_) isWorkflow_Node_Type() {}

func (*Workflow_Node_Filter_) isWorkflow_Node_Type() {}

func (m *Workflow_Node) GetType() isWorkflow_Node_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Workflow_Node) GetResult() *Workflow_Node_Result {
	if x, ok := m.GetType().(*Workflow_Node_Result_); ok {
		return x.Result
	}
	return nil
}

func (m *Workflow_Node) GetEvent() *Workflow_Node_Event {
	if x, ok := m.GetType().(*Workflow_Node_Event_); ok {
		return x.Event
	}
	return nil
}

func (m *Workflow_Node) GetTask() *Workflow_Node_Task {
	if x, ok := m.GetType().(*Workflow_Node_Task_); ok {
		return x.Task
	}
	return nil
}

func (m *Workflow_Node) GetMap() *Workflow_Node_Map {
	if x, ok := m.GetType().(*Workflow_Node_Map_); ok {
		return x.Map
	}
	return nil
}

func (m *Workflow_Node) GetFilter() *Workflow_Node_Filter {
	if x, ok := m.GetType().(*Workflow_Node_Filter_); ok {
		return x.Filter
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Workflow_Node) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Workflow_Node_Result_)(nil),
		(*Workflow_Node_Event_)(nil),
		(*Workflow_Node_Task_)(nil),
		(*Workflow_Node_Map_)(nil),
		(*Workflow_Node_Filter_)(nil),
	}
}

type Workflow_Node_Result struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	InstanceHash         string   `protobuf:"bytes,2,opt,name=instanceHash,proto3" json:"instanceHash,omitempty"`
	TaskKey              string   `protobuf:"bytes,3,opt,name=taskKey,proto3" json:"taskKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Workflow_Node_Result) Reset()         { *m = Workflow_Node_Result{} }
func (m *Workflow_Node_Result) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node_Result) ProtoMessage()    {}
func (*Workflow_Node_Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0, 0}
}

func (m *Workflow_Node_Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node_Result.Unmarshal(m, b)
}
func (m *Workflow_Node_Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node_Result.Marshal(b, m, deterministic)
}
func (m *Workflow_Node_Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node_Result.Merge(m, src)
}
func (m *Workflow_Node_Result) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node_Result.Size(m)
}
func (m *Workflow_Node_Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node_Result proto.InternalMessageInfo

func (m *Workflow_Node_Result) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Workflow_Node_Result) GetInstanceHash() string {
	if m != nil {
		return m.InstanceHash
	}
	return ""
}

func (m *Workflow_Node_Result) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

type Workflow_Node_Event struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	InstanceHash         string   `protobuf:"bytes,2,opt,name=instanceHash,proto3" json:"instanceHash,omitempty"`
	EventKey             string   `protobuf:"bytes,3,opt,name=eventKey,proto3" json:"eventKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Workflow_Node_Event) Reset()         { *m = Workflow_Node_Event{} }
func (m *Workflow_Node_Event) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node_Event) ProtoMessage()    {}
func (*Workflow_Node_Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0, 1}
}

func (m *Workflow_Node_Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node_Event.Unmarshal(m, b)
}
func (m *Workflow_Node_Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node_Event.Marshal(b, m, deterministic)
}
func (m *Workflow_Node_Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node_Event.Merge(m, src)
}
func (m *Workflow_Node_Event) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node_Event.Size(m)
}
func (m *Workflow_Node_Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node_Event proto.InternalMessageInfo

func (m *Workflow_Node_Event) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Workflow_Node_Event) GetInstanceHash() string {
	if m != nil {
		return m.InstanceHash
	}
	return ""
}

func (m *Workflow_Node_Event) GetEventKey() string {
	if m != nil {
		return m.EventKey
	}
	return ""
}

type Workflow_Node_Task struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	InstanceHash         string   `protobuf:"bytes,2,opt,name=instanceHash,proto3" json:"instanceHash,omitempty"`
	TaskKey              string   `protobuf:"bytes,3,opt,name=taskKey,proto3" json:"taskKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Workflow_Node_Task) Reset()         { *m = Workflow_Node_Task{} }
func (m *Workflow_Node_Task) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node_Task) ProtoMessage()    {}
func (*Workflow_Node_Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0, 2}
}

func (m *Workflow_Node_Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node_Task.Unmarshal(m, b)
}
func (m *Workflow_Node_Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node_Task.Marshal(b, m, deterministic)
}
func (m *Workflow_Node_Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node_Task.Merge(m, src)
}
func (m *Workflow_Node_Task) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node_Task.Size(m)
}
func (m *Workflow_Node_Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node_Task proto.InternalMessageInfo

func (m *Workflow_Node_Task) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Workflow_Node_Task) GetInstanceHash() string {
	if m != nil {
		return m.InstanceHash
	}
	return ""
}

func (m *Workflow_Node_Task) GetTaskKey() string {
	if m != nil {
		return m.TaskKey
	}
	return ""
}

type Workflow_Node_Map struct {
	Key                  string                      `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Outputs              []*Workflow_Node_Map_Output `protobuf:"bytes,2,rep,name=outputs,proto3" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *Workflow_Node_Map) Reset()         { *m = Workflow_Node_Map{} }
func (m *Workflow_Node_Map) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node_Map) ProtoMessage()    {}
func (*Workflow_Node_Map) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0, 3}
}

func (m *Workflow_Node_Map) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node_Map.Unmarshal(m, b)
}
func (m *Workflow_Node_Map) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node_Map.Marshal(b, m, deterministic)
}
func (m *Workflow_Node_Map) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node_Map.Merge(m, src)
}
func (m *Workflow_Node_Map) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node_Map.Size(m)
}
func (m *Workflow_Node_Map) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node_Map.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node_Map proto.InternalMessageInfo

func (m *Workflow_Node_Map) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Workflow_Node_Map) GetOutputs() []*Workflow_Node_Map_Output {
	if m != nil {
		return m.Outputs
	}
	return nil
}

type Workflow_Node_Map_Output struct {
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Types that are valid to be assigned to Value:
	//	*Workflow_Node_Map_Output_Ref
	Value                isWorkflow_Node_Map_Output_Value `protobuf_oneof:"value"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *Workflow_Node_Map_Output) Reset()         { *m = Workflow_Node_Map_Output{} }
func (m *Workflow_Node_Map_Output) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node_Map_Output) ProtoMessage()    {}
func (*Workflow_Node_Map_Output) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0, 3, 0}
}

func (m *Workflow_Node_Map_Output) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node_Map_Output.Unmarshal(m, b)
}
func (m *Workflow_Node_Map_Output) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node_Map_Output.Marshal(b, m, deterministic)
}
func (m *Workflow_Node_Map_Output) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node_Map_Output.Merge(m, src)
}
func (m *Workflow_Node_Map_Output) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node_Map_Output.Size(m)
}
func (m *Workflow_Node_Map_Output) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node_Map_Output.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node_Map_Output proto.InternalMessageInfo

func (m *Workflow_Node_Map_Output) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type isWorkflow_Node_Map_Output_Value interface {
	isWorkflow_Node_Map_Output_Value()
}

type Workflow_Node_Map_Output_Ref struct {
	Ref *Workflow_Node_Map_Output_Reference `protobuf:"bytes,2,opt,name=ref,proto3,oneof"`
}

func (*Workflow_Node_Map_Output_Ref) isWorkflow_Node_Map_Output_Value() {}

func (m *Workflow_Node_Map_Output) GetValue() isWorkflow_Node_Map_Output_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Workflow_Node_Map_Output) GetRef() *Workflow_Node_Map_Output_Reference {
	if x, ok := m.GetValue().(*Workflow_Node_Map_Output_Ref); ok {
		return x.Ref
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Workflow_Node_Map_Output) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Workflow_Node_Map_Output_Ref)(nil),
	}
}

type Workflow_Node_Map_Output_Reference struct {
	NodeKey              string   `protobuf:"bytes,1,opt,name=nodeKey,proto3" json:"nodeKey,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Workflow_Node_Map_Output_Reference) Reset()         { *m = Workflow_Node_Map_Output_Reference{} }
func (m *Workflow_Node_Map_Output_Reference) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node_Map_Output_Reference) ProtoMessage()    {}
func (*Workflow_Node_Map_Output_Reference) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0, 3, 0, 0}
}

func (m *Workflow_Node_Map_Output_Reference) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node_Map_Output_Reference.Unmarshal(m, b)
}
func (m *Workflow_Node_Map_Output_Reference) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node_Map_Output_Reference.Marshal(b, m, deterministic)
}
func (m *Workflow_Node_Map_Output_Reference) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node_Map_Output_Reference.Merge(m, src)
}
func (m *Workflow_Node_Map_Output_Reference) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node_Map_Output_Reference.Size(m)
}
func (m *Workflow_Node_Map_Output_Reference) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node_Map_Output_Reference.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node_Map_Output_Reference proto.InternalMessageInfo

func (m *Workflow_Node_Map_Output_Reference) GetNodeKey() string {
	if m != nil {
		return m.NodeKey
	}
	return ""
}

func (m *Workflow_Node_Map_Output_Reference) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type Workflow_Node_Filter struct {
	Key                  string                            `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Conditions           []*Workflow_Node_Filter_Condition `protobuf:"bytes,2,rep,name=conditions,proto3" json:"conditions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                          `json:"-"`
	XXX_unrecognized     []byte                            `json:"-"`
	XXX_sizecache        int32                             `json:"-"`
}

func (m *Workflow_Node_Filter) Reset()         { *m = Workflow_Node_Filter{} }
func (m *Workflow_Node_Filter) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node_Filter) ProtoMessage()    {}
func (*Workflow_Node_Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0, 4}
}

func (m *Workflow_Node_Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node_Filter.Unmarshal(m, b)
}
func (m *Workflow_Node_Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node_Filter.Marshal(b, m, deterministic)
}
func (m *Workflow_Node_Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node_Filter.Merge(m, src)
}
func (m *Workflow_Node_Filter) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node_Filter.Size(m)
}
func (m *Workflow_Node_Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node_Filter proto.InternalMessageInfo

func (m *Workflow_Node_Filter) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Workflow_Node_Filter) GetConditions() []*Workflow_Node_Filter_Condition {
	if m != nil {
		return m.Conditions
	}
	return nil
}

type Workflow_Node_Filter_Condition struct {
	Key                  string                                   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Predicate            Workflow_Node_Filter_Condition_Predicate `protobuf:"varint,2,opt,name=predicate,proto3,enum=types.Workflow_Node_Filter_Condition_Predicate" json:"predicate,omitempty"`
	Value                string                                   `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                 `json:"-"`
	XXX_unrecognized     []byte                                   `json:"-"`
	XXX_sizecache        int32                                    `json:"-"`
}

func (m *Workflow_Node_Filter_Condition) Reset()         { *m = Workflow_Node_Filter_Condition{} }
func (m *Workflow_Node_Filter_Condition) String() string { return proto.CompactTextString(m) }
func (*Workflow_Node_Filter_Condition) ProtoMessage()    {}
func (*Workflow_Node_Filter_Condition) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 0, 4, 0}
}

func (m *Workflow_Node_Filter_Condition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Node_Filter_Condition.Unmarshal(m, b)
}
func (m *Workflow_Node_Filter_Condition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Node_Filter_Condition.Marshal(b, m, deterministic)
}
func (m *Workflow_Node_Filter_Condition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Node_Filter_Condition.Merge(m, src)
}
func (m *Workflow_Node_Filter_Condition) XXX_Size() int {
	return xxx_messageInfo_Workflow_Node_Filter_Condition.Size(m)
}
func (m *Workflow_Node_Filter_Condition) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Node_Filter_Condition.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Node_Filter_Condition proto.InternalMessageInfo

func (m *Workflow_Node_Filter_Condition) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Workflow_Node_Filter_Condition) GetPredicate() Workflow_Node_Filter_Condition_Predicate {
	if m != nil {
		return m.Predicate
	}
	return Workflow_Node_Filter_Condition_Unknown
}

func (m *Workflow_Node_Filter_Condition) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Workflow_Edge struct {
	Src                  string   `protobuf:"bytes,1,opt,name=src,proto3" json:"src,omitempty"`
	Dst                  string   `protobuf:"bytes,2,opt,name=dst,proto3" json:"dst,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Workflow_Edge) Reset()         { *m = Workflow_Edge{} }
func (m *Workflow_Edge) String() string { return proto.CompactTextString(m) }
func (*Workflow_Edge) ProtoMessage()    {}
func (*Workflow_Edge) Descriptor() ([]byte, []int) {
	return fileDescriptor_980f671c228050a1, []int{0, 1}
}

func (m *Workflow_Edge) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Workflow_Edge.Unmarshal(m, b)
}
func (m *Workflow_Edge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Workflow_Edge.Marshal(b, m, deterministic)
}
func (m *Workflow_Edge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Workflow_Edge.Merge(m, src)
}
func (m *Workflow_Edge) XXX_Size() int {
	return xxx_messageInfo_Workflow_Edge.Size(m)
}
func (m *Workflow_Edge) XXX_DiscardUnknown() {
	xxx_messageInfo_Workflow_Edge.DiscardUnknown(m)
}

var xxx_messageInfo_Workflow_Edge proto.InternalMessageInfo

func (m *Workflow_Edge) GetSrc() string {
	if m != nil {
		return m.Src
	}
	return ""
}

func (m *Workflow_Edge) GetDst() string {
	if m != nil {
		return m.Dst
	}
	return ""
}

func init() {
	proto.RegisterEnum("types.Workflow_Node_Filter_Condition_Predicate", Workflow_Node_Filter_Condition_Predicate_name, Workflow_Node_Filter_Condition_Predicate_value)
	proto.RegisterType((*Workflow)(nil), "types.Workflow")
	proto.RegisterType((*Workflow_Node)(nil), "types.Workflow.Node")
	proto.RegisterType((*Workflow_Node_Result)(nil), "types.Workflow.Node.Result")
	proto.RegisterType((*Workflow_Node_Event)(nil), "types.Workflow.Node.Event")
	proto.RegisterType((*Workflow_Node_Task)(nil), "types.Workflow.Node.Task")
	proto.RegisterType((*Workflow_Node_Map)(nil), "types.Workflow.Node.Map")
	proto.RegisterType((*Workflow_Node_Map_Output)(nil), "types.Workflow.Node.Map.Output")
	proto.RegisterType((*Workflow_Node_Map_Output_Reference)(nil), "types.Workflow.Node.Map.Output.Reference")
	proto.RegisterType((*Workflow_Node_Filter)(nil), "types.Workflow.Node.Filter")
	proto.RegisterType((*Workflow_Node_Filter_Condition)(nil), "types.Workflow.Node.Filter.Condition")
	proto.RegisterType((*Workflow_Edge)(nil), "types.Workflow.Edge")
}

func init() { proto.RegisterFile("protobuf/types/workflow.proto", fileDescriptor_980f671c228050a1) }

var fileDescriptor_980f671c228050a1 = []byte{
	// 573 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xdf, 0x8a, 0xd3, 0x4e,
	0x14, 0xc7, 0x93, 0xe6, 0x4f, 0x37, 0xa7, 0x3f, 0x7e, 0x2c, 0xc3, 0x5e, 0xc4, 0x88, 0xb8, 0x2c,
	0x08, 0xeb, 0xa2, 0x89, 0x44, 0x44, 0xbc, 0xf0, 0xa6, 0x52, 0x29, 0x2c, 0xf5, 0xcf, 0xe0, 0xff,
	0xbb, 0x34, 0x39, 0x69, 0x43, 0xdb, 0x99, 0x90, 0x99, 0x6c, 0xe9, 0x1b, 0xf8, 0x0c, 0x7a, 0xeb,
	0x0b, 0xf9, 0x26, 0x3e, 0x82, 0xcc, 0xa4, 0x69, 0x5d, 0x8c, 0xab, 0x88, 0x77, 0x39, 0x27, 0x9f,
	0x73, 0xce, 0xf7, 0x7b, 0x66, 0x12, 0xb8, 0x51, 0x56, 0x5c, 0xf2, 0x69, 0x9d, 0x47, 0x72, 0x53,
	0xa2, 0x88, 0xd6, 0xbc, 0x5a, 0xe4, 0x4b, 0xbe, 0x0e, 0x75, 0x9e, 0x38, 0x3a, 0x7b, 0xf2, 0xd5,
	0x83, 0x83, 0xb7, 0xdb, 0x37, 0x84, 0x80, 0x3d, 0x4f, 0xc4, 0xdc, 0x37, 0x8f, 0xcd, 0x53, 0x8f,
	0xea, 0x67, 0x72, 0x08, 0xd6, 0x02, 0x37, 0x7e, 0x4f, 0xa7, 0xd4, 0x23, 0x39, 0x03, 0x87, 0xf1,
	0x0c, 0x85, 0x6f, 0x1f, 0x5b, 0xa7, 0x83, 0xf8, 0x28, 0xd4, 0x9d, 0xc2, 0xb6, 0x4b, 0xf8, 0x8c,
	0x67, 0x48, 0x1b, 0x44, 0xb1, 0x98, 0xcd, 0x50, 0xf8, 0x4e, 0x37, 0x3b, 0xca, 0x66, 0x48, 0x1b,
	0x24, 0xf8, 0x74, 0x00, 0xb6, 0xaa, 0x25, 0x0f, 0xc0, 0xad, 0x50, 0xd4, 0x4b, 0xa9, 0x85, 0x0c,
	0xe2, 0xeb, 0x5d, 0x13, 0x42, 0xaa, 0x91, 0xb1, 0x41, 0xb7, 0x30, 0x89, 0xc1, 0xc1, 0x0b, 0x64,
	0x52, 0x6b, 0x1d, 0xc4, 0x41, 0x67, 0xd5, 0x48, 0x11, 0x63, 0x83, 0x36, 0x28, 0x89, 0xc0, 0x96,
	0x89, 0x58, 0xf8, 0x96, 0x2e, 0xb9, 0xd6, 0x59, 0xf2, 0x2a, 0x11, 0x8b, 0xb1, 0x41, 0x35, 0x48,
	0xee, 0x80, 0xb5, 0x4a, 0x4a, 0xdf, 0xd6, 0xbc, 0xdf, 0xc9, 0x4f, 0x92, 0x72, 0x6c, 0x50, 0x85,
	0x29, 0x27, 0x79, 0xb1, 0x94, 0x58, 0xf9, 0xce, 0x15, 0x4e, 0x9e, 0x6a, 0x44, 0x39, 0x69, 0xe0,
	0xe0, 0x1d, 0xb8, 0x8d, 0xbb, 0x76, 0xfb, 0xe6, 0x7e, 0xfb, 0x27, 0xf0, 0x5f, 0xc1, 0x84, 0x4c,
	0x58, 0x8a, 0x63, 0x75, 0x56, 0xcd, 0xc1, 0x5c, 0xca, 0x11, 0x1f, 0xfa, 0x4a, 0xec, 0x39, 0x6e,
	0xb4, 0x31, 0x8f, 0xb6, 0x61, 0xf0, 0x1e, 0x1c, 0xbd, 0x81, 0xbf, 0x6c, 0x1c, 0xc0, 0x81, 0xde,
	0xdb, 0xbe, 0xf3, 0x2e, 0x0e, 0xde, 0x80, 0xad, 0x36, 0xf5, 0xcf, 0x25, 0x7f, 0x33, 0xc1, 0x9a,
	0x24, 0x65, 0x47, 0xdf, 0x47, 0xd0, 0xe7, 0xb5, 0x2c, 0x6b, 0x29, 0xfc, 0x9e, 0xbe, 0x5e, 0x37,
	0x7f, 0x75, 0x1e, 0xe1, 0x73, 0xcd, 0xd1, 0x96, 0x0f, 0x3e, 0x9b, 0xe0, 0x36, 0xb9, 0x8e, 0xbe,
	0x8f, 0xc1, 0xaa, 0x30, 0xdf, 0x5e, 0xa3, 0xdb, 0xbf, 0xe9, 0x19, 0x52, 0xcc, 0xb1, 0x42, 0xe5,
	0xc3, 0xa0, 0xaa, 0x2e, 0x78, 0x08, 0xde, 0x2e, 0xa7, 0x7c, 0xa9, 0x2f, 0xe1, 0x7c, 0x37, 0xa1,
	0x0d, 0x7f, 0xfe, 0xb0, 0x86, 0x7d, 0x70, 0x2e, 0x92, 0x65, 0x8d, 0xc1, 0xc7, 0x1e, 0xb8, 0xcd,
	0xa5, 0xe8, 0x50, 0x37, 0x02, 0x48, 0x39, 0xcb, 0x0a, 0x59, 0x70, 0xd6, 0x1a, 0xbf, 0x75, 0xc5,
	0xbd, 0x0a, 0x9f, 0xb4, 0x34, 0xfd, 0xa1, 0x30, 0xf8, 0x62, 0x82, 0xb7, 0x7b, 0xd3, 0x31, 0x66,
	0x02, 0x5e, 0x59, 0x61, 0x56, 0xa4, 0x89, 0x44, 0x2d, 0xf2, 0xff, 0x38, 0xfa, 0xa3, 0x29, 0xe1,
	0x8b, 0xb6, 0x8c, 0xee, 0x3b, 0x90, 0xa3, 0xad, 0xb7, 0xed, 0xe9, 0x36, 0xc1, 0xc9, 0x31, 0x78,
	0x3b, 0x9a, 0x0c, 0xa0, 0xff, 0x9a, 0x2d, 0x18, 0x5f, 0xb3, 0x43, 0x83, 0xb8, 0xd0, 0x1b, 0xbd,
	0x3c, 0x34, 0x87, 0x2e, 0xd8, 0x6a, 0x68, 0x70, 0x06, 0xb6, 0xfa, 0x57, 0x28, 0xa1, 0xa2, 0x4a,
	0x5b, 0xa1, 0xa2, 0x4a, 0x55, 0x26, 0x13, 0xb2, 0xdd, 0x63, 0x26, 0xe4, 0x30, 0xfe, 0x70, 0x6f,
	0x56, 0xc8, 0x79, 0x3d, 0x0d, 0x53, 0xbe, 0x8a, 0x56, 0x28, 0x66, 0x77, 0x73, 0x5e, 0xb3, 0x2c,
	0x51, 0xf2, 0x22, 0x64, 0xb3, 0x82, 0x61, 0x74, 0xf9, 0xef, 0x38, 0x75, 0x75, 0x7c, 0xff, 0x7b,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xdb, 0x06, 0xc8, 0x17, 0x36, 0x05, 0x00, 0x00,
}
