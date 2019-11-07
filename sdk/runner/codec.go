package runnersdk

import "github.com/mesg-foundation/engine/codec"

func init() {
	codec.RegisterConcrete(msgCreateRunner{}, "runner/create", nil)
	codec.RegisterConcrete(msgDeleteRunner{}, "runner/delete", nil)
}
