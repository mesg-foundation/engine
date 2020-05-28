package types

import (
	"fmt"

	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/instance"
)

// GenesisState - all instance state that must be provided at genesis
type GenesisState struct {
	Instances []*instance.Instance `json:"instances" yaml:"instances" validate:"dive"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(instances []*instance.Instance) GenesisState {
	return GenesisState{
		Instances: instances,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Instances: []*instance.Instance{},
	}
}

// ValidateGenesis validates the instance genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := xvalidator.Struct(data); err != nil {
		return fmt.Errorf("failed to validate %s genesis state: %w", ModuleName, err)
	}
	return nil
}
