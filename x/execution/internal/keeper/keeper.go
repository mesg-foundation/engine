package keeper

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	executionpb "github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	ownershippb "github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/protobuf/api"
	typespb "github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

// price share for the execution runner
const runnerShare = 0.9

// Keeper of the execution store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	supplyKeeper    types.SupplyKeeper
	serviceKeeper   types.ServiceKeeper
	instanceKeeper  types.InstanceKeeper
	runnerKeeper    types.RunnerKeeper
	processKeeper   types.ProcessKeeper
	ownershipKeeper types.OwnershipKeeper
}

// NewKeeper creates a execution keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper types.SupplyKeeper, serviceKeeper types.ServiceKeeper, instanceKeeper types.InstanceKeeper, runnerKeeper types.RunnerKeeper, processKeeper types.ProcessKeeper, ownershipKeeper types.OwnershipKeeper) Keeper {
	keeper := Keeper{
		storeKey:        key,
		cdc:             cdc,
		supplyKeeper:    supplyKeeper,
		serviceKeeper:   serviceKeeper,
		instanceKeeper:  instanceKeeper,
		runnerKeeper:    runnerKeeper,
		processKeeper:   processKeeper,
		ownershipKeeper: ownershipKeeper,
	}
	return keeper
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

	if err := k.supplyKeeper.DelegateCoinsFromAccountToModule(ctx, msg.Signer, types.ModuleName, price); err != nil {
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

	serviceOwners, err := k.ownershipKeeper.GetOwners(ctx, srv.Hash)
	if err != nil {
		return nil, err
	}

	addrs, err := parseAccAddresses(serviceOwners)
	if err != nil {
		return nil, err
	}

	if err := k.distributePriceShares(ctx, msg.Executor, addrs, exec.Price); err != nil {
		return nil, err
	}

	store.Set(exec.Hash, value)
	return exec, nil
}

func (k *Keeper) distributePriceShares(ctx sdk.Context, runnerOwner sdk.AccAddress, serviceOwners []sdk.AccAddress, coinsStr string) error {
	coins, err := sdk.ParseCoins(coinsStr)
	if err != nil {
		return fmt.Errorf("cannot parse coins: %w", err)
	}

	// send all to runner.
	if len(serviceOwners) == 0 {
		return k.supplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, runnerOwner, coins)
	}

	runnerCoins := coinsMulFloat(coins, runnerShare)

	if err := k.supplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, runnerOwner, runnerCoins); err != nil {
		return sdkerrors.Wrapf(err, "cannot send coins %s to runner owner %s", runnerCoins, runnerOwner)
	}

	serviceCoin, serviceRemCoin := coinsQuoRem(coins.Sub(runnerCoins), len(serviceOwners))

	// TODO: the service with rem coins should be randomized
	for i := 0; i < len(serviceOwners)-1; i++ {
		if err := k.supplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, serviceOwners[i], serviceCoin); err != nil {
			return sdkerrors.Wrapf(err, "cannot send coins %s to service owner %s", serviceCoin, serviceOwners[1])
		}
	}

	if err := k.supplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, serviceOwners[len(serviceOwners)-1], serviceRemCoin); err != nil {
		return sdkerrors.Wrapf(err, "cannot send rem coins %s to service owner %s", serviceRemCoin, serviceOwners[len(serviceOwners)-1])
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

// parseAccAddresses parses multiple accAddresses at once.
func parseAccAddresses(addresses []*ownershippb.Ownership) ([]sdk.AccAddress, error) {
	var out []sdk.AccAddress
	for i := range addresses {
		// send to first one for now
		addr, err := sdk.AccAddressFromBech32(addresses[i].Owner)
		if err != nil {
			return nil, err
		}
		out = append(out, addr)
	}
	return out, nil
}

// coinsQuoRem divides coins and returns quotient and remainder.
// NOTE: if reminder is zero then it will be set to quotient.
func coinsQuoRem(coins sdk.Coins, div int) (sdk.Coins, sdk.Coins) {
	var quo, rem sdk.Coins
	for i := range coins {
		q := sdk.Coin{
			Denom:  coins[i].Denom,
			Amount: coins[i].Amount.QuoRaw(int64(div)),
		}
		r := sdk.Coin{
			Denom:  coins[i].Denom,
			Amount: coins[i].Amount.ModRaw(int64(div)),
		}

		quo = append(quo, q)
		if r.IsZero() {
			rem = append(rem, q)
		} else {
			rem = append(rem, r)
		}
	}
	return quo, rem
}

// coinsProcentage multiply coins with float.
func coinsMulFloat(coins sdk.Coins, mul float64) sdk.Coins {
	var ret sdk.Coins
	for i := range coins {
		f := new(big.Float).SetInt(coins[i].Amount.BigInt())
		f.Mul(f, big.NewFloat(mul))
		amt := new(big.Int)
		f.Int(amt)

		c := sdk.Coin{
			Denom:  coins[i].Denom,
			Amount: sdk.NewIntFromBigInt(amt),
		}

		ret = append(ret, c)
	}
	return ret
}
