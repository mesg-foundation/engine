package servicesdk

import (
	"github.com/mesg-foundation/engine/codec"
)

func init() {
	codec.RegisterConcrete(msgCreateService{}, "service/create", nil)
}
