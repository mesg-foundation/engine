package servicesdk

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/mesg-foundation/engine/hash"
	serv "github.com/mesg-foundation/engine/service"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = Service{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type GenesisState struct {
	Services []serv.Service `json:"services"`
}

// app module Basics object
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return "service"
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	cdc := codec.New()
	return cdc.MustMarshalJSON(GenesisState{
		Services: []serv.Service{
			serv.Service{
				Hash:          hash.Int(1),
				Name:          "genesis-hash",
				Configuration: &serv.Configuration{},
			},
		},
	})
}

// Validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return nil
}

// Register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {}

// Get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

// Get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

func (s Service) Name() string {
	return "service"
}

func (s Service) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (s Service) Route() string {
	return "service"
}

func (s Service) NewHandler() sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		errmsg := fmt.Sprintf("Unrecognized service Msg type: %v", msg.Type())
		return sdk.ErrUnknownRequest(errmsg).Result()
	}
}
func (s Service) QuerierRoute() string {
	return "service"
}

// func (s Service) NewQuerierHandler() sdk.Querier {
// 	return NewQuerier(am.keeper)
// }

func (s Service) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {

}

func (s Service) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (s Service) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	cdc := codec.New()
	cdc.MustUnmarshalJSON(data, &genesisState)
	for _, service := range genesisState.Services {
		s.keeper.Set(ctx, service)
	}
	return []abci.ValidatorUpdate{}
}

func (s Service) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := GenesisState{}
	cdc := codec.New()
	return cdc.MustMarshalJSON(gs)
}
