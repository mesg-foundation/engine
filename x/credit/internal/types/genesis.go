package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/ext/xvalidator"
)

// GenesisState - all instance state that must be provided at genesis
type GenesisState struct {
	Params  Params             `json:"params" yaml:"params" validate:"dive"`
	Credits map[string]sdk.Int `json:"credits" yaml:"credits" validate:"dive,required,bigint"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, credits map[string]sdk.Int) GenesisState {
	return GenesisState{
		Params:  params,
		Credits: credits,
	}
}

// DefaultGenesisState is the default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:  DefaultParams(),
		Credits: map[string]sdk.Int{},
	}
}

// ValidateGenesis validates the instance genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := xvalidator.Struct(data); err != nil {
		return fmt.Errorf("failed to validate %s genesis state: %w", ModuleName, err)
	}
	return nil
}
