package execution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/cosmos/errors"
	"github.com/mesg-foundation/engine/execution"
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

// handleMsgCreateExecution creates a new execution.
func handleMsgCreateExecution(ctx sdk.Context, k Keeper, msg MsgCreateExecution) (*sdk.Result, error) {
	exec, err := k.Create(ctx, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(errors.ErrUnknown, err.Error())
	}

	// TODO: don't emit propsoed event to not break the stream listener in cosmos/client.go#152.
	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	// 		sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeExecutionProposed),
	// 		sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
	// 		sdk.NewAttribute(types.AttributeHash, exec.Hash.String()),
	// 	),
	// )

	// emit EventTypeExecutionInProgress event only when execution status is "in progress"
	if exec.Status == execution.Status_InProgress {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeExecutionInProgress),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
				sdk.NewAttribute(types.AttributeHash, exec.Hash.String()),
			),
		)
	}

	return &sdk.Result{
		Data:   exec.Hash,
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgUpdateExecution updates an execution.
func handleMsgUpdateExecution(ctx sdk.Context, k Keeper, msg MsgUpdateExecution) (*sdk.Result, error) {
	s, err := k.Update(ctx, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(errors.ErrUnknown, err.Error())
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
