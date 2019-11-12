package instancesdk

import (
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/protobuf/api"
)

func init() {
	codec.Codec.RegisterConcrete(&api.ListInstanceRequest_Filter{}, "mesg.types.ListInstanceRequest_Filter ", nil)
}
