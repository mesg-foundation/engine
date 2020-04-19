package execution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
)

// NewHandler creates an sdk.Handler for all the execution type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreate:
			return handleMsgCreate(ctx, k, msg)
		case MsgUpdate:
			return handleMsgUpdate(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// handleMsgCreate creates a new execution.
func handleMsgCreate(ctx sdk.Context, k Keeper, msg MsgCreate) (*sdk.Result, error) {
	exec, err := k.Create(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	)

	return &sdk.Result{
		Data:   exec.Hash,
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgUpdate updates an execution.
func handleMsgUpdate(ctx sdk.Context, k Keeper, msg MsgUpdate) (*sdk.Result, error) {
	exec, err := k.Update(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Executor.String()),
		),
	)

	return &sdk.Result{
		Data:   exec.Hash,
		Events: ctx.EventManager().Events(),
	}, nil
}
