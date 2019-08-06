package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Service is a struct that contains all the metadata of a service.
type Service struct {
	Hash       string         `json:"hash"`
	Definition string         `json:"definition"`
	Owner      sdk.AccAddress `json:"owner"`
}
