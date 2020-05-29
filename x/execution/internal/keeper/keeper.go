package keeper

import (
	"errors"
	"fmt"
	"math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	executionpb "github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	typespb "github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
	"github.com/tendermint/tendermint/libs/log"
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
	creditKeeper   types.CreditKeeper
	paramstore     params.Subspace
}

// NewKeeper creates a execution keeper
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	bankKeeper types.BankKeeper,
	serviceKeeper types.ServiceKeeper,
	instanceKeeper types.InstanceKeeper,
	runnerKeeper types.RunnerKeeper,
	processKeeper types.ProcessKeeper,
	creditKeeper types.CreditKeeper,
	paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:       key,
		cdc:            cdc,
		bankKeeper:     bankKeeper,
		serviceKeeper:  serviceKeeper,
		instanceKeeper: instanceKeeper,
		runnerKeeper:   runnerKeeper,
		processKeeper:  processKeeper,
		creditKeeper:   creditKeeper,
		paramstore:     paramstore.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Create creates a new execution with proposed status.
// The execution reaches consensus only when more than 2/3 of emitters proposed the same execution.
// TODO: we should split the message and keeper function of execution create from user and for process.
//nolint:gocyclo
func (k *Keeper) Create(ctx sdk.Context, msg types.MsgCreate) (*executionpb.Execution, error) {
	run, err := k.runnerKeeper.Get(ctx, msg.ExecutorHash)
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
	if err := srv.RequireTaskInputs(msg.TaskKey, msg.Inputs); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}
	var proc *process.Process
	if !msg.ProcessHash.IsZero() {
		proc, err = k.processKeeper.Get(ctx, msg.ProcessHash)
		if err != nil {
			return nil, err
		}
	}

	if proc == nil && run.Owner != msg.Signer.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "signer is not the execution's executor")
	}

	exec, err := executionpb.New(
		msg.ProcessHash,
		run.InstanceHash,
		msg.ParentHash,
		msg.EventHash,
		msg.NodeKey,
		msg.TaskKey,
		msg.Inputs,
		msg.Tags,
		msg.ExecutorHash,
	)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// check if exec already exists
	store := ctx.KVStore(k.storeKey)
	if store.Has(exec.Hash) {
		if exec, err = k.Get(ctx, exec.Hash); err != nil {
			return nil, err
		}
		if exec.Status != executionpb.Status_Proposed {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "the execution's status %q should be proposed", exec.Status)
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
				return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "no runner is running the instance that should have trigger this execution")
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
		if runEmitter.Owner == msg.Signer.String() {
			emitterIsPresent = true
			emitter.BlockHeight = ctx.BlockHeight()

			// emit event with action proposed
			event := sdk.NewEvent(
				types.EventType,
				sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeActionProposed),
				sdk.NewAttribute(types.AttributeKeyHash, exec.Hash.String()),
				sdk.NewAttribute(types.AttributeKeyAddress, exec.Address.String()),
				sdk.NewAttribute(types.AttributeKeyExecutor, exec.ExecutorHash.String()),
				sdk.NewAttribute(types.AttributeKeyInstance, exec.InstanceHash.String()),
			)
			if !exec.ProcessHash.IsZero() {
				event = event.AppendAttributes(
					sdk.NewAttribute(types.AttributeKeyExecutor, exec.ProcessHash.String()),
				)
			}
			ctx.EventManager().EmitEvent(event)

			break
		}
	}
	if !emitterIsPresent {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "message's signer is not in the execution's emitters")
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
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
		}

		// emit event
		event := sdk.NewEvent(
			types.EventType,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeActionCreated),
			sdk.NewAttribute(types.AttributeKeyHash, exec.Hash.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, exec.Address.String()),
			sdk.NewAttribute(types.AttributeKeyExecutor, exec.ExecutorHash.String()),
			sdk.NewAttribute(types.AttributeKeyInstance, exec.InstanceHash.String()),
		)
		if !exec.ProcessHash.IsZero() {
			event = event.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyExecutor, exec.ProcessHash.String()),
			)
		}
		ctx.EventManager().EmitEvent(event)

		if !ctx.IsCheckTx() {
			M.InProgress.Add(1)
		}
	}

	// save the exec
	value, err := k.cdc.MarshalBinaryLengthPrefixed(exec)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	store.Set(exec.Hash, value)

	return exec, nil
}

