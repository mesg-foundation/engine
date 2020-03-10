package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec register struct type in cdc.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*isValue_Kind)(nil), nil)
	cdc.RegisterConcrete(&Value_NullValue{}, "mesg.types.Value_NullValue", nil)
	cdc.RegisterConcrete(&Value_NumberValue{}, "mesg.types.Value_NumberValue", nil)
	cdc.RegisterConcrete(&Value_StringValue{}, "mesg.types.Value_StringValue", nil)
	cdc.RegisterConcrete(&Value_BoolValue{}, "mesg.types.Value_BoolValue", nil)
	cdc.RegisterConcrete(&Value_ListValue{}, "mesg.types.Value_ListValue", nil)
	cdc.RegisterConcrete(&Value_StructValue{}, "mesg.types.Value_StructValue", nil)
}
