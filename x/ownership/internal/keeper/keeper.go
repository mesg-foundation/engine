package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/x/ownership/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the ownership store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	bankKeeper types.BankKeeper
}

// NewKeeper creates a ownership keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankKeeper types.BankKeeper) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		bankKeeper: bankKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// List returns all ownerships.
func (k *Keeper) List(ctx sdk.Context) ([]*ownership.Ownership, error) {
	var (
		ownerships []*ownership.Ownership
		iter       = ctx.KVStore(k.storeKey).Iterator(nil, nil)
	)
	defer iter.Close()

	for iter.Valid() {
		var o *ownership.Ownership
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &o); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		ownerships = append(ownerships, o)
		iter.Next()
	}
	return ownerships, nil
}

// Set creates a new ownership.
func (k Keeper) Set(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash, resource ownership.Ownership_Resource, resourceAddress sdk.AccAddress) (*ownership.Ownership, error) {
	store := ctx.KVStore(k.storeKey)
	own, err := k.get(store, resourceHash)
	if err != nil {
		return nil, err
	}
	if own != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "resource %s:%q already has an owner", resource, resourceHash)
	}

	own, err = ownership.New(owner.String(), resource, resourceHash, resourceAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	data, err := k.cdc.MarshalBinaryLengthPrefixed(own)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	store.Set(own.Hash, data)
	return own, nil
}

// Delete deletes an ownership.
func (k Keeper) Delete(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash) error {
	store := ctx.KVStore(k.storeKey)
	own, err := k.get(store, resourceHash)
	if err != nil {
		return err
	}
	if own == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "resource %q do not have any ownership", resourceHash)
	}
	if own.Owner != owner.String() {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "resource %q is not owned by you", resourceHash)
	}

	// transfer all spendable coins from resource address to owner
	coins := k.bankKeeper.GetCoins(ctx, own.ResourceAddress)
	if !coins.IsZero() {
		if err := k.bankKeeper.SendCoins(ctx, own.ResourceAddress, owner, coins); err != nil {
			return err
		}
	}

	// remove ownership
	store.Delete(own.Hash)
	return nil
}

// Withdraw try to withdraw coins to owner rom specific resource.
func (k Keeper) Withdraw(ctx sdk.Context, msg types.MsgWithdraw) error {
	store := ctx.KVStore(k.storeKey)
	amount, err := sdk.ParseCoins(msg.Amount)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}
	own, err := k.get(store, msg.ResourceHash)
	if err != nil {
		return err
	}
	if own == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "resource %q does not have any ownership", msg.ResourceHash)
	}
	if own.Owner != msg.Owner.String() {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "resource %q is not owned by you", msg.ResourceHash)
	}
	return k.bankKeeper.SendCoins(ctx, own.ResourceAddress, msg.Owner, amount)
}

// get returns the ownership of a given resource.
func (k Keeper) get(store sdk.KVStore, resourceHash hash.Hash) (*ownership.Ownership, error) {
	iter := store.Iterator(nil, nil)
	var own *ownership.Ownership
	for iter.Valid() {
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &own); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		if own.ResourceHash.Equal(resourceHash) {
			iter.Close()
			return own, nil
		}
		iter.Next()
	}
	iter.Close()
	return nil, nil
}
