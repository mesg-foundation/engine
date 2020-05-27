package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/x/credit/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the credit store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper creates a credit keeper
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

// Add a number of credits to an address
func (k Keeper) Add(ctx sdk.Context, address sdk.AccAddress, amount sdk.Int) (sdk.Int, error) {
	value, err := k.Get(ctx, address)
	if err != nil {
		return sdk.Int{}, err
	}
	res := value.Add(amount)
	return k.Set(ctx, address, res)
}

// Sub a number of credits to an address
func (k Keeper) Sub(ctx sdk.Context, address sdk.AccAddress, amount sdk.Int) (sdk.Int, error) {
	value, err := k.Get(ctx, address)
	if err != nil {
		return sdk.Int{}, err
	}
	res := value.Sub(amount)
	return k.Set(ctx, address, res)
}

// Set a number of credit to a specific address
func (k Keeper) Set(ctx sdk.Context, address sdk.AccAddress, amount sdk.Int) (sdk.Int, error) {
	store := ctx.KVStore(k.storeKey)
	encoded, err := k.cdc.MarshalBinaryLengthPrefixed(amount)
	if err != nil {
		return sdk.Int{}, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	store.Set(address.Bytes(), encoded)

	// emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventType,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeActionCreated),
		),
	)
	return amount, nil
}

// Get the amount of a specific address
func (k Keeper) Get(ctx sdk.Context, address sdk.AccAddress) (sdk.Int, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(address.Bytes()) {
		return sdk.ZeroInt(), nil
	}
	value := store.Get(address.Bytes())
	var amount sdk.Int
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(value, &amount); err != nil {
		return sdk.Int{}, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return amount, nil
}
