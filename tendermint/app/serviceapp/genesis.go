package serviceapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	pbtypes "github.com/mesg-foundation/engine/protobuf/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	Services []*pbtypes.Service `json:"services"`
}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(gs GenesisState) error {
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Services: []*pbtypes.Service{},
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
