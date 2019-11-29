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

// Backend is the ownership backend.
type Backend struct {
	storeKey *cosmostypes.KVStoreKey
}

// NewBackend returns the backend of the ownership sdk.
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

	hashes := s.findOwnerships(store, resourceHash)
	if len(hashes) > 0 {
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
	return ownership, nil
}

// Delete deletes a ownership.
func (s *Backend) Delete(req cosmostypes.Request, owner cosmostypes.AccAddress, resourceHash hash.Hash) error {
	store := req.KVStore(s.storeKey)

	hashes := s.findOwnerships(store, resourceHash)
	if len(hashes) == 0 {
		return fmt.Errorf("resource %q do not have any ownership", resourceHash)
	}

	for _, hash := range hashes {
		store.Delete(hash)
	}
	return nil
}

// List returns all ownerships.
func (s *Backend) List(req cosmostypes.Request) ([]*ownership.Ownership, error) {
	var (
		ownerships []*ownership.Ownership
		iter       = req.KVStore(s.storeKey).Iterator(nil, nil)
	)
	defer iter.Close()

	for iter.Valid() {
		var o *ownership.Ownership
		if err := codec.UnmarshalBinaryBare(iter.Value(), &o); err != nil {
			return nil, err
		}
		ownerships = append(ownerships, o)
		iter.Next()
	}
	return ownerships, nil
}

// hasOwner checks if given resource has owner. Returns all ownership hash and true if has one
// nil and false otherwise.
func (s *Backend) findOwnerships(store cosmostypes.KVStore, resourceHash hash.Hash) []hash.Hash {
	var ownerships []hash.Hash
	iter := store.Iterator(nil, nil)

	for iter.Valid() {
		var o *ownership.Ownership
		if err := codec.UnmarshalBinaryBare(iter.Value(), &o); err == nil {
			if o.ResourceHash.Equal(resourceHash) {
				ownerships = append(ownerships, o.Hash)
			}
		}
		iter.Next()
	}

	iter.Close()
	return ownerships
}
