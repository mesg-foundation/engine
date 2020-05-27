package types

import (
	fmt "fmt"

	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/runner"
)

// GenesisState - all instance state that must be provided at genesis
type GenesisState struct {
	Runners []*runner.Runner `json:"runners" yaml:"runners" validate:"dive"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(execs []*runner.Runner) GenesisState {
	return GenesisState{
		Runners: execs,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Runners: nil,
	}
}

// ValidateGenesis validates the instance genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := xvalidator.Struct(data); err != nil {
		return fmt.Errorf("failed to validate %s genesis state: %w", ModuleName, err)
	}
	return nil
}
