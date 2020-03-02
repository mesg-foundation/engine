package process

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/x/process/internal/types"
)

// NewHandler creates an sdk.Handler for all the process type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateProcess:
			return handleMsgCreateProcess(ctx, k, &msg)
		case MsgDeleteProcess:
			return handleMsgDeleteProcess(ctx, k, &msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgCreateProcess creates a new process.
func handleMsgCreateProcess(ctx sdk.Context, k Keeper, msg *MsgCreateProcess) (*sdk.Result, error) {
	p, err := k.Create(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateProcess),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeHash, p.Hash.String()),
		),
	)

	return &sdk.Result{
		Data:   p.Hash,
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgDeleteProcess deletes a process.
func handleMsgDeleteProcess(ctx sdk.Context, k Keeper, msg *MsgDeleteProcess) (*sdk.Result, error) {
	if err := k.Delete(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeDeleteProcess),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeHash, msg.Request.Hash.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
