package servicesdk

import (
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/service"
)

// Service is the interface of this sdk
type Service interface {
	Create(req *api.CreateServiceRequest, accountName string, accountPassword string) (*service.Service, error)
	Get(hash hash.Hash) (*service.Service, error)
	List() ([]*service.Service, error)
	Exists(hash hash.Hash) (bool, error)
	Hash(req *api.CreateServiceRequest) (hash.Hash, error)
}
