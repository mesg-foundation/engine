package servicesdk

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

type Service interface {
	Create(ctx context.Context, srv *service.Service) (*service.Service, error)
	Delete(ctx context.Context, hash hash.Hash) error
	Get(ctx context.Context, hash hash.Hash) (*service.Service, error)
	List(ctx context.Context) ([]*service.Service, error)
}

// AlreadyExistsError is an not found error.
type AlreadyExistsError struct {
	Hash hash.Hash
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("service %q already exists", e.Hash.String())
}

type KeeperFactory func(context.Context) (*database.ServiceKeeper, error)
