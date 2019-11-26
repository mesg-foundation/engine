package ownershipsdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	abci "github.com/tendermint/tendermint/abci/types"
)

const backendName = "ownership"

const storeHashKey = "ownership_hash"

// Backend is the ownership backend.
type Backend struct {
	storeKey     *cosmostypes.KVStoreKey
	storeHashKey *cosmostypes.KVStoreKey
}

// NewBackend returns the backend of the ownership sdk.
func NewBackend(appFactory *cosmos.AppFactory) *Backend {
	backend := &Backend{
		storeKey:     cosmostypes.NewKVStoreKey(backendName),
		storeHashKey: cosmostypes.NewKVStoreKey(storeHashKey),
	}
	appBackendBasic := cosmos.NewAppModuleBasic(backendName)
	appBackend := cosmos.NewAppModule(appBackendBasic, backend.handler, backend.querier)
	appFactory.RegisterModule(appBackend)
	appFactory.RegisterStoreKey(backend.storeKey)
	appFactory.RegisterStoreKey(backend.storeHashKey)
	return backend
}

func (s *Backend) handler(request cosmostypes.Request, msg cosmostypes.Msg) (hash.Hash, error) {
	errmsg := fmt.Sprintf("Unrecognized ownership Msg type: %v", msg.Type())
	return nil, cosmostypes.ErrUnknownRequest(errmsg)
}

func (s *Backend) querier(request cosmostypes.Request, path []string, req abci.RequestQuery) (interface{}, error) {
	switch path[0] {
	case "list":
		return s.List(request)
	default:
		return nil, errors.New("unknown ownership query endpoint" + path[0])
	}
}

// Create creates a new ownership.
func (s *Backend) Create(req cosmostypes.Request, owner cosmostypes.AccAddress, resourceHash hash.Hash, resource ownership.Ownership_Resource) (*ownership.Ownership, error) {
	store := req.KVStore(s.storeKey)
	storeHash := req.KVStore(s.storeHashKey)

	if storeHash.Has(resourceHash) {
		return nil, fmt.Errorf("resource %s:%q already has an owner", resource, resourceHash)
	}
	ownership := &ownership.Ownership{
		Owner:        owner.String(),
		Resource:     resource,
		ResourceHash: resourceHash,
	}
	ownership.Hash = hash.Dump(ownership)

	value, err := codec.MarshalBinaryBare(ownership)
	if err != nil {
		return nil, err
	}
	store.Set(ownership.Hash, value)
	storeHash.Set(resourceHash, ownership.Hash)
	return ownership, nil
}

// Delete deletes a ownership.
func (s *Backend) Delete(req cosmostypes.Request, owner cosmostypes.AccAddress, resourceHash hash.Hash) error {
	store := req.KVStore(s.storeKey)
	storeHash := req.KVStore(s.storeHashKey)

	if !storeHash.Has(resourceHash) {
		return fmt.Errorf("resource %q do not exist", resourceHash)
	}

	ownershipHash := storeHash.Get(resourceHash)
	store.Delete(ownershipHash)
	storeHash.Delete(resourceHash)
	return nil
}

// CreateServiceOwnership creates a new ownership.
func (s *Backend) CreateServiceOwnership(request cosmostypes.Request, serviceHash hash.Hash, owner cosmostypes.AccAddress) (*ownership.Ownership, error) {
	store := request.KVStore(s.storeKey)
	// check if owner is authorized to create the ownership
	allOwnshp, err := s.List(request)
	if err != nil {
		return nil, err
	}
	ownshpSrv := ownershipsOfService(allOwnshp, serviceHash)
	// check if service already have an owner.
	if len(ownshpSrv) > 0 {
		return nil, fmt.Errorf("service %q has already an owner", serviceHash)
	}
	ownership := &ownership.Ownership{
		Owner:        owner.String(),
		Resource:     ownership.Ownership_Service,
		ResourceHash: serviceHash,
	}
	ownership.Hash = hash.Dump(ownership)
	value, err := codec.MarshalBinaryBare(ownership)
	if err != nil {
		return nil, err
	}
	store.Set(ownership.Hash, value)
	return ownership, nil
}

// List returns all ownerships.
func (s *Backend) List(request cosmostypes.Request) ([]*ownership.Ownership, error) {
	var (
		ownerships []*ownership.Ownership
		iter       = request.KVStore(s.storeKey).Iterator(nil, nil)
	)

	for iter.Valid() {
		var o *ownership.Ownership
		if err := codec.UnmarshalBinaryBare(iter.Value(), &o); err != nil {
			return nil, err
		}
		ownerships = append(ownerships, o)
		iter.Next()
	}
	iter.Close()
	return ownerships, nil
}

// ownershipsOfService only returns the ownership concerning the specify service.
func ownershipsOfService(allOwnshp []*ownership.Ownership, serviceHash hash.Hash) []*ownership.Ownership {
	ownshpSrv := make([]*ownership.Ownership, 0)
	for _, o := range allOwnshp {
		switch o.Resource {
		case ownership.Ownership_Service:
			if o.ResourceHash.Equal(serviceHash) {
				ownshpSrv = append(ownshpSrv, o)
			}
		default:
			continue
		}
	}
	return ownshpSrv
}
