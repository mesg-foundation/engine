package errors

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// CodespaceMesg is a cosmos codespace for all mesg errors.
	CodespaceMesg = "mesg"
)

// Base mesg codes.
const (
	CodeValidation uint32 = 2000
	MissingHash    uint32 = 2001
	Unknown        uint32 = 3000
)

// mesg errors
var (
	ErrMissingHash = sdkerrors.Register(CodespaceMesg, MissingHash, "bad request: missing hash")

	ErrValidation = sdkerrors.Register(CodespaceMesg, CodeValidation, "validation failed")

	ErrUnknown = sdkerrors.Register(CodespaceMesg, Unknown, "unknown error")
)

// NewMesgWrapError creates error with given code type and mesg codespace.
func NewMesgWrapError(code uint32, err error) *sdkerrors.Error {
	return sdkerrors.New(CodespaceMesg, code, err.Error())
}
