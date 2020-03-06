package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos/address"
	ownershippb "github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/protobuf/api"
	servicepb "github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/service/validator"
	"github.com/mesg-foundation/engine/x/service/internal/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the service store
type Keeper struct {
	storeKey        sdk.StoreKey
	cdc             *codec.Codec
	ownershipKeeper types.OwnershipKeeper
}

// NewKeeper creates a service keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ownershipKeeper types.OwnershipKeeper) Keeper {
	keeper := Keeper{
		storeKey:        key,
		cdc:             cdc,
		ownershipKeeper: ownershipKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Create creates a new service.
func (k Keeper) Create(ctx sdk.Context, msg *types.MsgCreateService) (*servicepb.Service, error) {
	store := ctx.KVStore(k.storeKey)
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

	if _, err := k.ownershipKeeper.Set(ctx, msg.Owner, srv.Hash, ownershippb.Ownership_Service); err != nil {
		return nil, err
	}

	value, err := k.cdc.MarshalBinaryLengthPrefixed(srv)
	if err != nil {
		return nil, err
	}
	store.Set(srv.Hash, value)
	return srv, nil
}

// Get returns the service that matches given hash.
func (k Keeper) Get(ctx sdk.Context, hash address.ServAddress) (*servicepb.Service, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("service %q not found", hash)
	}
	var sv *servicepb.Service
	return sv, k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(hash), &sv)
}

// Exists returns true if a specific set of data exists in the database, false otherwise
func (k Keeper) Exists(ctx sdk.Context, hash address.ServAddress) (bool, error) {
	return ctx.KVStore(k.storeKey).Has(hash), nil
}

// Hash returns the hash of a service request.
func (k Keeper) Hash(_ sdk.Context, serviceRequest *api.CreateServiceRequest) address.ServAddress {
	return initializeService(serviceRequest).Hash
}

// List returns all services.
func (k Keeper) List(ctx sdk.Context) ([]*servicepb.Service, error) {
	var (
		services []*servicepb.Service
		iter     = ctx.KVStore(k.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var sv *servicepb.Service
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &sv); err != nil {
			return nil, err
		}
		services = append(services, sv)
		iter.Next()
	}
	iter.Close()
	return services, nil
}

func initializeService(req *api.CreateServiceRequest) *servicepb.Service {
	// create service
	srv := &servicepb.Service{
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
	srv.Hash = address.ServAddress(crypto.AddressHash([]byte(srv.HashSerialize())))
	return srv
}
