package keeper

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	executionpb "github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	typespb "github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

// price share for the execution runner
// TODO: this should be a param
var (
	runnerShare  = sdk.NewDecWithPrec(8, 1)
	serviceShare = sdk.NewDecWithPrec(1, 1)
)

// Keeper of the execution store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	bankKeeper     types.BankKeeper
	serviceKeeper  types.ServiceKeeper
	instanceKeeper types.InstanceKeeper
	runnerKeeper   types.RunnerKeeper
	processKeeper  types.ProcessKeeper
	paramstore     params.Subspace
}

// NewKeeper creates a execution keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankKeeper types.BankKeeper, serviceKeeper types.ServiceKeeper, instanceKeeper types.InstanceKeeper, runnerKeeper types.RunnerKeeper, processKeeper types.ProcessKeeper, paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:       key,
		cdc:            cdc,
		bankKeeper:     bankKeeper,
		serviceKeeper:  serviceKeeper,
		instanceKeeper: instanceKeeper,
		runnerKeeper:   runnerKeeper,
		processKeeper:  processKeeper,
		paramstore:     paramstore.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Create creates a new execution with proposed status.
// The execution reaches consensus only when even emitter proposed the same execution.
// TODO: we should split the message and keeper function of execution create from user and for process.
//nolint:gocyclo
func (k *Keeper) Create(ctx sdk.Context, msg types.MsgCreateExecution) (*executionpb.Execution, error) {
	price, err := sdk.ParseCoins(msg.Request.Price)
	if err != nil {
		return nil, err
	}

	minPriceCoin, err := sdk.ParseCoins(k.MinPrice(ctx))
	if err != nil {
		return nil, err
	}
	if !price.IsAllGTE(minPriceCoin) {
		return nil, fmt.Errorf("execution price too low. Min value: %q", minPriceCoin.String())
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
	if err := srv.RequireTaskInputs(msg.Request.TaskKey, msg.Request.Inputs); err != nil {
		return nil, err
	}
	var proc *process.Process
	if !msg.Request.ProcessHash.IsZero() {
		proc, err = k.processKeeper.Get(ctx, msg.Request.ProcessHash)
		if err != nil {
			return nil, err
		}
	}

	if proc == nil && run.Address != msg.Signer.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "signer is not the execution's executor")
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

	// check if exec already exists
	store := ctx.KVStore(k.storeKey)
	if store.Has(exec.Hash) {
		if exec, err = k.Get(ctx, exec.Hash); err != nil {
			return nil, err
		}
		if exec.Status != executionpb.Status_Proposed {
			return nil, fmt.Errorf("the execution's status %q should be proposed", exec.Status)
		}
	} else {
		// set everything that should be put at the creation of the exec

		// init execution's emitters
		if proc == nil {
			// no process, set the signer as only emitter
			exec.Emitters = []*executionpb.Execution_Emitter{{
				RunnerHash:  run.Hash,
				BlockHeight: 0,
			}}
		} else {
			matchedRuns, err := k.fetchEmitters(ctx, proc, exec.NodeKey)
			if err != nil {
				return nil, err
			}
			if len(matchedRuns) == 0 {
				return nil, fmt.Errorf("no runner is running the instance that should have trigger this execution")
			}
			for _, matchedRun := range matchedRuns {
				exec.Emitters = append(exec.Emitters, &executionpb.Execution_Emitter{
					RunnerHash:  matchedRun.Hash,
					BlockHeight: 0,
				})
			}
		}
	}

	// check if signer is in emitters, set emitter's block height, return error if not present.
	emitterIsPresent := false
	for _, emitter := range exec.Emitters {
		runEmitter, err := k.runnerKeeper.Get(ctx, emitter.RunnerHash)
		if err != nil {
			return nil, err
		}
		if runEmitter.Address == msg.Signer.String() {
			emitterIsPresent = true
			emitter.BlockHeight = ctx.BlockHeight()
			break
		}
	}
	if !emitterIsPresent {
		return nil, fmt.Errorf("message's signer is not in the execution's emitters")
	}

	// define consensus requirements
	nbrEmitterRequired := int(math.Ceil(float64(len(exec.Emitters)) * 2 / 3))

	// calculate emitter already proposed
	nbrProposedEmitters := 0
	for _, emitter := range exec.Emitters {
		if emitter.BlockHeight > 0 {
			nbrProposedEmitters++
		}
	}

	// check if emitter consensus is reached
	if nbrProposedEmitters == nbrEmitterRequired {
		// set block height when consensus is reached
		exec.BlockHeight = ctx.BlockHeight()

		// change the status of the exec
		if err := exec.Execute(); err != nil {
			return nil, err
		}
		if !ctx.IsCheckTx() {
			M.InProgress.Add(1)
		}

		// transfer the coin either from the process or from the signer
		from := msg.Signer
		if !msg.Request.ProcessHash.IsZero() {
			from = sdk.AccAddress(crypto.AddressHash(msg.Request.ProcessHash))
		}
		execAddress := sdk.AccAddress(crypto.AddressHash(exec.Hash))
		if err := k.bankKeeper.SendCoins(ctx, from, execAddress, price); err != nil {
			return nil, err
		}
	}

	// save the exec
	value, err := k.cdc.MarshalBinaryLengthPrefixed(exec)
	if err != nil {
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

	// check if signer is the executor
	runExecutor, err := k.runnerKeeper.Get(ctx, exec.ExecutorHash)
	if err != nil {
		return nil, err
	}
	if runExecutor.Address != msg.Executor.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "signer is not the execution's executor")
	}

	switch res := msg.Request.Result.(type) {
	case *api.UpdateExecutionRequest_Outputs:
		if err := k.validateExecutionOutput(ctx, exec.InstanceHash, exec.TaskKey, res.Outputs); err != nil {
			if err1 := exec.Fail(err); err1 != nil {
				return nil, err1
			}
		} else if err := exec.Complete(res.Outputs); err != nil {
			return nil, err
		}
	case *api.UpdateExecutionRequest_Error:
		if err := exec.Fail(errors.New(res.Error)); err != nil {
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

	if err := k.distributePriceShares(ctx, exec.Hash, exec.Emitters, exec.ExecutorHash, inst.ServiceHash, exec.Price); err != nil {
		return nil, err
	}

	store.Set(exec.Hash, value)
	return exec, nil
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

func (k *Keeper) distributePriceShares(ctx sdk.Context, execHash hash.Hash, emitters []*executionpb.Execution_Emitter, runnerHash, serviceHash hash.Hash, price string) error {
	coins, err := sdk.ParseCoins(price)
	if err != nil {
		return fmt.Errorf("cannot parse coins: %w", err)
	}
	if coins.Empty() {
		return nil
	}

	// addresses
	execAddress := sdk.AccAddress(crypto.AddressHash(execHash))
	runnerAddress := sdk.AccAddress(crypto.AddressHash(runnerHash))
	serviceAddress := sdk.AccAddress(crypto.AddressHash(serviceHash))

	// inputs
	inputs := []bank.Input{bank.NewInput(execAddress, coins)}

	// outputs
	runnerCoins, _ := sdk.NewDecCoinsFromCoins(coins...).MulDecTruncate(runnerShare).TruncateDecimal()
	serviceCoins, _ := sdk.NewDecCoinsFromCoins(coins...).MulDecTruncate(serviceShare).TruncateDecimal()

	outputs := []bank.Output{
		bank.NewOutput(runnerAddress, runnerCoins),
		bank.NewOutput(serviceAddress, serviceCoins),
	}

	// emitters get all the remaining coins
	emittersCoins := coins.Sub(runnerCoins).Sub(serviceCoins)
	distributedEmittersCoins := sdk.NewCoins()

	for i, emitter := range emitters {
		emitterAddress := sdk.AccAddress(crypto.AddressHash(emitter.RunnerHash))
		if err != nil {
			return err
		}
		// 1 emitter share = 1 / len(emitters)
		emitterShare := sdk.NewDecFromBigInt(big.NewInt(1).Div(big.NewInt(1), big.NewInt(int64(len(emitters)))))
		emitterCoins, _ := sdk.NewDecCoinsFromCoins(emittersCoins...).MulDecTruncate(emitterShare).TruncateDecimal()
		distributedEmittersCoins = distributedEmittersCoins.Add(emitterCoins...)

		// give the remaining emitter to the last emitter
		if i == len(emitters)-1 {
			emitterCoins = emitterCoins.Add(emittersCoins.Sub(distributedEmittersCoins)...)
		}

		outputs = append(outputs, bank.NewOutput(emitterAddress, emitterCoins))
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

// fetchEmitters returns the runners running the instance that was responsible for creating this execution from the process.
func (k *Keeper) fetchEmitters(ctx sdk.Context, proc *process.Process, nodeKey string) ([]*runner.Runner, error) {
	// get parent node's instance hash
	parentNode, err := proc.FindParentWithType(nodeKey, func(node *process.Process_Node) bool {
		switch node.GetType().(type) {
		case *process.Process_Node_Event_, *process.Process_Node_Result_:
			return true
		default:
			return false
		}
	})
	if err != nil {
		return nil, err
	}
	var instanceHash hash.Hash
	switch n := parentNode.GetType().(type) {
	case *process.Process_Node_Event_:
		instanceHash = n.Event.InstanceHash
	case *process.Process_Node_Result_:
		instanceHash = n.Result.InstanceHash
	default:
		return nil, fmt.Errorf("parent node type should be an event or a result")
	}

	// get runners of this instance
	runs, err := k.runnerKeeper.List(ctx)
	if err != nil {
		return nil, err
	}
	matchedRuns := make([]*runner.Runner, 0)
	for _, run := range runs {
		if run.InstanceHash.Equal(instanceHash) {
			matchedRuns = append(matchedRuns, run)
		}
	}
	return matchedRuns, nil
}
