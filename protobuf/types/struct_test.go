package types

import (
	"bytes"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
)

func TestStructMarshal(t *testing.T) {
	cdc := codec.New()
	RegisterCodec(cdc)
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
	t.Run("BinaryLengthPrefixed", func(t *testing.T) {
		var (
			structValueSort1 []byte
			structValueSort2 []byte
			err              error
			structUnm1       *Struct
			structUnm2       *Struct
		)
		t.Run("Marshal", func(t *testing.T) {
			structValueSort1, err = cdc.MarshalBinaryLengthPrefixed(structSort1)
			require.NoError(t, err)
			structValueSort2, err = cdc.MarshalBinaryLengthPrefixed(structSort2)
			require.NoError(t, err)
			require.True(t, bytes.Equal(structValueSort1, structValueSort2))
		})
		t.Run("Unmarshal", func(t *testing.T) {
			require.NoError(t, cdc.UnmarshalBinaryLengthPrefixed(structValueSort1, &structUnm1))
			require.True(t, structSort1.Equal(structUnm1))
			require.NoError(t, cdc.UnmarshalBinaryLengthPrefixed(structValueSort2, &structUnm2))
			require.True(t, structSort2.Equal(structUnm2))
			require.True(t, structUnm1.Equal(structUnm2))
		})
	})
	t.Run("JSON", func(t *testing.T) {
		var (
			structValueSort1 []byte
			structValueSort2 []byte
			err              error
			structUnm1       *Struct
			structUnm2       *Struct
		)
		t.Run("Marshal", func(t *testing.T) {
			structValueSort1, err = cdc.MarshalJSON(structSort1)
			require.NoError(t, err)
			structValueSort2, err = cdc.MarshalJSON(structSort2)
			require.NoError(t, err)
			require.True(t, bytes.Equal(structValueSort1, structValueSort2))
		})
		t.Run("Unmarshal", func(t *testing.T) {
			require.NoError(t, cdc.UnmarshalJSON(structValueSort1, &structUnm1))
			require.True(t, structSort1.Equal(structUnm1))
			require.NoError(t, cdc.UnmarshalJSON(structValueSort2, &structUnm2))
			require.True(t, structSort2.Equal(structUnm2))
			require.True(t, structUnm1.Equal(structUnm2))
		})
	})
}
