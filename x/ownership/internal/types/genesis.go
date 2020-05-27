package types

import (
	fmt "fmt"

	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/ownership"
)

// GenesisState - all ownership state that must be provided at genesis
type GenesisState struct {
	Ownerships []*ownership.Ownership `json:"ownerships" yaml:"ownerships" validate:"dive"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(ownerships []*ownership.Ownership) GenesisState {
	return GenesisState{
		Ownerships: ownerships,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Ownerships: nil,
	}
}

// ValidateGenesis validates the ownership genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := xvalidator.Struct(data); err != nil {
		return fmt.Errorf("failed to validate %s genesis state: %w", ModuleName, err)
	}
	return nil
}
