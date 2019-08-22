package servicesdk

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	abci "github.com/tendermint/tendermint/abci/types"
)

type sdk struct {
	container container.Container

	name     string
	cdc      *codec.Codec
	storeKey *cosmostypes.KVStoreKey
}

func New(app *cosmos.App, c container.Container) Service {
	sdk := &sdk{
		container: c,
		name:      "service",
		cdc:       app.Cdc(),
		storeKey:  cosmostypes.NewKVStoreKey("service"),
	}
	appModuleBasic := cosmos.NewAppModuleBasic("service")
	appModule := cosmos.NewAppModule(appModuleBasic, app.Cdc(), sdk.handler, sdk.querier)
	app.RegisterModule(appModule)
	app.RegisterStoreKey(sdk.storeKey)
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

// ------------------------------------
//   Handlers
// ------------------------------------

// func (s *sdk) dbFromRequest(request cosmostypes.Request) *database.ServiceDB {
// 	return database.NewServiceDB(store.NewCosmosStore(request.KVStore(s.storeKey)))
// }

func (s *sdk) handler(request cosmostypes.Request, msg cosmostypes.Msg) cosmostypes.Result {
	panic("to implement")
	// db := s.dbFromRequest(request)
	// switch msg := msg.(type) {
	// case MsgCreateService:
	// 	_, err := create(s.container, db, msg.Service)
	// 	if err != nil {
	// 		return cosmostypes.ErrInternal(err.Error()).Result()
	// 	}
	// 	return cosmostypes.Result{}
	// case MsgRemoveService:
	// 	err := db.Delete(msg.Hash)
	// 	if err != nil {
	// 		return cosmostypes.ErrInternal(err.Error()).Result()
	// 	}
	// 	return cosmostypes.Result{}
	// default:
	// 	errmsg := fmt.Sprintf("Unrecognized service Msg type: %v", msg.Type())
	// 	return cosmostypes.ErrUnknownRequest(errmsg).Result()
	// }
}

func (s *sdk) querier(request cosmostypes.Request, path []string, req abci.RequestQuery) (interface{}, error) {
	panic("to implement")
	// db := s.dbFromRequest(request)
	// switch path[0] {
	// case "get":
	// 	hash, err := hash.Decode(path[1])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return db.Get(hash)
	// case "list":
	// 	return db.All()
	// default:
	// 	return nil, fmt.Errorf("unknown service query endpoint %q", path[0])
	// }
}
