package execution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
)

// NewHandler creates an sdk.Handler for all the execution type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateExecution:
			return handleMsgCreateExecution(ctx, k, msg)
		case MsgUpdateExecution:
			return handleMsgUpdateExecution(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgCreateExecution creates a new process.
func handleMsgCreateExecution(ctx sdk.Context, k Keeper, msg MsgCreateExecution) (*sdk.Result, error) {
	s, err := k.Create(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateExecution),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
			sdk.NewAttribute(types.AttributeHash, s.Hash.String()),
		),
	)

	return &sdk.Result{
		Data:   s.Hash,
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgUpdateExecution updates an execution.
func handleMsgUpdateExecution(ctx sdk.Context, k Keeper, msg MsgUpdateExecution) (*sdk.Result, error) {
	s, err := k.Update(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeUpdateExecution),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Executor.String()),
			sdk.NewAttribute(types.AttributeHash, s.Hash.String()),
		),
	)

	return &sdk.Result{
		Data:   s.Hash,
		Events: ctx.EventManager().Events(),
	}, nil
}
