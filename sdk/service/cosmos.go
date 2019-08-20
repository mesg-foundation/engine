package servicesdk

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

type CosmosClient struct {
}

func NewCosmosClient(app *cosmos.App, c container.Container, keeperFactory KeeperFactory) *CosmosClient {
	app.RegisterModule(newAppModule(c, keeperFactory))
	return &CosmosClient{}
}

// Create creates a new service from definition.
func (s *CosmosClient) Create(ctx context.Context, srv *service.Service) (*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// Delete deletes the service by hash.
func (s *CosmosClient) Delete(ctx context.Context, hash hash.Hash) error {
	return fmt.Errorf("not implemented")
}

// Get returns the service that matches given hash.
func (s *CosmosClient) Get(ctx context.Context, hash hash.Hash) (*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// List returns all services.
func (s *CosmosClient) List(ctx context.Context) ([]*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}
