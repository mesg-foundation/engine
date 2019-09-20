package servicesdk

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// msgCreateService defines a state transition to create a service.
type msgCreateService struct {
	Request *api.CreateServiceRequest `json:"request"`
	// Owner   cosmostypes.AccAddress   `json:"owner"`
	cdc *codec.Codec
}

// newMsgCreateService is a constructor function for msgCreateService.
func newMsgCreateService(cdc *codec.Codec, req *api.CreateServiceRequest) *msgCreateService {
	return &msgCreateService{
		Request: req,
		// Owner:   owner,
		cdc: cdc,
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
	// if msg.Owner.Empty() {
	// 	return cosmostypes.ErrInvalidAddress(msg.Owner.String())
	// }
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg msgCreateService) GetSignBytes() []byte {
	return cosmostypes.MustSortJSON(msg.cdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg msgCreateService) GetSigners() []cosmostypes.AccAddress {
	return []cosmostypes.AccAddress{}
	// return []cosmostypes.AccAddress{msg.Owner}
}

// --------------------------------------------------------
//
// --------------------------------------------------------

// msgDeleteService defines a state transition to delete a service.
type msgDeleteService struct {
	Hash hash.Hash `json:"hash"`
	cdc  *codec.Codec
}

// newMsgDeleteService is a constructor function for msgSetService.
func newMsgDeleteService(cdc *codec.Codec, hash hash.Hash) *msgDeleteService {
	return &msgDeleteService{
		Hash: hash,
		cdc:  cdc,
	}
}

// Route should return the name of the module.
func (msg msgDeleteService) Route() string {
	return backendName
}

// Type returns the action.
func (msg msgDeleteService) Type() string {
	return "delete_service"
}

// ValidateBasic runs stateless checks on the message.
func (msg msgDeleteService) ValidateBasic() cosmostypes.Error {
	// if msg.Owner.Empty() {
	// 	return cosmostypes.ErrInvalidAddress(msg.Owner.String())
	// }
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg msgDeleteService) GetSignBytes() []byte {
	return cosmostypes.MustSortJSON(msg.cdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg msgDeleteService) GetSigners() []cosmostypes.AccAddress {
	return []cosmostypes.AccAddress{}
}
