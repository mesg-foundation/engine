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
	"github.com/tendermint/tendermint/crypto"
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
func (k Keeper) Create(ctx sdk.Context, msg *types.MsgCreateProcess) (*processpb.Process, error) {
	store := ctx.KVStore(k.storeKey)
	p := &process.Process{
		Name:  msg.Request.Name,
		Nodes: msg.Request.Nodes,
		Edges: msg.Request.Edges,
	}
	p.Hash = hash.Dump(p)
	p.Address = sdk.AccAddress(crypto.AddressHash(p.Hash)).String()
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

	procAddress := sdk.AccAddress(crypto.AddressHash(p.Hash))
	procInitBal, err := sdk.ParseCoins(processCreateInitialBalance)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if err := k.bankKeeper.SendCoins(ctx, msg.Owner, procAddress, procInitBal); err != nil {
		return nil, err
	}

	value, err := k.cdc.MarshalBinaryLengthPrefixed(p)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}

	if _, err := k.ownershipKeeper.Set(ctx, msg.Owner, p.Hash, ownershippb.Ownership_Process); err != nil {
		return nil, err
	}

	store.Set(p.Hash, value)
	return p, nil
}

// Delete deletes a process.
func (k Keeper) Delete(ctx sdk.Context, msg *types.MsgDeleteProcess) error {
	if err := k.ownershipKeeper.Delete(ctx, msg.Owner, msg.Request.Hash); err != nil {
		return err
	}
	ctx.KVStore(k.storeKey).Delete(msg.Request.Hash)
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
