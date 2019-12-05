package executionsdk

import "github.com/mesg-foundation/engine/codec"

func init() {
	codec.RegisterConcrete(msgCreateExecution{}, "execution/create", nil)
}
