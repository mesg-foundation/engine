package servicesdk

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/store"
	"github.com/mesg-foundation/engine/tendermint"
	abci "github.com/tendermint/tendermint/abci/types"
)

type sdk struct {
	container container.Container

	name string
	cdc  *codec.Codec
}

func New(app *cosmos.App, c container.Container) Service {
	sdk := &sdk{
		container: c,
		name:      "service",
		cdc:       app.Cdc(),
	}
	appModuleBasic := tendermint.NewAppModuleBasic("service")
	appModule := tendermint.NewAppModule(appModuleBasic, sdk.handler, sdk.querier)
	app.RegisterModule(appModule)
	return sdk
}

// Create creates a new service from definition.
func (s *sdk) Create(srv *service.Service) (*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// Delete deletes the service by hash.
func (s *sdk) Delete(hash hash.Hash) error {
	return fmt.Errorf("not implemented")
}

// Get returns the service that matches given hash.
func (s *sdk) Get(hash hash.Hash) (*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// List returns all services.
func (s *sdk) List() ([]*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// ------------------------------------
//   Handlers
// ------------------------------------

func (s *sdk) keeperFromRequest(request cosmostypes.Request) *database.ServiceKeeper {
	return database.NewServiceKeeper(store.NewCosmosStore(request.KVStore(cosmostypes.NewKVStoreKey(s.name))))
}

func (s *sdk) handler(request cosmostypes.Request, msg cosmostypes.Msg) cosmostypes.Result {
	panic("to implement")
	return cosmostypes.Result{}
}

func (s *sdk) querier(request cosmostypes.Request, path []string, req abci.RequestQuery) (res []byte, err cosmostypes.Error) {
	panic("to implement")
	keeper := s.keeperFromRequest(request)
	switch path[0] {
	case "get":
		hash, err := hash.Decode(path[1])
		if err != nil {
			return nil, cosmostypes.ErrInternal(err.Error())
		}
		service, err := keeper.Get(hash)
		if err != nil {
			return nil, cosmostypes.ErrInternal(err.Error())
		}
		res, err := s.cdc.MarshalJSON(service)
		if err != nil {
			return nil, cosmostypes.ErrInternal(err.Error())
		}
		return res, nil
	case "list":
		services, err := keeper.All()
		if err != nil {
			return nil, cosmostypes.ErrInternal(err.Error())
		}
		res, err := s.cdc.MarshalJSON(services)
		if err != nil {
			return nil, cosmostypes.ErrInternal(err.Error())
		}
		return res, nil
	default:
		return nil, cosmostypes.ErrUnknownRequest("unknown service query endpoint" + path[0])
	}
}

// Create creates a new service from definition.
func (s *sdk) create(request cosmostypes.Request, srv *service.Service) (*service.Service, error) {
	return create(s.container, s.keeperFromRequest(request), srv)
}

// Delete deletes the service by hash.
func (s *sdk) delete(request cosmostypes.Request, hash hash.Hash) error {
	return s.keeperFromRequest(request).Delete(hash)
}
