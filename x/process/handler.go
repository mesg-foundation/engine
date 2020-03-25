package process

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/x/process/internal/types"
)

// NewHandler creates an sdk.Handler for all the process type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreate:
			return handleMsgCreate(ctx, k, &msg)
		case MsgDelete:
			return handleMsgDelete(ctx, k, &msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// handleMsgCreate creates a new process.
func handleMsgCreate(ctx sdk.Context, k Keeper, msg *MsgCreate) (*sdk.Result, error) {
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

// handleMsgDelete deletes a process.
func handleMsgDelete(ctx sdk.Context, k Keeper, msg *MsgDelete) (*sdk.Result, error) {
	if err := k.Delete(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeDeleteProcess),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeHash, msg.Hash.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
