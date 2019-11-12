package servicesdk

import (
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/service"
)

func init() {
	codec.RegisterConcrete(msgCreateService{}, "service/create", nil)
	codec.RegisterConcrete(&service.Service{}, "mesg.types.Service", nil)
}
