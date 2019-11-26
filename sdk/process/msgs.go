package processsdk

import (
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/x/xvalidator"
)

// msgCreateProcess defines a state transition to create a service.
type msgCreateProcess struct {
	Owner   cosmostypes.AccAddress    `json:"owner" validate:"required,accaddress"`
	Request *api.CreateProcessRequest `json:"request" validate:"required"`
}

// newMsgCreateProcess is a constructor function for msgCreateProcess.
func newMsgCreateProcess(owner cosmostypes.AccAddress, req *api.CreateProcessRequest) *msgCreateProcess {
	return &msgCreateProcess{
		Request: req,
		Owner:   owner,
	}
}

// Route should return the name of the module.
func (msg msgCreateProcess) Route() string {
	return backendName
}

// Type returns the action.
func (msg msgCreateProcess) Type() string {
	return "create_process"
}

// ValidateBasic runs stateless checks on the message.
func (msg msgCreateProcess) ValidateBasic() cosmostypes.Error {
	if msg.Owner.Empty() {
		return cosmostypes.ErrInvalidAddress("owner is missing")
	}
	if err := xvalidator.Validate.Struct(msg); err != nil {
		return cosmostypes.ErrInternal(err.Error())
	}
	p := &process.Process{
		Key:   msg.Request.Key,
		Nodes: msg.Request.Nodes,
		Edges: msg.Request.Edges,
	}
	if err := p.Validate(); err != nil {
		return cosmostypes.ErrInternal(err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg msgCreateProcess) GetSignBytes() []byte {
	return cosmostypes.MustSortJSON(codec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg msgCreateProcess) GetSigners() []cosmostypes.AccAddress {
	return []cosmostypes.AccAddress{msg.Owner}
}

// msgDeleteProcess defines a state transition to create a service.
type msgDeleteProcess struct {
	Owner   cosmostypes.AccAddress    `json:"owner" validate:"required,accaddress"`
	Request *api.DeleteProcessRequest `json:"request" validate:"required"`
}

// newMsgDeleteProcess is a constructor function for msgDeleteProcess.
func newMsgDeleteProcess(owner cosmostypes.AccAddress, request *api.DeleteProcessRequest) *msgDeleteProcess {
	return &msgDeleteProcess{
		Owner:   owner,
		Request: request,
	}
}

// Route should return the name of the module.
func (msg msgDeleteProcess) Route() string {
	return backendName
}

// Type returns the action.
func (msg msgDeleteProcess) Type() string {
	return "delete_process"
}

// ValidateBasic runs stateless checks on the message.
func (msg msgDeleteProcess) ValidateBasic() cosmostypes.Error {
	if msg.Owner.Empty() {
		return cosmostypes.ErrInvalidAddress("owner is missing")
	}
	if err := xvalidator.Validate.Struct(msg); err != nil {
		return cosmostypes.ErrInternal(err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg msgDeleteProcess) GetSignBytes() []byte {
	return cosmostypes.MustSortJSON(codec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg msgDeleteProcess) GetSigners() []cosmostypes.AccAddress {
	return []cosmostypes.AccAddress{msg.Owner}
}
