package types

import (
	fmt "fmt"

	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/process"
)

// GenesisState - all instance state that must be provided at genesis
type GenesisState struct {
	Processes []*process.Process `json:"processes" yaml:"processes" validate:"dive"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(processes []*process.Process) GenesisState {
	return GenesisState{
		Processes: processes,
	}
}

// DefaultGenesisState is the default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Processes: []*process.Process{},
	}
}

// ValidateGenesis validates the instance genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := xvalidator.Struct(data); err != nil {
		return fmt.Errorf("failed to validate %s genesis state: %w", ModuleName, err)
	}
	return nil
}
