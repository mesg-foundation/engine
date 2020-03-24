package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	processpb "github.com/mesg-foundation/engine/process"
)

// Route should return the name of the module route.
func (msg MsgCreate) Route() string {
	return RouterKey
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
	p := processpb.New(msg.Name, msg.Nodes, msg.Edges)
	if err := p.Validate(); err != nil {
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
	return []sdk.AccAddress{msg.Owner}
}

// Route should return the name of the module.
func (msg MsgDelete) Route() string {
	return ModuleName
}

// Type returns the action.
func (msg MsgDelete) Type() string {
	return "delete"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgDelete) ValidateBasic() error {
	if err := xvalidator.Struct(msg); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgDelete) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgDelete) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
