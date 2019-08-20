package tendermint

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

func NewAppModuleBasic(name string) AppModuleBasic {
	return AppModuleBasic{
		name: name,
	}
}

func NewAppModule(moduleBasic AppModuleBasic, handler sdk.Handler, querier sdk.Querier) AppModule {
	return AppModule{
		AppModuleBasic: moduleBasic,
		handler:        handler,
		querier:        querier,
	}
}

type AppModuleBasic struct {
	name string
}

func (m AppModuleBasic) Name() string {
	return m.name
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return []byte("{}")
}

// ValidateGenesis checks the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {}

// GetQueryCmd returns the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

// GetTxCmd returns the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

type AppModule struct {
	AppModuleBasic
	handler sdk.Handler
	querier sdk.Querier
}

func (m AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (m AppModule) Route() string {
	return m.name
}

func (m AppModule) NewHandler() sdk.Handler {
	return m.handler
	// return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// errmsg := fmt.Sprintf("Unrecognized service Msg type: %v", msg.Type())
	// return sdk.ErrUnknownRequest(errmsg).Result()
	// }
}
func (m AppModule) QuerierRoute() string {
	return m.name
}

func (m AppModule) NewQuerierHandler() sdk.Querier {
	return m.querier
}

func (m AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (m AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (m AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (m AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return []byte("{}")
}
