package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/ext/xvalidator"
)

// GenesisState - all instance state that must be provided at genesis
type GenesisState struct {
	Credits map[string]sdk.Int `json:"credits" yaml:"credits" validate:"dive"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(credits map[string]sdk.Int) GenesisState {
	return GenesisState{
		Credits: credits,
	}
}

// DefaultGenesisState is the default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
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
