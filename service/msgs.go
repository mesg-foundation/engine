package service

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/types"
)

// MsgSetService ...
type MsgSetService struct {
	Hash    string
	Service Service
	Owner   types.AccAddress
}

// NewMsgSetService is a constructor function for MsgSetService
func NewMsgSetService(hash string, service Service, owner types.AccAddress) MsgSetService {
	return MsgSetService{
		Hash:    hash,
		Service: service,
		Owner:   owner,
	}
}

// Route should return the name of the module
func (msg MsgSetService) Route() string { return "service" }

// Type should return the action
func (msg MsgSetService) Type() string { return "set_service" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetService) ValidateBasic() types.Error {
	if msg.Owner.Empty() {
		return types.ErrInvalidAddress(msg.Owner.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetService) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return types.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSetService) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.Owner}
}
