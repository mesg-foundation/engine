package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/hash"
	ownershippb "github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/process"
	processpb "github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/x/process/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

var processCreateInitialBalance = "10000000atto"

// Keeper of the process store
type Keeper struct {
	storeKey        sdk.StoreKey
	cdc             *codec.Codec
	ownershipKeeper types.OwnershipKeeper
	instanceKeeper  types.InstanceKeeper
	bankKeeper      types.BankKeeper
}

// NewKeeper creates a process keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, instanceKeeper types.InstanceKeeper, ownershipKeeper types.OwnershipKeeper, bankKeeper types.BankKeeper) Keeper {
	keeper := Keeper{
		storeKey:        key,
		cdc:             cdc,
		instanceKeeper:  instanceKeeper,
		ownershipKeeper: ownershipKeeper,
		bankKeeper:      bankKeeper,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Create creates a new process.
func (k Keeper) Create(ctx sdk.Context, msg *types.MsgCreate) (*processpb.Process, error) {
	store := ctx.KVStore(k.storeKey)

	p, err := process.New(msg.Name, msg.Nodes, msg.Edges)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if store.Has(p.Hash) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "process %q already exists", p.Hash)
	}

	for _, node := range p.Nodes {
		switch n := node.Type.(type) {
		case *processpb.Process_Node_Result_:
			if _, err := k.instanceKeeper.Get(ctx, n.Result.InstanceHash); err != nil {
				return nil, err
			}
		case *processpb.Process_Node_Event_:
			if _, err := k.instanceKeeper.Get(ctx, n.Event.InstanceHash); err != nil {
				return nil, err
			}
		case *processpb.Process_Node_Task_:
			if _, err := k.instanceKeeper.Get(ctx, n.Task.InstanceHash); err != nil {
				return nil, err
			}
		}
	}

	procInitBal, err := sdk.ParseCoins(processCreateInitialBalance)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, err.Error())
	}
	if err := k.bankKeeper.SendCoins(ctx, msg.Owner, p.Address, procInitBal); err != nil {
		return nil, err
	}

	value, err := k.cdc.MarshalBinaryLengthPrefixed(p)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}

	if _, err := k.ownershipKeeper.Set(ctx, msg.Owner, p.Hash, ownershippb.Ownership_Process, p.Address); err != nil {
		return nil, err
	}

	store.Set(p.Hash, value)

	// emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventType,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeActionCreated),
			sdk.NewAttribute(types.AttributeKeyHash, p.Hash.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, p.Address.String()),
		),
	)

	return p, nil
}

// Delete deletes a process.
func (k Keeper) Delete(ctx sdk.Context, msg *types.MsgDelete) error {
	p, err := k.Get(ctx, msg.Hash)
	if err != nil {
		return err
	}
	if err := k.ownershipKeeper.Delete(ctx, msg.Owner, msg.Hash); err != nil {
		return err
	}
	ctx.KVStore(k.storeKey).Delete(msg.Hash)

	// emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventType,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeActionDeleted),
			sdk.NewAttribute(types.AttributeKeyHash, p.Hash.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, p.Address.String()),
		),
	)

	return nil
}

// Get returns the process that matches given hash.
func (k Keeper) Get(ctx sdk.Context, hash hash.Hash) (*processpb.Process, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "process %q not found", hash)
	}
	value := store.Get(hash)
	var p *processpb.Process
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(value, &p); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return p, nil
}

// Exists returns true if a specific set of data exists in the database, false otherwise
func (k Keeper) Exists(ctx sdk.Context, hash hash.Hash) (bool, error) {
	return ctx.KVStore(k.storeKey).Has(hash), nil
}

// List returns all processes.
func (k Keeper) List(ctx sdk.Context) ([]*processpb.Process, error) {
	var (
		processes []*processpb.Process
		iter      = ctx.KVStore(k.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var p *processpb.Process
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &p); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		processes = append(processes, p)
		iter.Next()
	}
	iter.Close()
	return processes, nil
}

// Import imports a list of processes into the store.
func (k *Keeper) Import(ctx sdk.Context, processes []*process.Process) error {
	store := ctx.KVStore(k.storeKey)
	for _, proc := range processes {
		value, err := k.cdc.MarshalBinaryLengthPrefixed(proc)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
		}
		store.Set(proc.Hash, value)
	}
	return nil
}
