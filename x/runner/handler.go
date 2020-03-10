package runner

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/x/runner/internal/types"
)

// NewHandler creates an sdk.Handler for all the runner type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateRunner:
			return handleMsgCreateRunner(ctx, k, &msg)
		case MsgDeleteRunner:
			return handleMsgDeleteRunner(ctx, k, &msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgCreateRunner creates a new runner.
func handleMsgCreateRunner(ctx sdk.Context, k Keeper, msg *MsgCreateRunner) (*sdk.Result, error) {
	runner, err := k.Create(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateRunner),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
			sdk.NewAttribute(types.AttributeHash, runner.Hash.String()),
		),
	)

	return &sdk.Result{
		Data:   runner.Hash,
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgDeleteRunner deletes a runner.
func handleMsgDeleteRunner(ctx sdk.Context, k Keeper, msg *MsgDeleteRunner) (*sdk.Result, error) {
	if err := k.Delete(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeDeleteRunner),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
			sdk.NewAttribute(types.AttributeHash, msg.RunnerHash.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
