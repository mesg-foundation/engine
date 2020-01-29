package servicesdk

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	ownershipsdk "github.com/mesg-foundation/engine/sdk/ownership"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/service/validator"
)

// Keeper holds the logic to read and write data.
type Keeper struct {
	storeKey   *cosmostypes.KVStoreKey
	ownerships *ownershipsdk.Keeper
}

// NewKeeper initialize a new keeper.
func NewKeeper(storeKey *cosmostypes.KVStoreKey, ownerships *ownershipsdk.Keeper) *Keeper {
	return &Keeper{
		storeKey:   storeKey,
		ownerships: ownerships,
	}
}

// Create creates a new service.
func (k *Keeper) Create(request cosmostypes.Request, msg *msgCreateService) (*service.Service, error) {
	store := request.KVStore(k.storeKey)
	// create service
	srv := initializeService(msg.Request)

	// check if service already exists.
	if store.Has(srv.Hash) {
		return nil, fmt.Errorf("service %q already exists", srv.Hash)
	}

	// TODO: the following test should be moved in New function
	if srv.Sid == "" {
		// make sure that sid doesn't have the same length with id.
		srv.Sid = "_" + srv.Hash.String()
	}

	if err := validator.ValidateService(srv); err != nil {
		return nil, err
	}

	if _, err := k.ownerships.Create(request, msg.Owner, srv.Hash, ownership.Ownership_Service); err != nil {
		return nil, err
	}

	value, err := codec.MarshalBinaryBare(srv)
	if err != nil {
		return nil, err
	}
	store.Set(srv.Hash, value)
	return srv, nil
}

// Get returns the service that matches given hash.
func (k *Keeper) Get(request cosmostypes.Request, hash hash.Hash) (*service.Service, error) {
	var sv *service.Service
	store := request.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("service %q not found", hash)
	}
	value := store.Get(hash)
	return sv, codec.UnmarshalBinaryBare(value, &sv)
}

// Exists returns true if a specific set of data exists in the database, false otherwise
func (k *Keeper) Exists(request cosmostypes.Request, hash hash.Hash) (bool, error) {
	return request.KVStore(k.storeKey).Has(hash), nil
}

// Hash returns the hash of a service request.
func (k *Keeper) Hash(serviceRequest *api.CreateServiceRequest) hash.Hash {
	return initializeService(serviceRequest).Hash
}

// List returns all services.
func (k *Keeper) List(request cosmostypes.Request) ([]*service.Service, error) {
	var (
		services []*service.Service
		iter     = request.KVStore(k.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var sv *service.Service
		if err := codec.UnmarshalBinaryBare(iter.Value(), &sv); err != nil {
			return nil, err
		}
		services = append(services, sv)
		iter.Next()
	}
	iter.Close()
	return services, nil
}

func initializeService(req *api.CreateServiceRequest) *service.Service {
	// create service
	srv := &service.Service{
		Sid:           req.Sid,
		Name:          req.Name,
		Description:   req.Description,
		Configuration: req.Configuration,
		Tasks:         req.Tasks,
		Events:        req.Events,
		Dependencies:  req.Dependencies,
		Repository:    req.Repository,
		Source:        req.Source,
	}

	// calculate and apply hash to service.
	srv.Hash = hash.Dump(srv)
	return srv
}
