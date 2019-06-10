package service

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) types.Handler {
	return func(ctx types.Context, msg types.Msg) types.Result {
		switch msg := msg.(type) {
		case MsgSetService:
			return handleMsgSetService(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return types.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgSetService(ctx types.Context, keeper Keeper, msg MsgSetService) types.Result {
	owner, _ := keeper.GetOwner(ctx, msg.Hash)
	if !msg.Owner.Equals(owner) { // Checks if the the msg sender is the same as the current owner
		return types.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.SetService(ctx, msg.Hash, msg.Service) // If so, set the name to the value specified in the msg.
	return types.Result{}                         // return
}
