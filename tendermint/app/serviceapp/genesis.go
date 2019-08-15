package serviceapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/service/validator"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	Services []*service.Service `json:"services"`
}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(gs GenesisState) error {
	for _, service := range gs.Services {
		if err := validator.ValidateService(service); err != nil {
			return err
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Services: []*service.Service{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, gs GenesisState) []abci.ValidatorUpdate {
	for _, service := range gs.Services {
		keeper.SetService(ctx, service)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{Services: keeper.GetServices(ctx)}
}
