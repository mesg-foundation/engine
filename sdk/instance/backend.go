package instancesdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	abci "github.com/tendermint/tendermint/abci/types"
)

const backendName = "instance"

// Backend is the instance backend.
type Backend struct {
	storeKey *cosmostypes.KVStoreKey
}

// NewBackend returns the backend of the instance sdk.
func NewBackend(appFactory *cosmos.AppFactory) *Backend {
	backend := &Backend{
		storeKey: cosmostypes.NewKVStoreKey(backendName),
	}
	appBackendBasic := cosmos.NewAppModuleBasic(backendName)
	appBackend := cosmos.NewAppModule(appBackendBasic, backend.handler, backend.querier)
	appFactory.RegisterModule(appBackend)
	appFactory.RegisterStoreKey(backend.storeKey)

	return backend
}

func (s *Backend) db(request cosmostypes.Request) *database.InstanceDB {
	return database.NewInstanceDB(store.NewCosmosStore(request.KVStore(s.storeKey)))
}

func (s *Backend) handler(request cosmostypes.Request, msg cosmostypes.Msg) cosmostypes.Result {
	errmsg := fmt.Sprintf("Unrecognized instance Msg type: %v", msg.Type())
	return cosmostypes.ErrUnknownRequest(errmsg).Result()
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
	case "exists":
		hash, err := hash.Decode(path[1])
		if err != nil {
			return nil, err
		}
		return s.Exists(request, hash)

	default:
		return nil, errors.New("unknown instance query endpoint" + path[0])
	}
}

// FetchOrCreate creates a new instance if needed.
func (s *Backend) FetchOrCreate(request cosmostypes.Request, serviceHash hash.Hash, envHash hash.Hash) (*instance.Instance, error) {
	db := s.db(request)

	// create instance hash
	inst := &instance.Instance{
		ServiceHash: serviceHash,
		EnvHash:     envHash,
	}
	inst.Hash = hash.Dump(inst)

	// create instance if needed
	if exists, err := db.Exist(inst.Hash); err != nil {
		return nil, err
	} else if !exists {
		if err := db.Save(inst); err != nil {
			return nil, err
		}
	}

	return inst, nil
}

// Get returns the instance that matches given hash.
func (s *Backend) Get(request cosmostypes.Request, hash hash.Hash) (*instance.Instance, error) {
	return s.db(request).Get(hash)
}

// Exists returns true if a specific set of data exists in the database, false otherwise
func (s *Backend) Exists(request cosmostypes.Request, hash hash.Hash) (bool, error) {
	return s.db(request).Exist(hash)
}

// List returns all instances.
func (s *Backend) List(request cosmostypes.Request) ([]*instance.Instance, error) {
	return s.db(request).All()
}