// Update updates a new execution from definition.
func (k *Keeper) Update(ctx sdk.Context, msg types.MsgUpdate) (*executionpb.Execution, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(msg.Hash) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "execution %q doesn't exist", msg.Hash)
	}
	var exec *executionpb.Execution
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(msg.Hash), &exec); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// check if signer is the executor
	runExecutor, err := k.runnerKeeper.Get(ctx, exec.ExecutorHash)
	if err != nil {
		return nil, err
	}
	if runExecutor.Owner != msg.Executor.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "signer is not the execution's executor")
	}

	eventAction := types.AttributeActionFailed
	switch res := msg.Result.(type) {
	case *types.MsgUpdate_Outputs:
		if err := k.validateOutput(ctx, exec.InstanceHash, exec.TaskKey, res.Outputs); err != nil {
			if err1 := exec.Fail(err); err1 != nil {
				return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err1.Error())
			}
		} else if err := exec.Complete(res.Outputs); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
		}
		eventAction = types.AttributeActionCompleted
	case *types.MsgUpdate_Error:
		if err := exec.Fail(errors.New(res.Error)); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
		}
	default:
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "no execution result supplied")
	}

	exec.Start = msg.Start
	exec.Stop = msg.Stop
	price, err := k.calculatePrice(ctx, exec)
	if err != nil {
		return nil, err
	}
	exec.Price = price.String()

	value, err := k.cdc.MarshalBinaryLengthPrefixed(exec)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}

	if !ctx.IsCheckTx() {
		M.Completed.Add(1)
	}

	from := msg.Executor
	if !exec.ProcessHash.IsZero() {
		proc, err := k.processKeeper.Get(ctx, exec.ProcessHash)
		if err != nil {
			return nil, err
		}
		from = proc.PaymentAddress
	}

	if err := k.creditKeeper.Transfer(ctx, from, runExecutor.Address, price); err != nil {
		return nil, err
	}

	store.Set(exec.Hash, value)

	// emit event
	event := sdk.NewEvent(
		types.EventType,
		sdk.NewAttribute(sdk.AttributeKeyAction, eventAction),
		sdk.NewAttribute(types.AttributeKeyHash, exec.Hash.String()),
		sdk.NewAttribute(types.AttributeKeyAddress, exec.Address.String()),
		sdk.NewAttribute(types.AttributeKeyExecutor, exec.ExecutorHash.String()),
		sdk.NewAttribute(types.AttributeKeyInstance, exec.InstanceHash.String()),
	)
	if !exec.ProcessHash.IsZero() {
		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyExecutor, exec.ProcessHash.String()),
		)
	}
	ctx.EventManager().EmitEvent(event)

	return exec, nil
}

// Get returns the execution that matches given hash.
func (k *Keeper) Get(ctx sdk.Context, hash hash.Hash) (*executionpb.Execution, error) {
	var exec *executionpb.Execution
	store := ctx.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "execution %q not found", hash)
	}
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(hash), &exec); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return exec, nil
}

// List returns all executions.
func (k *Keeper) List(ctx sdk.Context, filter types.ListFilter) ([]*executionpb.Execution, error) {
	var (
		execs []*executionpb.Execution
		iter  = ctx.KVStore(k.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var exec *executionpb.Execution
		value := iter.Value()
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(value, &exec); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		if filter.Match(exec) {
			execs = append(execs, exec)
		}
		iter.Next()
	}
	iter.Close()
	return execs, nil
}

func (k *Keeper) validateOutput(ctx sdk.Context, instanceHash hash.Hash, taskKey string, outputs *typespb.Struct) error {
	inst, err := k.instanceKeeper.Get(ctx, instanceHash)
	if err != nil {
		return err
	}
	srv, err := k.serviceKeeper.Get(ctx, inst.ServiceHash)
	if err != nil {
		return err
	}
	if err := srv.RequireTaskOutputs(taskKey, outputs); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
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
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}
	var instanceHash hash.Hash
	switch n := parentNode.GetType().(type) {
	case *process.Process_Node_Event_:
		instanceHash = n.Event.InstanceHash
	case *process.Process_Node_Result_:
		instanceHash = n.Result.InstanceHash
	default:
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "parent node type should be an event or a result")
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

func (k *Keeper) calculatePrice(ctx sdk.Context, exec *executionpb.Execution) (sdk.Int, error) {
	inst, err := k.instanceKeeper.Get(ctx, exec.InstanceHash)
	if err != nil {
		return sdk.Int{}, err
	}
	srv, err := k.serviceKeeper.Get(ctx, inst.ServiceHash)
	if err != nil {
		return sdk.Int{}, err
	}
	task, err := srv.GetTask(exec.TaskKey)
	if err != nil {
		return sdk.Int{}, err
	}
	perCall := sdk.NewInt(200)
	if task.Price != nil && task.Price.PerCall != "" {
		price, ok := sdk.NewIntFromString(task.Price.PerCall)
		if !ok {
			return sdk.Int{}, fmt.Errorf("invalid price per call %q", task.Price.PerCall)
		}
		if price.GT(perCall) {
			perCall = price
		}
	}
	perSec := sdk.NewInt(20000)
	if task.Price != nil && task.Price.PerSec != "" {
		price, ok := sdk.NewIntFromString(task.Price.PerSec)
		if !ok {
			return sdk.Int{}, fmt.Errorf("invalid price per sec %q", task.Price.PerSec)
		}
		if price.GT(perSec) {
			perSec = price
		}
	}
	perKB := sdk.NewInt(100)
	if task.Price != nil && task.Price.PerKB != "" {
		price, ok := sdk.NewIntFromString(task.Price.PerKB)
		if !ok {
			return sdk.Int{}, fmt.Errorf("invalid price per kb %q", task.Price.PerKB)
		}
		if price.GT(perKB) {
			perKB = price
		}
	}
	duration := sdk.NewInt(exec.GetDuration())
	inputs, err := k.cdc.MarshalBinaryLengthPrefixed(exec.Inputs)
	if err != nil {
		return sdk.Int{}, err
	}
	outputs, err := k.cdc.MarshalBinaryLengthPrefixed(exec.Outputs)
	if err != nil {
		return sdk.Int{}, err
	}
	datasize := sdk.NewInt(int64(math.Ceil(float64(len(inputs)+len(outputs)) / 1000)))
	return perCall.Add(perSec.Mul(duration)).Add(perKB.Mul(datasize)), nil
} 
 
// Import imports a list of executions into the store.
func (k *Keeper) Import(ctx sdk.Context, execs []*executionpb.Execution) error {
	store := ctx.KVStore(k.storeKey)
	for _, exec := range execs {
		value, err := k.cdc.MarshalBinaryLengthPrefixed(exec)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
		}
		store.Set(exec.Hash, value)
	}
	return nil
}
