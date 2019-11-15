package types

import (
	"testing"

	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
)

func TestStructHash(t *testing.T) {
	hashes := make(map[string]bool)
	structs := []Struct{
		{Fields: map[string]*Value{}},
		{Fields: map[string]*Value{"v": {Kind: &Value_NullValue{}}}},
		{Fields: map[string]*Value{"v": {Kind: &Value_NumberValue{}}}},
		{Fields: map[string]*Value{"v": {Kind: &Value_StringValue{}}}},
		{Fields: map[string]*Value{"v": {Kind: &Value_BoolValue{}}}},
		{Fields: map[string]*Value{"v": {Kind: &Value_StructValue{}}}},
		{Fields: map[string]*Value{"v": {Kind: &Value_ListValue{}}}},
	}
	for _, s := range structs {
		hashes[string(hash.Dump(s))] = true
	}
	require.Equal(t, len(hashes), len(structs))
}

func TestStructMarshal(t *testing.T) {
	var (
		structSort1 = &Struct{
			Fields: map[string]*Value{
				"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
				"1":      {Kind: &Value_StringValue{StringValue: "valuea"}},
				"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
				"number": {Kind: &Value_NumberValue{NumberValue: 10}},
				"string": {Kind: &Value_StringValue{StringValue: "valuea"}},
				"bool":   {Kind: &Value_BoolValue{BoolValue: true}},
				"struct": {Kind: &Value_StructValue{StructValue: &Struct{Fields: map[string]*Value{
					"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
					"1":      {Kind: &Value_StringValue{StringValue: "valuea"}},
					"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
					"number": {Kind: &Value_NumberValue{NumberValue: 10}},
					"struct": {Kind: &Value_StructValue{StructValue: &Struct{Fields: map[string]*Value{
						"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
						"1":      {Kind: &Value_StringValue{StringValue: "valuea"}},
						"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
						"number": {Kind: &Value_NumberValue{NumberValue: 10}},
					}}}},
				}}}},
				"list": {Kind: &Value_ListValue{ListValue: &ListValue{Values: []*Value{
					{Kind: &Value_StringValue{StringValue: "dvaluea"}},
				}}}},
			},
		}
		structSort2 = &Struct{
			Fields: map[string]*Value{
				"1": {Kind: &Value_StringValue{StringValue: "valuea"}},
				"struct": {Kind: &Value_StructValue{StructValue: &Struct{Fields: map[string]*Value{
					"1":    {Kind: &Value_StringValue{StringValue: "valuea"}},
					"null": {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
					"struct": {Kind: &Value_StructValue{StructValue: &Struct{Fields: map[string]*Value{
						"number": {Kind: &Value_NumberValue{NumberValue: 10}},
						"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
						"1":      {Kind: &Value_StringValue{StringValue: "valuea"}},
						"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
					}}}},
					"number": {Kind: &Value_NumberValue{NumberValue: 10}},
					"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
				}}}},
				"10000":  {Kind: &Value_NumberValue{NumberValue: 10}},
				"number": {Kind: &Value_NumberValue{NumberValue: 10}},
				"null":   {Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}},
				"bool":   {Kind: &Value_BoolValue{BoolValue: true}},
				"string": {Kind: &Value_StringValue{StringValue: "valuea"}},
				"list": {Kind: &Value_ListValue{ListValue: &ListValue{Values: []*Value{
					{Kind: &Value_StringValue{StringValue: "dvaluea"}},
				}}}},
			},
		}
		structValueSort1 []byte
		structValueSort2 []byte
		err              error
		structUnm1       *Struct
		structUnm2       *Struct
	)
	t.Run("Marshal", func(t *testing.T) {
		structValueSort1, err = codec.MarshalBinaryBare(structSort1)
		require.NoError(t, err)
		structValueSort2, err = codec.MarshalBinaryBare(structSort2)
		require.NoError(t, err)
		require.Equal(t, hash.Dump(structValueSort1), hash.Dump(structValueSort2))
	})
	t.Run("Unmarshal", func(t *testing.T) {
		require.NoError(t, codec.UnmarshalBinaryBare(structValueSort1, &structUnm1))
		require.Equal(t, structSort1, structUnm1)
		require.NoError(t, codec.UnmarshalBinaryBare(structValueSort2, &structUnm2))
		require.Equal(t, structSort2, structUnm2)
		require.Equal(t, structUnm1, structUnm2)
	})
}
