package types

import (
	fmt "fmt"

	executionpb "github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xvalidator"
)

// GenesisState - all instance state that must be provided at genesis
type GenesisState struct {
	Executions []*executionpb.Execution `json:"executions" yaml:"executions" validate:"dive"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(execs []*executionpb.Execution) GenesisState {
	return GenesisState{
		Executions: execs,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Executions: []*executionpb.Execution{},
	}
}

// ValidateGenesis validates the instance genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := xvalidator.Struct(data); err != nil {
		return fmt.Errorf("failed to validate %s genesis state: %w", ModuleName, err)
	}
	return nil
}
