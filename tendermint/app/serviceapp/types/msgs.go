package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName

// MsgSetService defines a SetService message.
type MsgSetService struct {
	Hash       string         `json:"hash"`
	Definition string         `json:"definition"`
	Owner      sdk.AccAddress `json:"owner"`
}

// NewMsgSetService is a constructor function for MsgSetService.
func NewMsgSetService(hash, definition string, owner sdk.AccAddress) MsgSetService {
	return MsgSetService{
		Hash:       hash,
		Definition: definition,
		Owner:      owner,
	}
}

// Route should return the name of the module.
func (msg MsgSetService) Route() string {
	return RouterKey
}

// Type returns the action.
func (msg MsgSetService) Type() string {
	return "set_service"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgSetService) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if msg.Hash == "" {
		return sdk.ErrUnknownRequest("hash cannot be empty")
	}
	if msg.Definition == "" {
		return sdk.ErrUnknownRequest("definition cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgSetService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgSetService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgRemoveService defines the RemoveService message.
type MsgRemoveService struct {
	Hash    string         `json:"hash"`
	Remover sdk.AccAddress `json:"buyer"`
}

// NewMsgRemoveService is the constructor function for MsgRemoveService.
func NewMsgRemoveService(hash string, remover sdk.AccAddress) MsgRemoveService {
	return MsgRemoveService{
		Hash:    hash,
		Remover: remover,
	}
}

// Route should return the name of the module.
func (msg MsgRemoveService) Route() string {
	return RouterKey
}

// Type returns the action.
func (msg MsgRemoveService) Type() string {
	return "remove_service"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgRemoveService) ValidateBasic() sdk.Error {
	if msg.Remover.Empty() {
		return sdk.ErrInvalidAddress(msg.Remover.String())
	}
	if msg.Hash == "" {
		return sdk.ErrUnknownRequest("hash cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgRemoveService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgRemoveService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Remover}
}
