package instancesdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/api"
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

func (s *Backend) handler(request cosmostypes.Request, msg cosmostypes.Msg) (hash.Hash, cosmostypes.Error) {
	errmsg := fmt.Sprintf("Unrecognized instance Msg type: %v", msg.Type())
	return nil, cosmostypes.ErrUnknownRequest(errmsg)
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
		var f api.ListInstanceRequest_Filter
		if err := codec.UnmarshalBinaryBare(req.Data, &f); err != nil {
			return nil, err
		}

		return s.List(request, &f)
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
	inst := &instance.Instance{
		ServiceHash: serviceHash,
		EnvHash:     envHash,
	}
	inst.Hash = hash.Dump(inst)

	if store := request.KVStore(s.storeKey); !store.Has(inst.Hash) {
		value, err := codec.MarshalBinaryBare(inst)
		if err != nil {
			return nil, err
		}
		store.Set(inst.Hash, value)
	}

	return inst, nil
}

// Get returns the instance that matches given hash.
func (s *Backend) Get(request cosmostypes.Request, hash hash.Hash) (*instance.Instance, error) {
	var i *instance.Instance
	store := request.KVStore(s.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("instance %q not found", hash)
	}
	value := store.Get(hash)
	return i, codec.UnmarshalBinaryBare(value, &i)
}

// Exists returns true if a specific set of data exists in the database, false otherwise
func (s *Backend) Exists(request cosmostypes.Request, hash hash.Hash) (bool, error) {
	return request.KVStore(s.storeKey).Has(hash), nil
}

// List returns all instances.
func (s *Backend) List(request cosmostypes.Request, f *api.ListInstanceRequest_Filter) ([]*instance.Instance, error) {
	var (
		instances []*instance.Instance
		iter      = request.KVStore(s.storeKey).Iterator(nil, nil)
	)

	// filter results
	for iter.Valid() {
		var i *instance.Instance
		if err := codec.UnmarshalBinaryBare(iter.Value(), &i); err != nil {
			return nil, err
		}
		if f == nil || f.ServiceHash.IsZero() || i.ServiceHash.Equal(f.ServiceHash) {
			instances = append(instances, i)
		}
		iter.Next()
	}
	iter.Close()
	return instances, nil
}
