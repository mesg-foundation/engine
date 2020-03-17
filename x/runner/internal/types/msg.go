package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/hash"
)

// MsgCreateRunner defines a state transition to create a runner.
type MsgCreateRunner struct {
	// TODO: address should be renamed to owner
	Address     sdk.AccAddress `json:"address" validate:"required,accaddress"`
	ServiceHash hash.Hash      `json:"serviceHash" validate:"required,hash"`
	EnvHash     hash.Hash      `json:"envHash" validate:"omitempty,hash"`
}

// NewMsgCreateRunner is a constructor function for MsgCreateRunner.
func NewMsgCreateRunner(address sdk.AccAddress, serviceHash hash.Hash, envHash hash.Hash) *MsgCreateRunner {
	return &MsgCreateRunner{
		Address:     address,
		ServiceHash: serviceHash,
		EnvHash:     envHash,
	}
}

// Route should return the name of the module route.
func (msg MsgCreateRunner) Route() string {
	return RouterKey
}

// Type returns the action.
func (msg MsgCreateRunner) Type() string {
	return "CreateRunner"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgCreateRunner) ValidateBasic() error {
	if err := xvalidator.Validate.Struct(msg); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if msg.ServiceHash.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "serviceHash is missing")
	}
	if msg.EnvHash.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "envHash is missing")
	}
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address is missing")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgCreateRunner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgCreateRunner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// MsgDeleteRunner defines a state transition to delete a runner.
type MsgDeleteRunner struct {
	Address    sdk.AccAddress `json:"address" validate:"required,accaddress"`
	RunnerHash hash.Hash      `json:"runnerHash" validate:"required,hash"`
}

// NewMsgDeleteRunner is a constructor function for MsgDeleteRunner.
func NewMsgDeleteRunner(address sdk.AccAddress, runnerHash hash.Hash) *MsgDeleteRunner {
	return &MsgDeleteRunner{
		Address:    address,
		RunnerHash: runnerHash,
	}
}

// Route should return the name of the module.
func (msg MsgDeleteRunner) Route() string {
	return ModuleName
}

// Type returns the action.
func (msg MsgDeleteRunner) Type() string {
	return "DeleteRunner"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgDeleteRunner) ValidateBasic() error {
	if err := xvalidator.Validate.Struct(msg); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if msg.RunnerHash.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "runnerHash is missing")
	}
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address is missing")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgDeleteRunner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgDeleteRunner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
