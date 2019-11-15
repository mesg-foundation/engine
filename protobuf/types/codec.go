package types

import (
	"github.com/mesg-foundation/engine/codec"
)

func init() {
	codec.RegisterInterface((*isValue_Kind)(nil), nil)
	codec.RegisterConcrete(&Value_NullValue{}, "mesg.types.Value_NullValue", nil)
	codec.RegisterConcrete(&Value_NumberValue{}, "mesg.types.Value_NumberValue", nil)
	codec.RegisterConcrete(&Value_StringValue{}, "mesg.types.Value_StringValue", nil)
	codec.RegisterConcrete(&Value_BoolValue{}, "mesg.types.Value_BoolValue", nil)
	codec.RegisterConcrete(&Value_ListValue{}, "mesg.types.Value_ListValue", nil)
	codec.RegisterConcrete(&Value_StructValue{}, "mesg.types.Value_StructValue", nil)
}
