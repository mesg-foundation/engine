package servicesdk

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Module is the service module.
type Module struct {
	container container.Container
	name      string
	cdc       *codec.Codec
	storeKey  *cosmostypes.KVStoreKey
}

// NewModule returns the module of the service sdk.
func NewModule(app *cosmos.App, c container.Container) *Module {
	name := "service"
	module := &Module{
		container: c,
		name:      name,
		cdc:       app.Cdc(),
		storeKey:  cosmostypes.NewKVStoreKey(name),
	}
	appModuleBasic := cosmos.NewAppModuleBasic(name)
	appModule := cosmos.NewAppModule(appModuleBasic, module.cdc, module.handler, module.querier)
	app.RegisterModule(appModule)
	app.RegisterStoreKey(module.storeKey)
	return module
}

func (s *Module) db(request cosmostypes.Request) *database.ServiceDB {
	return database.NewServiceDB(store.NewCosmosStore(request.KVStore(s.storeKey)))
}

func (s *Module) handler(request cosmostypes.Request, msg cosmostypes.Msg) cosmostypes.Result {
	panic("to implement")
	// switch msg := msg.(type) {
	// case MsgCreateService:
	// 	_, err := s.Create(request, msg.Service)
	// 	if err != nil {
	// 		return cosmostypes.ErrInternal(err.Error()).Result()
	// 	}
	// 	return cosmostypes.Result{}
	// case MsgRemoveService:
	// 	err := s.Delete(request, msg.Hash)
	// 	if err != nil {
	// 		return cosmostypes.ErrInternal(err.Error()).Result()
	// 	}
	// 	return cosmostypes.Result{}
	// default:
	// 	errmsg := fmt.Sprintf("Unrecognized service Msg type: %v", msg.Type())
	// 	return cosmostypes.ErrUnknownRequest(errmsg).Result()
	// }
}

func (s *Module) querier(request cosmostypes.Request, path []string, req abci.RequestQuery) (interface{}, error) {
	panic("to implement")
	// switch path[0] {
	// case "get":
	// 	hash, err := hash.Decode(path[1])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return s.Get(request, hash)
	// case "list":
	// 	return s.All(request)
	// default:
	// 	return nil, fmt.Errorf("unknown service query endpoint %q", path[0])
	// }
}

// Create creates a new service from definition.
func (s *Module) Create(request cosmostypes.Request, srv *service.Service) (*service.Service, error) {
	return create(s.container, s.db(request), srv)
}

// Delete deletes the service by hash.
func (s *Module) Delete(request cosmostypes.Request, hash hash.Hash) error {
	return s.db(request).Delete(hash)
}

// Get returns the service that matches given hash.
func (s *Module) Get(request cosmostypes.Request, hash hash.Hash) (*service.Service, error) {
	return s.db(request).Get(hash)
}

// List returns all services.
func (s *Module) List(request cosmostypes.Request) ([]*service.Service, error) {
	return s.db(request).All()
}
