package serviceapp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetService:
			return handleMsgSetService(ctx, keeper, msg)
		case MsgRemoveService:
			return handleMsgRemoveService(ctx, keeper, msg)
		default:
			errmsg := fmt.Sprintf("Unrecognized service Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errmsg).Result()
		}
	}
}

// handleMsgSetService handles a message to set service.
func handleMsgSetService(ctx sdk.Context, keeper Keeper, msg MsgSetService) sdk.Result {
	keeper.SetService(ctx, Service(msg))
	return sdk.Result{}
}

// handleMsgRemoveService handles a message to remove service.
func handleMsgRemoveService(ctx sdk.Context, keeper Keeper, msg MsgRemoveService) sdk.Result {
	keeper.RemoveService(ctx, msg.Hash)
	return sdk.Result{}
}
