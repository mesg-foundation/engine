package servicesdk

import (
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	cosmoserrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/validator"
)

// msgCreateService defines a state transition to create a service.
type msgCreateService struct {
	Request *api.CreateServiceRequest `json:"request" validate:"required"`
	Owner   cosmostypes.AccAddress    `json:"owner" validate:"required,accaddress"`
}

// newMsgCreateService is a constructor function for msgCreateService.
func newMsgCreateService(req *api.CreateServiceRequest, owner cosmostypes.AccAddress) *msgCreateService {
	return &msgCreateService{
		Request: req,
		Owner:   owner,
	}
}

// Route should return the name of the module.
func (msg msgCreateService) Route() string {
	return ModuleName
}

// Type returns the action.
func (msg msgCreateService) Type() string {
	return "create_service"
}

// ValidateBasic runs stateless checks on the message.
func (msg msgCreateService) ValidateBasic() error {
	if err := validator.Validate.Struct(msg); err != nil {
		return err
	}
	if msg.Owner.Empty() {
		return cosmoserrors.Wrap(cosmoserrors.ErrInvalidAddress, "owner is missing")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg msgCreateService) GetSignBytes() []byte {
	return cosmostypes.MustSortJSON(codec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg msgCreateService) GetSigners() []cosmostypes.AccAddress {
	return []cosmostypes.AccAddress{msg.Owner}
}
