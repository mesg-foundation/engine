package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/mesg-foundation/engine/app"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cdc := app.MakeCodec()

	// init the config of cosmos
	cosmos.InitConfig()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               version.ServerName,
		Short:             "Engine Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.MigrateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
			auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
	)
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(flags.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(debug.Cmd(cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "MESG", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	if err := executor.Execute(); err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	var cache sdk.MultiStorePersistentCache

	if viper.GetBool(server.FlagInterBlockCache) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	initApp, err := app.NewInitApp(
		logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
		baseapp.SetInterBlockCache(cache),
	)
	if err != nil {
		panic(err)
	}
	return initApp
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		aApp, err := app.NewInitApp(logger, db, traceStore, false, uint(1))
		if err != nil {
			return nil, nil, err
		}
		if err := aApp.LoadHeight(height); err != nil {
			return nil, nil, err
		}
		return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	aApp, err := app.NewInitApp(logger, db, traceStore, true, uint(1))
	if err != nil {
		return nil, nil, err
	}
	return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

/*
genaccounts
func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string) *cobra.Command {

genutil
func CollectGenTxsCmd(ctx *server.Context, cdc *codec.Codec, genAccIterator types.GenesisAccountsIterator, defaultNodeHome string) *cobra.Command {
func GenTxCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager, smbh StakingMsgBuildingHelpers, genAccIterator types.GenesisAccountsIterator, defaultNodeHome, defaultCLIHome string) *cobra.Command {
func InitCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
func ValidateGenesisCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager) *cobra.Command {

server
func AddCommands(ctx *Context, cdc *codec.Codec, rootCmd *cobra.Command, appCreator AppCreator, appExport AppExporter) {
	func StartCmd(ctx *Context, appCreator AppCreator) *cobra.Command {
	func ExportCmd(ctx *Context, cdc *codec.Codec, appExporter AppExporter) *cobra.Command {
	func ShowNodeIDCmd(ctx *Context) *cobra.Command {
	func ShowValidatorCmd(ctx *Context) *cobra.Command {
	func ShowAddressCmd(ctx *Context) *cobra.Command {
	func VersionCmd(ctx *Context) *cobra.Command {
	func UnsafeResetAllCmd(ctx *Context) *cobra.Command {
*/
