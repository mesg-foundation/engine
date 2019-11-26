package processsdk

import (
	"github.com/mesg-foundation/engine/codec"
)

func init() {
	codec.RegisterConcrete(msgCreateProcess{}, "process/create", nil)
	codec.RegisterConcrete(msgDeleteProcess{}, "process/delete", nil)
}
