package keeper

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	executionpb "github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	typespb "github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

// price share for the execution runner
var runnerShare = sdk.NewDecWithPrec(9, 1)

// Keeper of the execution store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	bankKeeper     types.BankKeeper
	serviceKeeper  types.ServiceKeeper
	instanceKeeper types.InstanceKeeper
	runnerKeeper   types.RunnerKeeper
	processKeeper  types.ProcessKeeper
}

// NewKeeper creates a execution keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankKeeper types.BankKeeper, serviceKeeper types.ServiceKeeper, instanceKeeper types.InstanceKeeper, runnerKeeper types.RunnerKeeper, processKeeper types.ProcessKeeper) Keeper {
	return Keeper{
		storeKey:       key,
		cdc:            cdc,
		bankKeeper:     bankKeeper,
		serviceKeeper:  serviceKeeper,
		instanceKeeper: instanceKeeper,
		runnerKeeper:   runnerKeeper,
		processKeeper:  processKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Create creates a new execution from definition.
func (k *Keeper) Create(ctx sdk.Context, msg types.MsgCreateExecution) (*executionpb.Execution, error) {
	price, err := sdk.ParseCoins(msg.Request.Price)
	if err != nil {
		return nil, err
	}
	run, err := k.runnerKeeper.Get(ctx, msg.Request.ExecutorHash)
	if err != nil {
		return nil, err
	}
	inst, err := k.instanceKeeper.Get(ctx, run.InstanceHash)
	if err != nil {
		return nil, err
	}
	srv, err := k.serviceKeeper.Get(ctx, inst.ServiceHash)
	if err != nil {
		return nil, err
	}
	if !msg.Request.ProcessHash.IsZero() {
		if _, err := k.processKeeper.Get(ctx, msg.Request.ProcessHash); err != nil {
			return nil, err
		}
	}
	if err := srv.RequireTaskInputs(msg.Request.TaskKey, msg.Request.Inputs); err != nil {
		return nil, err
	}
	exec := executionpb.New(
		msg.Request.ProcessHash,
		run.InstanceHash,
		msg.Request.ParentHash,
		msg.Request.EventHash,
		msg.Request.NodeKey,
		msg.Request.TaskKey,
		msg.Request.Price,
		msg.Request.Inputs,
		msg.Request.Tags,
		msg.Request.ExecutorHash,
	)
	store := ctx.KVStore(k.storeKey)
	if store.Has(exec.Hash) {
		return nil, fmt.Errorf("execution %q already exists", exec.Hash)
	}
	if err := exec.Execute(); err != nil {
		return nil, err
	}
	value, err := k.cdc.MarshalBinaryLengthPrefixed(exec)
	if err != nil {
		return nil, err
	}

	if !ctx.IsCheckTx() {
		M.InProgress.Add(1)
	}

	execAddress := sdk.AccAddress(crypto.AddressHash(exec.Hash))
	if err := k.bankKeeper.SendCoins(ctx, msg.Signer, execAddress, price); err != nil {
		return nil, err
	}

	store.Set(exec.Hash, value)
	return exec, nil
}

// Update updates a new execution from definition.
func (k *Keeper) Update(ctx sdk.Context, msg types.MsgUpdateExecution) (*executionpb.Execution, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(msg.Request.Hash) {
		return nil, fmt.Errorf("execution %q doesn't exist", msg.Request.Hash)
	}
	var exec *executionpb.Execution
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(msg.Request.Hash), &exec); err != nil {
		return nil, err
	}
	switch res := msg.Request.Result.(type) {
	case *api.UpdateExecutionRequest_Outputs:
		if err := k.validateExecutionOutput(ctx, exec.InstanceHash, exec.TaskKey, res.Outputs); err != nil {
			if err1 := exec.Failed(err); err1 != nil {
				return nil, err1
			}
		} else if err := exec.Complete(res.Outputs); err != nil {
			return nil, err
		}
	case *api.UpdateExecutionRequest_Error:
		if err := exec.Failed(errors.New(res.Error)); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("no execution result supplied")
	}
	value, err := k.cdc.MarshalBinaryLengthPrefixed(exec)
	if err != nil {
		return nil, err
	}

	if !ctx.IsCheckTx() {
		M.Completed.Add(1)
	}

	inst, err := k.instanceKeeper.Get(ctx, exec.InstanceHash)
	if err != nil {
		return nil, err
	}
	srv, err := k.serviceKeeper.Get(ctx, inst.ServiceHash)
	if err != nil {
		return nil, err
	}
	run, err := k.runnerKeeper.Get(ctx, exec.ExecutorHash)
	if err != nil {
		return nil, err
	}

	execAddress := sdk.AccAddress(crypto.AddressHash(exec.Hash))
	runnerAddress := sdk.AccAddress(crypto.AddressHash(run.Hash))
	serviceAddress := sdk.AccAddress(crypto.AddressHash(srv.Hash))
	if err := k.distributePriceShares(ctx, execAddress, runnerAddress, serviceAddress, exec.Price); err != nil {
		return nil, err
	}

	store.Set(exec.Hash, value)
	return exec, nil
}

func (k *Keeper) distributePriceShares(ctx sdk.Context, execAddress, runnerAddress, serviceAddress sdk.AccAddress, price string) error {
	coins, err := sdk.ParseCoins(price)
	if err != nil {
		return fmt.Errorf("cannot parse coins: %w", err)
	}
	if coins.Empty() {
		return nil
	}

	runnerCoins, _ := sdk.NewDecCoinsFromCoins(coins...).MulDecTruncate(runnerShare).TruncateDecimal()
	serviceCoins := coins.Sub(runnerCoins)

	inputs := []bank.Input{bank.NewInput(execAddress, coins)}
	outputs := []bank.Output{
		bank.NewOutput(runnerAddress, runnerCoins),
		bank.NewOutput(serviceAddress, serviceCoins),
	}

	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return sdkerrors.Wrapf(err, "cannot distribute coins from execution adddress %s with inputs %s and outputs %s", execAddress, inputs, outputs)
	}
	return nil
}

func (k *Keeper) validateExecutionOutput(ctx sdk.Context, instanceHash hash.Hash, taskKey string, outputs *typespb.Struct) error {
	inst, err := k.instanceKeeper.Get(ctx, instanceHash)
	if err != nil {
		return err
	}
	srv, err := k.serviceKeeper.Get(ctx, inst.ServiceHash)
	if err != nil {
		return err
	}
	return srv.RequireTaskOutputs(taskKey, outputs)
}

// Get returns the execution that matches given hash.
func (k *Keeper) Get(ctx sdk.Context, hash hash.Hash) (*executionpb.Execution, error) {
	var exec *executionpb.Execution
	store := ctx.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("execution %q not found", hash)
	}
	return exec, k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(hash), &exec)
}

// List returns all executions.
func (k *Keeper) List(ctx sdk.Context) ([]*executionpb.Execution, error) {
	var (
		execs []*executionpb.Execution
		iter  = ctx.KVStore(k.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var exec *executionpb.Execution
		value := iter.Value()
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(value, &exec); err != nil {
			return nil, err
		}
		execs = append(execs, exec)
		iter.Next()
	}
	iter.Close()
	return execs, nil
}
