package types

import (
	fmt "fmt"

	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/service"
)

// GenesisState - all instance state that must be provided at genesis
type GenesisState struct {
	Services []*service.Service `json:"services" yaml:"services" validate:"dive"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(services []*service.Service) GenesisState {
	return GenesisState{
		Services: services,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Services: nil,
	}
}

// ValidateGenesis validates the instance genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := xvalidator.Struct(data); err != nil {
		return fmt.Errorf("failed to validate %s genesis state: %w", ModuleName, err)
	}
	return nil
}
