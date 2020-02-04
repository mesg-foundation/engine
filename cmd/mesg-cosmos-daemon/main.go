package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/cosmos"
	enginesdk "github.com/mesg-foundation/engine/sdk"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
)

var (
	defaultCLIHome = os.ExpandEnv("$HOME/.mesg-cosmos-cli")
	basicManager   module.BasicManager
)

func main() {
	cobra.EnableCommandSorting = false

	// init app and codec
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	cosmos.CustomizeConfig(cfg)

	basicManager = *enginesdk.NewBasicManager()
	cdc := codec.Codec

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "mesg-cosmos-daemon",
		Short:             "ClI Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, basicManager, cfg.Tendermint.Config.RootDir),
		genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, cfg.Tendermint.Config.RootDir),
		genutilcli.GenTxCmd(
			ctx, cdc, basicManager, staking.AppModuleBasic{},
			auth.GenesisAccountIterator{}, cfg.Tendermint.Config.RootDir, defaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, basicManager),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "MESG", cfg.Tendermint.Config.RootDir)
	err = executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db db.DB, traceStore io.Writer) abci.Application {
	// TODO: could start the engine here
	return nil
}

func exportAppStateAndTMValidators(logger log.Logger, db db.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	return nil, nil, nil
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
