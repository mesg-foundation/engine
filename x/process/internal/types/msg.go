package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/cosmos/address"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	processpb "github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/tendermint/tendermint/crypto"
)

// MsgCreateProcess defines a state transition to create a process.
type MsgCreateProcess struct {
	Owner   sdk.AccAddress            `json:"address" validate:"required,accaddress"`
	Request *api.CreateProcessRequest `json:"request" validate:"required"`
}

// NewMsgCreateProcess is a constructor function for MsgCreateProcess.
func NewMsgCreateProcess(owner sdk.AccAddress, request *api.CreateProcessRequest) *MsgCreateProcess {
	return &MsgCreateProcess{
		Owner:   owner,
		Request: request,
	}
}

// Route should return the name of the module route.
func (msg MsgCreateProcess) Route() string {
	return RouterKey
}

// Type returns the action.
func (msg MsgCreateProcess) Type() string {
	return "CreateProcess"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgCreateProcess) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner is missing")
	}
	if err := xvalidator.Validate.Struct(msg); err != nil {
		return err
	}
	p := &processpb.Process{
		Name:  msg.Request.Name,
		Nodes: msg.Request.Nodes,
		Edges: msg.Request.Edges,
	}
	// TODO: i don't think we need to generate the hash here
	p.Hash = address.ProcAddress(crypto.AddressHash([]byte(p.HashSerialize())))
	if err := p.Validate(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgCreateProcess) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgCreateProcess) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgDeleteProcess defines a state transition to delete a process.
type MsgDeleteProcess struct {
	Owner   sdk.AccAddress            `json:"address" validate:"required,accaddress"`
	Request *api.DeleteProcessRequest `json:"request" validate:"required"`
}

// NewMsgDeleteProcess is a constructor function for MsgDeleteProcess.
func NewMsgDeleteProcess(owner sdk.AccAddress, request *api.DeleteProcessRequest) *MsgDeleteProcess {
	return &MsgDeleteProcess{
		Owner:   owner,
		Request: request,
	}
}

// Route should return the name of the module.
func (msg MsgDeleteProcess) Route() string {
	return ModuleName
}

// Type returns the action.
func (msg MsgDeleteProcess) Type() string {
	return "DeleteProcess"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgDeleteProcess) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner is missing")
	}
	if err := xvalidator.Validate.Struct(msg); err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgDeleteProcess) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgDeleteProcess) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
