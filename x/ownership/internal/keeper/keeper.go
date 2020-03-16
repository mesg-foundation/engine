package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/x/ownership/internal/types"
	"github.com/tendermint/tendermint/crypto"
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
func (k Keeper) Set(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash, resource ownership.Ownership_Resource) (*ownership.Ownership, error) {
	store := ctx.KVStore(k.storeKey)
	hashes, err := k.findOwnerships(store, "", resourceHash)
	if err != nil {
		return nil, err
	}
	if len(hashes) > 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "resource %s:%q already has an owner", resource, resourceHash)
	}

	ownership := &ownership.Ownership{
		Owner:        owner.String(),
		Resource:     resource,
		ResourceHash: resourceHash,
	}
	ownership.Hash = hash.Dump(ownership)

	data, err := k.cdc.MarshalBinaryLengthPrefixed(ownership)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	store.Set(ownership.Hash, data)
	return ownership, nil
}

// Delete deletes an ownership.
func (k Keeper) Delete(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash) error {
	store := ctx.KVStore(k.storeKey)
	hashes, err := k.findOwnerships(store, owner.String(), resourceHash)
	if err != nil {
		return err
	}
	if len(hashes) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "resource %q do not have any ownership", resourceHash)
	}

	// transfer all spendable coins from resource address to owner
	addr := sdk.AccAddress(crypto.AddressHash(resourceHash))
	coins := k.bankKeeper.GetCoins(ctx, addr)
	if !coins.IsZero() {
		if err := k.bankKeeper.SendCoins(ctx, addr, owner, coins); err != nil {
			return err
		}
	}

	// remove all ownerships
	for _, hash := range hashes {
		store.Delete(hash)
	}
	return nil
}

// WithdrawCoins try to withdraw coins to owner rom specific resource.
func (k Keeper) WithdrawCoins(ctx sdk.Context, msg types.MsgWithdrawCoins) error {
	ownerships, err := k.findOwnerships(ctx.KVStore(k.storeKey), msg.Owner.String(), msg.Hash)
	if err != nil {
		return err
	}
	if len(ownerships) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "address %s is not owner of resource %s", msg.Owner, msg.Hash)
	}
	addr := sdk.AccAddress(crypto.AddressHash(msg.Hash))
	return k.bankKeeper.SendCoins(ctx, addr, msg.Owner, msg.Amount)
}

// hasOwner checks if given resource has owner. Returns all ownership hash and true if has one
// nil and false otherwise.
func (k Keeper) findOwnerships(store sdk.KVStore, owner string, resourceHash hash.Hash) ([]hash.Hash, error) {
	var (
		ownerships []hash.Hash
		iter       = store.Iterator(nil, nil)
	)
	for iter.Valid() {
		var o *ownership.Ownership
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &o); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		if (owner == "" || o.Owner == owner) && o.ResourceHash.Equal(resourceHash) {
			ownerships = append(ownerships, o.Hash)
		}
		iter.Next()
	}
	iter.Close()
	return ownerships, nil
}
