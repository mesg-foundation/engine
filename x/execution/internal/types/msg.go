package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// MsgCreateExecution defines a state transition to create a execution.
type MsgCreateExecution struct {
	Request *api.CreateExecutionRequest `json:"request"`
	Signer  sdk.AccAddress              `json:"signer"`
}

// NewMsgCreateExecution is a constructor function for MsgCreateExecution.
func NewMsgCreateExecution(req *api.CreateExecutionRequest, signer sdk.AccAddress) *MsgCreateExecution {
	return &MsgCreateExecution{
		Request: req,
		Signer:  signer,
	}
}

// Route should return the name of the module.
func (msg MsgCreateExecution) Route() string {
	return ModuleName
}

// Type returns the action.
func (msg MsgCreateExecution) Type() string {
	return "CreateExecution"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgCreateExecution) ValidateBasic() error {
	price, err := sdk.ParseCoins(msg.Request.Price)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "cannot parse price")
	}
	if price.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "price must be positive")
	}
	if err := xvalidator.Struct(msg); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if !msg.Request.ParentHash.IsZero() && !msg.Request.EventHash.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "cannot have both parent and event hash")
	}
	if msg.Request.ParentHash.IsZero() && msg.Request.EventHash.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "should have at least an event hash or parent hash")
	}
	if msg.Request.ExecutorHash.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "should have a executor hash")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgCreateExecution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgCreateExecution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// MsgUpdateExecution defines a state transition to update a execution.
type MsgUpdateExecution struct {
	Request  *api.UpdateExecutionRequest `json:"request"`
	Executor sdk.AccAddress              `json:"executor"`
}

// NewMsgUpdateExecution is a constructor function for MsgUpdateExecution.
func NewMsgUpdateExecution(req *api.UpdateExecutionRequest, executor sdk.AccAddress) *MsgUpdateExecution {
	return &MsgUpdateExecution{
		Request:  req,
		Executor: executor,
	}
}

// Route should return the name of the module.
func (msg MsgUpdateExecution) Route() string {
	return ModuleName
}

// Type returns the action.
func (msg MsgUpdateExecution) Type() string {
	return "UpdateExecution"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgUpdateExecution) ValidateBasic() error {
	if msg.Executor.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "executor is missing")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgUpdateExecution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgUpdateExecution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Executor}
}
