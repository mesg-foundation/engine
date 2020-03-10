package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// MsgCreateService defines a state transition to create a service.
type MsgCreateService struct {
	Owner   sdk.AccAddress            `json:"owner" validate:"required,accaddress"`
	Request *api.CreateServiceRequest `json:"request" validate:"required"`
}

// NewMsgCreateService is a constructor function for MsgCreateService.
func NewMsgCreateService(owner sdk.AccAddress, request *api.CreateServiceRequest) *MsgCreateService {
	return &MsgCreateService{
		Owner:   owner,
		Request: request,
	}
}

// Route should return the name of the module route.
func (msg MsgCreateService) Route() string {
	return RouterKey
}

// Type returns the action.
func (msg MsgCreateService) Type() string {
	return "CreateService"
}

// ValidateBasic runs stateless checks on the message.
func (msg MsgCreateService) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner is missing")
	}
	if err := xvalidator.Validate.Struct(msg); err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg MsgCreateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg MsgCreateService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
