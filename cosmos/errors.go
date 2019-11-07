package cosmos

import (
	"github.com/cosmos/cosmos-sdk/types"
)

const (
	// CodespaceMesg is a cosmos codespace for all mesg errors.
	CodespaceMesg types.CodespaceType = "mesg"
)

const (
	// Base mesg codes.
	CodeMesgInternal   types.CodeType = 1000
	CodeMesgValidation                = 2000
)

// NewMesgErrorf creates error with given code type and mesg codespace.
func NewMesgErrorf(ct types.CodeType, format string, a ...interface{}) types.Error {
	return types.NewError(CodespaceMesg, ct, format, a...)
}

// NewMesgWrapError creates error with given code type and mesg codespace.
func NewMesgWrapError(ct types.CodeType, err error) types.Error {
	return types.NewError(CodespaceMesg, ct, err.Error())
}
