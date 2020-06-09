package credit

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/x/credit/internal/types"
)

// NewHandler creates an sdk.Handler for all the instance type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgAdd:
			return handleMsgAdd(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgAdd adds credits to an address.
func handleMsgAdd(ctx sdk.Context, k Keeper, msg MsgAdd) (*sdk.Result, error) {
	// check if signer is a minter
	minters := k.Minters(ctx)
	isMinter := false
	for _, minter := range minters {
		if msg.Signer.Equals(minter) {
			isMinter = true
			break
		}
	}
	if !isMinter {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "the signer is not a minter")
	}
	if _, err := k.Add(ctx, msg.Address, msg.Amount); err != nil {
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
		Events: ctx.EventManager().Events(),
	}, nil
}
