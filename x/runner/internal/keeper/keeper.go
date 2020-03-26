package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/hash"
	ownershippb "github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/x/runner/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the runner store
type Keeper struct {
	storeKey        sdk.StoreKey
	cdc             *codec.Codec
	instanceKeeper  types.InstanceKeeper
	ownershipKeeper types.OwnershipKeeper
}

// NewKeeper creates a runner keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, instanceKeeper types.InstanceKeeper, ownershipKeeper types.OwnershipKeeper) Keeper {
	keeper := Keeper{
		storeKey:        key,
		cdc:             cdc,
		instanceKeeper:  instanceKeeper,
		ownershipKeeper: ownershipKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Create creates a new runner.
func (k Keeper) Create(ctx sdk.Context, msg *types.MsgCreate) (*runner.Runner, error) {
	store := ctx.KVStore(k.storeKey)
	inst, err := k.instanceKeeper.FetchOrCreate(ctx, msg.ServiceHash, msg.EnvHash)
	if err != nil {
		return nil, err
	}

	r, err := runner.New(msg.Owner.String(), inst.Hash)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if store.Has(r.Hash) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "runner %q already exists", r.Hash)
	}

	value, err := k.cdc.MarshalBinaryLengthPrefixed(r)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	if _, err := k.ownershipKeeper.Set(ctx, msg.Owner, r.Hash, ownershippb.Ownership_Runner, r.Address); err != nil {
		return nil, err
	}

	store.Set(r.Hash, value)
	return r, nil
}

// Delete deletes a runner.
func (k Keeper) Delete(ctx sdk.Context, msg *types.MsgDelete) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(msg.Hash) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "runner %q not found", msg.Hash)
	}

	value := store.Get(msg.Hash)
	var r *runner.Runner
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(value, &r); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if r.Owner != msg.Owner.String() {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only the runner owner can remove itself")
	}

	if err := k.ownershipKeeper.Delete(ctx, msg.Owner, r.Hash); err != nil {
		return err
	}
	store.Delete(msg.Hash)
	return nil
}

// Get returns the runner that matches given hash.
func (k Keeper) Get(ctx sdk.Context, hash hash.Hash) (*runner.Runner, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "runner %q not found", hash)
	}
	value := store.Get(hash)
	var r *runner.Runner
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(value, &r); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return r, nil
}

// List returns all runners.
func (k Keeper) List(ctx sdk.Context) ([]*runner.Runner, error) {
	var (
		runners []*runner.Runner
		iter    = ctx.KVStore(k.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var r *runner.Runner
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &r); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		runners = append(runners, r)
		iter.Next()
	}
	iter.Close()
	return runners, nil
}
