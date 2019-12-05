package resultsdk

import "github.com/mesg-foundation/engine/codec"

func init() {
	codec.RegisterConcrete(msgCreateResult{}, "result/create", nil)
}
