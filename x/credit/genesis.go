package credit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/x/credit/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize the keeper with the data from the genesis file.
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	if err := k.Import(ctx, data.Credits); err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k Keeper) (data types.GenesisState) {
	execs, err := k.Export(ctx)
	if err != nil {
		panic(err)
	}
	return types.NewGenesisState(execs)
}
