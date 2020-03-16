package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/hash"
)

// MsgWithdrawCoins defines a state transition to create a execution.
type MsgWithdrawCoins struct {
	Hash   hash.Hash      `json:"hash"`
	Amount sdk.Coins      `json:"amount"`
	Owner  sdk.AccAddress `json:"owner"`
}

// NewMsgWithdrawCoins is a constructor function for MsgWithdrawCoins.
func NewMsgWithdrawCoins(h hash.Hash, amount sdk.Coins, owner sdk.AccAddress) *MsgWithdrawCoins {
	return &MsgWithdrawCoins{
		Hash:   h,
		Amount: amount,
		Owner:  owner,
	}
}

// Route should return the name of the module.
func (msg MsgWithdrawCoins) Route() string {
	return ModuleName
}

// Type returns the action.
func (msg MsgWithdrawCoins) Type() string {
	return "WithdrawCoins"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgWithdrawCoins) ValidateBasic() error {
	if msg.Amount.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "price must be positive")
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msg.Owner.String())
	}
	if msg.Hash.IsZero() || !msg.Hash.Valid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid resource address: %s", msg.Owner.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgWithdrawCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgWithdrawCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
