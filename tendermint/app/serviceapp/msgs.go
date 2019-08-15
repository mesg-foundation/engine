package serviceapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/service/validator"
)

const RouterKey = ModuleName

// MsgSetService defines a SetService message.
type MsgSetService struct {
	Service *service.Service
	Owner   sdk.AccAddress `json:"owner"`
}

// NewMsgSetService is a constructor function for MsgSetService.
func NewMsgSetService(servcie *service.Service, owner sdk.AccAddress) MsgSetService {
	return MsgSetService{
		Service: servcie,
		Owner:   owner,
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

	if err := validator.ValidateService(msg.Service); err != nil {
		return sdk.ErrUnknownRequest(err.Error())
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
	Hash    hash.Hash      `json:"hash"`
	Remover sdk.AccAddress `json:"buyer"`
}

// NewMsgRemoveService is the constructor function for MsgRemoveService.
func NewMsgRemoveService(hash hash.Hash, remover sdk.AccAddress) MsgRemoveService {
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
	if msg.Hash.IsZero() {
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
