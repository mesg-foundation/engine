package servicesdk

import (
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// msgCreateService defines a state transition to create a service.
type msgCreateService struct {
	Request *api.CreateServiceRequest `json:"request"`
	Owner   cosmostypes.AccAddress    `json:"owner"`
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
	return backendName
}

// Type returns the action.
func (msg msgCreateService) Type() string {
	return "create_service"
}

// ValidateBasic runs stateless checks on the message.
func (msg msgCreateService) ValidateBasic() cosmostypes.Error {
	if msg.Owner.Empty() {
		return cosmostypes.ErrInvalidAddress("owner is missing")
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
