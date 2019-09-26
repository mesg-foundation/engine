package servicesdk

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	ownershipsdk "github.com/mesg-foundation/engine/sdk/ownership"
	abci "github.com/tendermint/tendermint/abci/types"
)

const backendName = "service"

// Backend is the service backend.
type Backend struct {
	container container.Container
	cdc       *codec.Codec
	storeKey  *cosmostypes.KVStoreKey
	ownerships *ownershipsdk.Backend
}

// NewBackend returns the backend of the service sdk.
func NewBackend(appFactory *cosmos.AppFactory, c container.Container, ownerships *ownershipsdk.Backend) *Backend {
	backend := &Backend{
		container: c,
		cdc:       appFactory.Cdc(),
		storeKey:  cosmostypes.NewKVStoreKey(backendName),
		ownerships: ownerships,
	}
	appBackendBasic := cosmos.NewAppModuleBasic(backendName)
	appBackend := cosmos.NewAppModule(appBackendBasic, backend.cdc, backend.handler, backend.querier)
	appFactory.RegisterModule(appBackend)
	appFactory.RegisterStoreKey(backend.storeKey)

	backend.cdc.RegisterConcrete(msgCreateService{}, "service/create", nil)
	backend.cdc.RegisterConcrete(msgDeleteService{}, "service/delete", nil)

	return backend
}

func (s *Backend) db(request cosmostypes.Request) *database.ServiceDB {
	return database.NewServiceDB(store.NewCosmosStore(request.KVStore(s.storeKey)), s.cdc)
}

func (s *Backend) handler(request cosmostypes.Request, msg cosmostypes.Msg) cosmostypes.Result {
	switch msg := msg.(type) {
	case msgCreateService:
		srv, err := s.Create(request, &msg)
		if err != nil {
			return cosmostypes.ErrInternal(err.Error()).Result()
		}
		return cosmostypes.Result{
			Data: srv.Hash,
		}
	case msgDeleteService:
		err := s.Delete(request, msg.Hash)
		if err != nil {
			return cosmostypes.ErrInternal(err.Error()).Result()
		}
		return cosmostypes.Result{}
	default:
		errmsg := fmt.Sprintf("Unrecognized service Msg type: %v", msg.Type())
		return cosmostypes.ErrUnknownRequest(errmsg).Result()
	}
}

func (s *Backend) querier(request cosmostypes.Request, path []string, req abci.RequestQuery) (interface{}, error) {
	switch path[0] {
	case "get":
		hash, err := hash.Decode(path[1])
		if err != nil {
			return nil, err
		}
		return s.Get(request, hash)
	case "list":
		return s.List(request)
	default:
		return nil, errors.New("unknown service query endpoint" + path[0])
	}
}

// Create creates a new service from definition.
func (s *Backend) Create(request cosmostypes.Request, msg *msgCreateService) (*service.Service, error) {
	return create(s.container, s.db(request), msg.Request, msg.Owner, s.ownerships, request)
}

// Delete deletes the service by hash.
func (s *Backend) Delete(request cosmostypes.Request, hash hash.Hash) error {
	return s.db(request).Delete(hash)
}

// Get returns the service that matches given hash.
func (s *Backend) Get(request cosmostypes.Request, hash hash.Hash) (*service.Service, error) {
	return s.db(request).Get(hash)
}

// List returns all services.
func (s *Backend) List(request cosmostypes.Request) ([]*service.Service, error) {
	return s.db(request).All()
}
