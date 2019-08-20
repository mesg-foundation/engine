package servicesdk

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

type appModule struct {
	appModuleBasic

	logic *logic
}

func newAppModule(c container.Container, keeperFactory KeeperFactor) *appModule {
	return &appModule{
		logic: newLogic(c, keeperFactory),
	}
}

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = appModule{}
	_ module.AppModuleBasic = appModuleBasic{}
)

type genesisState struct {
	Services []*service.Service `json:"services"`
}

// app module Basics object
type appModuleBasic struct{}

func (appModuleBasic) Name() string {
	return "service"
}

func (appModuleBasic) RegisterCodec(cdc *codec.Codec) {
}

func (appModuleBasic) DefaultGenesis() json.RawMessage {
	cdc := codec.New()
	return cdc.MustMarshalJSON(genesisState{
		Services: []*service.Service{
			&service.Service{
				Hash:          hash.Int(1),
				Name:          "genesis-hash",
				Configuration: &service.Configuration{},
			},
		},
	})
}

// Validation check of the Genesis
func (appModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return nil
}

// Register rest routes
func (appModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {}

// Get the root query command of this module
func (appModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

// Get the root tx command of this module
func (appModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

func (s appModule) Name() string {
	return "service"
}

func (s appModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (s appModule) Route() string {
	return "service"
}

func (s appModule) NewHandler() sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		errmsg := fmt.Sprintf("Unrecognized service Msg type: %v", msg.Type())
		return sdk.ErrUnknownRequest(errmsg).Result()
	}
}
func (s appModule) QuerierRoute() string {
	return "service"
}

// NewQuerier is the module level router for state queries.
func (s appModule) NewQuerierHandler() sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case "get":
			hash, err := hash.Decode(path[1])
			if err != nil {
				return nil, sdk.ErrInternal(err.Error())
			}
			service, err := s.logic.get(ctx, hash)
			if err != nil {
				return nil, sdk.ErrInternal(err.Error())
			}
			cdc := codec.New()
			res, err := cdc.MarshalJSON(service)
			if err != nil {
				return nil, sdk.ErrInternal(err.Error())
			}
			return res, nil
		default:
			return nil, sdk.ErrUnknownRequest("unknown service query endpoint" + path[0])
		}
	}
}

func (s appModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (s appModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (s appModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState genesisState
	cdc := codec.New()
	cdc.MustUnmarshalJSON(data, &genesisState)
	for _, service := range genesisState.Services {
		s.logic.create(ctx, service)
	}
	return []abci.ValidatorUpdate{}
}

func (s appModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := genesisState{}
	cdc := codec.New()
	return cdc.MustMarshalJSON(gs)
}
