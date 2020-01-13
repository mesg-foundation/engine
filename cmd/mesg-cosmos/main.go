package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/logger"
	enginesdk "github.com/mesg-foundation/engine/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	db "github.com/tendermint/tm-db"
)

var (
	defaultCLIHome = os.ExpandEnv("$HOME/.mesg-cosmos-cli")
	app            *cosmos.App
)

func main() {
	// init app and codec
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	db, err := db.NewGoLevelDB("app", filepath.Join(cfg.Path, cfg.Cosmos.RelativePath))
	if err != nil {
		panic(err)
	}
	appFactory := cosmos.NewAppFactory(logger.TendermintLogger(), db, cfg.Cosmos.MinGasPrices)
	enginesdk.NewBackend(appFactory)
	app, err = cosmos.NewApp(appFactory)
	if err != nil {
		panic(err)
	}
	cosmos.CustomizeConfig()
	cdc := codec.Codec

	rootCmd := &cobra.Command{
		Use:   "mesg-cosmos",
		Short: "Cosmos Client",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(defaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		lcd.ServeCommand(cdc, registerRoutes),
		keys.Commands(),
		version.Cmd,
		client.NewCompletionCmd(rootCmd, true),
	)

	executor := cli.PrepareMainCmd(rootCmd, "MESG", defaultCLIHome)
	err = executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.BasicManager().RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
	)

	// add modules' query commands
	app.BasicManager().AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
	)

	// add modules' tx commands
	app.BasicManager().AddTxCommands(txCmd, cdc)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

/*
List of commands that could be used in this cli

distribution
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	func GetCmdWithdrawRewards(cdc *codec.Codec) *cobra.Command {
	func GetCmdWithdrawAllRewards(cdc *codec.Codec, queryRoute string) *cobra.Command {
	func GetCmdSetWithdrawAddr(cdc *codec.Codec) *cobra.Command {
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryValidatorOutstandingRewards(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryValidatorCommission(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryValidatorSlashes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryDelegatorRewards(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryCommunityPool(queryRoute string, cdc *codec.Codec) *cobra.Command {
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {

gov
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryProposal(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryProposals(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryVote(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryVotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryDeposit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryDeposits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryTally(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryParam(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryProposer(queryRoute string, cdc *codec.Codec) *cobra.Command {
func GetTxCmd(storeKey string, cdc *codec.Codec, pcmds []*cobra.Command) *cobra.Command {
	func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	func GetCmdVote(cdc *codec.Codec) *cobra.Command {

bank
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	func SendTxCmd(cdc *codec.Codec) *cobra.Command {

params
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {

auth
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	func GetAccountCmd(cdc *codec.Codec) *cobra.Command {
func QueryTxsByEventsCmd(cdc *codec.Codec) *cobra.Command {
func QueryTxCmd(cdc *codec.Codec) *cobra.Command {
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	func GetSignCommand(codec *codec.Codec) *cobra.Command {
	func GetMultiSignCommand(cdc *codec.Codec) *cobra.Command {
func GetBroadcastCommand(cdc *codec.Codec) *cobra.Command {
func GetEncodeCommand(cdc *codec.Codec) *cobra.Command {

staking
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryValidator(storeName string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryValidators(storeName string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryValidatorUnbondingDelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryValidatorRedelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryDelegation(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryDelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryValidatorDelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryUnbondingDelegation(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryUnbondingDelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryRedelegation(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryRedelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryPool(storeName string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	func GetCmdCreateValidator(cdc *codec.Codec) *cobra.Command {
	func GetCmdEditValidator(cdc *codec.Codec) *cobra.Command {
	func GetCmdDelegate(cdc *codec.Codec) *cobra.Command {
	func GetCmdRedelegate(storeName string, cdc *codec.Codec) *cobra.Command {
	func GetCmdUnbond(storeName string, cdc *codec.Codec) *cobra.Command {

crisis
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	func GetCmdInvariantBroken(cdc *codec.Codec) *cobra.Command {

supply
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryTotalSupply(cdc *codec.Codec) *cobra.Command {

slashing
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQuerySigningInfo(storeName string, cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {

mint
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryInflation(cdc *codec.Codec) *cobra.Command {
	func GetCmdQueryAnnualProvisions(cdc *codec.Codec) *cobra.Command {

client
func ConfigCmd(defaultCLIHome string) *cobra.Command {

client/lcd
func ServeCommand(cdc *codec.Codec, registerRoutesFn func(*RestServer)) *cobra.Command {

client/keys
func Commands() *cobra.Command {

client/rpc
func ValidatorCommand(cdc *codec.Codec) *cobra.Command {
func StatusCommand() *cobra.Command {
func BlockCommand() *cobra.Command {
*/
