package service

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/gogo/protobuf/codec"
)

// Keeper ...
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey types.StoreKey // Unexposed key to access store from sdk.Context

	cdc codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey types.StoreKey, cdc codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// GetService ...
func (k Keeper) GetService(ctx types.Context, hash string) (Service, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(hash)) {
		return NewService(), nil
	}
	bz := store.Get([]byte(hash))
	var service Service
	return service, k.cdc.Unmarshal(bz, &service)
}

// GetOwner ...
func (k Keeper) GetOwner(ctx types.Context, hash string) (types.AccAddress, error) {
	service, err := k.GetService(ctx, hash)
	if err != nil {
		return nil, err
	}
	return service.Owner, nil
}

// SetService ...
func (k Keeper) SetService(ctx types.Context, hash string, service Service) error {
	store := ctx.KVStore(k.storeKey)
	s, err := k.cdc.Marshal(service)
	if err != nil {
		return err
	}
	store.Set([]byte(hash), s)
	return nil
}
