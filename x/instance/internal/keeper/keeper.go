package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/x/instance/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the instance store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper creates a instance keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// FetchOrCreate creates a new instance if needed.
func (k Keeper) FetchOrCreate(ctx sdk.Context, serviceHash hash.Hash, envHash hash.Hash) (*instance.Instance, error) {
	store := ctx.KVStore(k.storeKey)

	inst := instance.New(serviceHash, envHash)

	if !store.Has(inst.Hash) {
		value, err := k.cdc.MarshalBinaryLengthPrefixed(inst)
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
		}
		store.Set(inst.Hash, value)
	}

	return inst, nil
}

// Get returns the instance from the keeper.
func (k Keeper) Get(ctx sdk.Context, instanceHash hash.Hash) (*instance.Instance, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(instanceHash) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "instance %q not found", instanceHash)
	}

	var item *instance.Instance
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(instanceHash), &item); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return item, nil
}

// List returns instances from the keeper.
func (k Keeper) List(ctx sdk.Context, f *api.ListInstanceRequest_Filter) ([]*instance.Instance, error) {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(nil, nil)
	var items []*instance.Instance

	for iter.Valid() {
		var item *instance.Instance
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &item); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		if f == nil || f.ServiceHash.IsZero() || item.ServiceHash.Equal(f.ServiceHash) {
			items = append(items, item)
		}
		iter.Next()
	}
	iter.Close()
	return items, nil
}
