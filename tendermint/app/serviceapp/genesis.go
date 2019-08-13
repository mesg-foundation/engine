package serviceapp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	Services []Service `json:"services"`
}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(gs GenesisState) error {
	for _, service := range gs.Services {
		if service.Owner == nil {
			return fmt.Errorf("invalid service genesis: missing owner")
		}
		if service.Definition == "" {
			return fmt.Errorf("invalid service genesis: missing definition")
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Services: []Service{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, gs GenesisState) []abci.ValidatorUpdate {
	for _, service := range gs.Services {
		keeper.SetService(ctx, service)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	var services []Service
	for it := keeper.GetHashesIterator(ctx); it.Valid(); it.Next() {
		hash := string(it.Key())
		services = append(services, keeper.GetService(ctx, hash))
	}
	return GenesisState{Services: services}
}
