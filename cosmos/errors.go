package cosmos

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// CodespaceMesg is a cosmos codespace for all mesg errors.
	CodespaceMesg = "mesg"
)

// Base mesg codes.
const (
	CodeInternal   uint32 = 1000
	CodeValidation uint32 = 2000
	CodeMarshal    uint32 = 2001
	CodeUnmarshal  uint32 = 2002
)

// mesg errors
var (
	ErrValidation = sdkerrors.Register(CodespaceMesg, CodeValidation, "validation failed")
	ErrMarshal    = sdkerrors.Register(CodespaceMesg, CodeMarshal, "failed to marshal")     // TODO: to replace by cosmoserrors.ErrJSONMarshal if it makes sense
	ErrUnmarshal  = sdkerrors.Register(CodespaceMesg, CodeUnmarshal, "failed to unmarshal") // TODO: to replace by cosmoserrors.ErrJSONUnmarshal if it makes sense
)

// NewMesgWrapError creates error with given code type and mesg codespace.
func NewMesgWrapError(code uint32, err error) *sdkerrors.Error {
	return sdkerrors.New(CodespaceMesg, code, err.Error())
}
