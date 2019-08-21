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
	appModuleBasic := cosmos.NewAppModuleBasic("service")
	appModule := cosmos.NewAppModule(appModuleBasic, app.Cdc(), sdk.handler, sdk.querier)
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
	// keeper := s.keeperFromRequest(request)
	switch msg := msg.(type) {
	// case MsgCreateService:
	// 	_, err := create(s.container, keeper, msg.Service)
	// 	if err != nil {
	// 		return cosmostypes.ErrInternal(err.Error()).Result()
	// 	}
	// 	return cosmostypes.Result{}
	// case MsgRemoveService:
	// 	err := keeper.Delete(msg.Hash)
	// 	if err != nil {
	// 		return cosmostypes.ErrInternal(err.Error()).Result()
	// 	}
	// 	return cosmostypes.Result{}
	default:
		errmsg := fmt.Sprintf("Unrecognized service Msg type: %v", msg.Type())
		return cosmostypes.ErrUnknownRequest(errmsg).Result()
	}
}

func (s *sdk) querier(request cosmostypes.Request, path []string, req abci.RequestQuery) (interface{}, error) {
	panic("to implement")
	keeper := s.keeperFromRequest(request)
	switch path[0] {
	case "get":
		hash, err := hash.Decode(path[1])
		if err != nil {
			return nil, err
		}
		return keeper.Get(hash)
	case "list":
		return keeper.All()
	default:
		return nil, fmt.Errorf("unknown service query endpoint %q", path[0])
	}
}
