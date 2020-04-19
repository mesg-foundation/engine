package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/hash"
	ownershippb "github.com/mesg-foundation/engine/ownership"
	servicepb "github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/x/service/internal/types"
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
func (k Keeper) Create(ctx sdk.Context, msg *types.MsgCreate) (*servicepb.Service, error) {
	store := ctx.KVStore(k.storeKey)
	// create service
	srv, err := servicepb.New(
		msg.Sid,
		msg.Name,
		msg.Description,
		msg.Configuration,
		msg.Tasks,
		msg.Events,
		msg.Dependencies,
		msg.Repository,
		msg.Source,
	)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// check if service already exists.
	if store.Has(srv.Hash) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "service %q already exists", srv.Hash)
	}

	if _, err := k.ownershipKeeper.Set(ctx, msg.Owner, srv.Hash, ownershippb.Ownership_Service, srv.Address); err != nil {
		return nil, err
	}

	value, err := k.cdc.MarshalBinaryLengthPrefixed(srv)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	store.Set(srv.Hash, value)

	// emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventType,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeActionCreated),
			sdk.NewAttribute(types.AttributeKeyHash, srv.Hash.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, srv.Address.String()),
		),
	)

	return srv, nil
}

// Get returns the service that matches given hash.
func (k Keeper) Get(ctx sdk.Context, hash hash.Hash) (*servicepb.Service, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "service %q not found", hash)
	}
	var sv *servicepb.Service
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(hash), &sv); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return sv, nil
}

// Exists returns true if a specific set of data exists in the database, false otherwise
func (k Keeper) Exists(ctx sdk.Context, hash hash.Hash) (bool, error) {
	return ctx.KVStore(k.storeKey).Has(hash), nil
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
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		services = append(services, sv)
		iter.Next()
	}
	iter.Close()
	return services, nil
}
