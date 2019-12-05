package executionsdk

import (
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/x/xvalidator"
)

// msgCreateExecution defines a state transition to create an execution.
type msgCreateExecution struct {
	Request *api.CreateExecutionRequest `json:"request" validate:"required"`
	Signer  cosmostypes.AccAddress      `json:"signer" validate:"required,accaddress"`
}

const msgCreateExecutionType = "create_execution"

// newMsgCreateExecution is a constructor function for msgCreateExecution.
func newMsgCreateExecution(req *api.CreateExecutionRequest, signer cosmostypes.AccAddress) *msgCreateExecution {
	return &msgCreateExecution{
		Request: req,
		Signer:  signer,
	}
}

// Route should return the name of the module.
func (msg msgCreateExecution) Route() string {
	return backendName
}

// Type returns the action.
func (msg msgCreateExecution) Type() string {
	return msgCreateExecutionType
}

// ValidateBasic runs stateless checks on the message.
func (msg msgCreateExecution) ValidateBasic() cosmostypes.Error {
	if err := xvalidator.Validate.Struct(msg); err != nil {
		return cosmostypes.ErrInternal(err.Error())
	}
	if !msg.Request.ParentResultHash.IsZero() && !msg.Request.EventHash.IsZero() {
		return cosmostypes.ErrInternal("cannot have both parent and event hash")
	}
	if msg.Request.ParentResultHash.IsZero() && msg.Request.EventHash.IsZero() {
		return cosmostypes.ErrInternal("should have at least an event hash or parent hash")
	}
	if msg.Request.ExecutorHash.IsZero() {
		return cosmostypes.ErrInternal("should have a executor hash")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg msgCreateExecution) GetSignBytes() []byte {
	return cosmostypes.MustSortJSON(codec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg msgCreateExecution) GetSigners() []cosmostypes.AccAddress {
	return []cosmostypes.AccAddress{msg.Signer}
}
