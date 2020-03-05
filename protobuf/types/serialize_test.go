package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	structSort1 = &Struct{
		Fields: map[string]*Value{
			"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
			"1":      {Kind: &Value_StringValue{StringValue: "value"}},
			"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
			"number": {Kind: &Value_NumberValue{NumberValue: 10}},
			"string": {Kind: &Value_StringValue{StringValue: "value"}},
			"bool":   {Kind: &Value_BoolValue{BoolValue: true}},
			"struct": {Kind: &Value_StructValue{StructValue: &Struct{Fields: map[string]*Value{
				"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
				"1":      {Kind: &Value_StringValue{StringValue: "value"}},
				"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
				"number": {Kind: &Value_NumberValue{NumberValue: 10}},
				"struct": {Kind: &Value_StructValue{StructValue: &Struct{Fields: map[string]*Value{
					"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
					"1":      {Kind: &Value_StringValue{StringValue: "value"}},
					"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
					"number": {Kind: &Value_NumberValue{NumberValue: 10}},
				}}}},
			}}}},
			"list": {Kind: &Value_ListValue{ListValue: &ListValue{Values: []*Value{
				{Kind: &Value_StringValue{StringValue: "value"}},
			}}}},
		},
	}
	structSort2 = &Struct{
		Fields: map[string]*Value{
			"1": {Kind: &Value_StringValue{StringValue: "value"}},
			"struct": {Kind: &Value_StructValue{StructValue: &Struct{Fields: map[string]*Value{
				"1":    {Kind: &Value_StringValue{StringValue: "value"}},
				"null": {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
				"struct": {Kind: &Value_StructValue{StructValue: &Struct{Fields: map[string]*Value{
					"number": {Kind: &Value_NumberValue{NumberValue: 10}},
					"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
					"1":      {Kind: &Value_StringValue{StringValue: "value"}},
					"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
				}}}},
				"number": {Kind: &Value_NumberValue{NumberValue: 10}},
				"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
			}}}},
			"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
			"number": {Kind: &Value_NumberValue{NumberValue: 10}},
			"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
			"bool":   {Kind: &Value_BoolValue{BoolValue: true}},
			"string": {Kind: &Value_StringValue{StringValue: "value"}},
			"list": {Kind: &Value_ListValue{ListValue: &ListValue{Values: []*Value{
				{Kind: &Value_StringValue{StringValue: "value"}},
			}}}},
		},
	}
)

func TestHashSerialize(t *testing.T) {
	require.Equal(t, "1:1:3:value;;10000:2:10;;bool:4:true;;list:6:1:0:3:value;;;;;number:2:10;;string:3:value;;struct:5:1:1:3:value;;10000:2:10;;number:2:10;;struct:5:1:1:3:value;;10000:2:10;;number:2:10;;;;;;;;;", structSort1.HashSerialize())
	require.Equal(t, structSort1.HashSerialize(), structSort2.HashSerialize())
}

func BenchmarkHashSerialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		structSort1.HashSerialize()
	}
}
