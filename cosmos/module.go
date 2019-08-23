package cosmos

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

// AppModuleBasic is a basic element of an cosmos app.
type AppModuleBasic struct {
	name string
}

// AppModule is a main element of an cosmos app.
type AppModule struct {
	AppModuleBasic
	handler sdk.Handler
	querier Querier
	cdc     *codec.Codec
}

// Querier is responsible to answer to ABCI queries.
type Querier func(request sdk.Request, path []string, req abci.RequestQuery) (res interface{}, err error)

// NewAppModuleBasic inits an AppModuleBasic using a name.
func NewAppModuleBasic(name string) AppModuleBasic {
	return AppModuleBasic{
		name: name,
	}
}

// NewAppModule inits an AppModule using an AppModuleBasic, Codec, Handler and Querier.
func NewAppModule(moduleBasic AppModuleBasic, cdc *codec.Codec, handler sdk.Handler, querier Querier) AppModule {
	return AppModule{
		AppModuleBasic: moduleBasic,
		handler:        handler,
		querier:        querier,
		cdc:            cdc,
	}
}

// ----------------------------------------------
// AppModuleBasic
// ----------------------------------------------

// Name returns the name of the module.
func (m AppModuleBasic) Name() string {
	return m.name
}

// RegisterCodec registers the module's structs in the codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
}

// DefaultGenesis returns the default genesis of the module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return []byte("{}")
}

// ValidateGenesis checks a Genesis.
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

// ----------------------------------------------
// AppModule
// ----------------------------------------------

// RegisterInvariants registers invariants to the registry.
func (m AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route returns the route prefix for transaction of the module.
func (m AppModule) Route() string {
	return m.name
}

// NewHandler returns the handler used to apply transactions.
func (m AppModule) NewHandler() sdk.Handler {
	return m.handler
}

// QuerierRoute the route prefix for query of the module.
func (m AppModule) QuerierRoute() string {
	return m.name
}

// NewQuerierHandler returns the handler used to reply ABCI query.
func (m AppModule) NewQuerierHandler() sdk.Querier {
	return func(request sdk.Request, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		data, err := m.querier(request, path, req)
		if err != nil {
			return nil, sdk.ErrInternal(err.Error())
		}
		res, err := m.cdc.MarshalJSON(data)
		if err != nil {
			return nil, sdk.ErrInternal(err.Error())
		}
		return res, nil
	}
}

// BeginBlock is called at the beginning of the process of a new block.
func (m AppModule) BeginBlock(_ sdk.Request, _ abci.RequestBeginBlock) {}

// EndBlock is called at the end of the process of a new block.
func (m AppModule) EndBlock(sdk.Request, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// InitGenesis initializes the genesis from a request and data.
func (m AppModule) InitGenesis(request sdk.Request, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the current state of the app.
func (m AppModule) ExportGenesis(request sdk.Request) json.RawMessage {
	return []byte("{}")
}
