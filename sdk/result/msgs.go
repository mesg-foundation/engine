package resultsdk

import (
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/x/xvalidator"
)

// msgCreateResult defines a state transition to create an result.
type msgCreateResult struct {
	Request  *api.CreateResultRequest `json:"request" validate:"required"`
	Executor cosmostypes.AccAddress   `json:"executor" validate:"required,accaddress"`
}

const msgCreateResultType = "create_execution"

// newMsgCreateResult is a constructor function for msgCreateResult.
func newMsgCreateResult(req *api.CreateResultRequest, executor cosmostypes.AccAddress) *msgCreateResult {
	return &msgCreateResult{
		Request:  req,
		Executor: executor,
	}
}

// Route should return the name of the module.
func (msg msgCreateResult) Route() string {
	return backendName
}

// Type returns the action.
func (msg msgCreateResult) Type() string {
	return msgCreateResultType
}

// ValidateBasic runs stateless checks on the message.
func (msg msgCreateResult) ValidateBasic() cosmostypes.Error {
	if err := xvalidator.Validate.Struct(msg); err != nil {
		return cosmostypes.ErrInternal(err.Error())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (msg msgCreateResult) GetSignBytes() []byte {
	return cosmostypes.MustSortJSON(codec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required.
func (msg msgCreateResult) GetSigners() []cosmostypes.AccAddress {
	return []cosmostypes.AccAddress{msg.Executor}
}
