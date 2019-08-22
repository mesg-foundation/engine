package servicesdk

import (
	"fmt"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

type sdk struct {
	app *cosmos.App
}

// NewSDK returns the service sdk.
func NewSDK(app *cosmos.App) Service {
	sdk := &sdk{
		app: app,
	}
	return sdk
}

// Create creates a new service from definition.
func (s *sdk) Create(srv *service.Service) (*service.Service, error) {
	return nil, fmt.Errorf("create not implemented")
}

// Delete deletes the service by hash.
func (s *sdk) Delete(hash hash.Hash) error {
	return fmt.Errorf("delete not implemented")
}

// Get returns the service that matches given hash.
func (s *sdk) Get(hash hash.Hash) (*service.Service, error) {
	return nil, fmt.Errorf("get not implemented")
}

// List returns all services.
func (s *sdk) List() ([]*service.Service, error) {
	return nil, fmt.Errorf("list not implemented")
}
