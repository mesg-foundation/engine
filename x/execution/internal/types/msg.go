package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/ext/xvalidator"
)

// Route should return the name of the module.
func (msg MsgCreate) Route() string {
	return ModuleName
}

// Type returns the action.
func (msg MsgCreate) Type() string {
	return "create"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgCreate) ValidateBasic() error {
	if err := xvalidator.Struct(msg); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// Route should return the name of the module.
func (msg MsgUpdate) Route() string {
	return ModuleName
}

// Type returns the action.
func (msg MsgUpdate) Type() string {
	return "update"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgUpdate) ValidateBasic() error {
	if err := xvalidator.Struct(msg); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgUpdate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Executor}
}
