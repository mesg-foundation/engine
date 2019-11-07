package ownershipsdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/database/store"
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

func (s *Backend) db(request cosmostypes.Request) *database.OwnershipDB {
	return database.NewOwnershipDB(store.NewCosmosStore(request.KVStore(s.storeKey)))
}

func (s *Backend) handler(request cosmostypes.Request, msg cosmostypes.Msg) cosmostypes.Result {
	errmsg := fmt.Sprintf("Unrecognized ownership Msg type: %v", msg.Type())
	return cosmostypes.ErrUnknownRequest(errmsg).Result()
}

func (s *Backend) querier(request cosmostypes.Request, path []string, req abci.RequestQuery) (interface{}, error) {
	switch path[0] {
	case "list":
		return s.List(request)
	default:
		return nil, errors.New("unknown ownership query endpoint" + path[0])
	}
}

// CreateServiceOwnership creates a new ownership.
func (s *Backend) CreateServiceOwnership(request cosmostypes.Request, serviceHash hash.Hash, owner cosmostypes.AccAddress) (*ownership.Ownership, error) {
	db := s.db(request)
	// check if owner is authorized to create the ownership
	allOwnshp, err := db.All()
	if err != nil {
		return nil, err
	}
	ownshpSrv := ownershipsOfService(allOwnshp, serviceHash)
	// check if service already have an owner.
	if len(ownshpSrv) > 0 {
		return nil, fmt.Errorf("service %q has already an owner", serviceHash.String())
	}
	ownership := &ownership.Ownership{
		Owner: owner.String(),
		Resource: &ownership.Ownership_ServiceHash{
			ServiceHash: serviceHash,
		},
	}
	ownership.Hash = hash.Dump(ownership)
	return ownership, db.Save(ownership)
}

// List returns all ownerships.
func (s *Backend) List(request cosmostypes.Request) ([]*ownership.Ownership, error) {
	return s.db(request).All()
}

// ownershipsOfService only returns the ownership concerning the specify service.
func ownershipsOfService(allOwnshp []*ownership.Ownership, serviceHash hash.Hash) []*ownership.Ownership {
	ownshpSrv := make([]*ownership.Ownership, 0)
	for _, o := range allOwnshp {
		switch x := o.Resource.(type) {
		case *ownership.Ownership_ServiceHash:
			if x.ServiceHash.Equal(serviceHash) {
				ownshpSrv = append(ownshpSrv, o)
			}
		default:
			continue
		}
	}
	return ownshpSrv
}
