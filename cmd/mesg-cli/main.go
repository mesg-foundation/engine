package main

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/mesg-foundation/engine/app"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/types"
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()

	// Read in the configuration file for the sdk
	// TODO: change back when more refactor is done
	// config := sdk.GetConfig()
	// config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	// config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	// config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	// config.Seal()
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	cosmos.CustomizeConfig(cfg)

	rootCmd := &cobra.Command{
		Use:   "mesg-cli",
		Short: "Command line interface for interacting with appd",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		exportTxsCmd(cdc),
		importTxsCmd(cdc),
		flags.LineBreak,
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		flags.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCmd, true),
	)

	// Add flags and prefix all env exposed with MESG
	executor := cli.PrepareMainCmd(rootCmd, "MESG", app.DefaultCLIHome)

	err = executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		flags.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		authcmd.GetDecodeCommand(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
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
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

func exportTxsCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exportTxs",
		Short: "Export all txs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			status, err := cliCtx.Client.Status()
			if err != nil {
				return err
			}
			fmt.Println("last block", status.SyncInfo.LatestBlockHeight)

			numJobs := status.SyncInfo.LatestBlockHeight - 1
			jobs := make(chan int64, numJobs)
			jobsDone := make(chan bool, numJobs)
			results := make([]types.Tx, 0)
			resultsMux := sync.Mutex{}
			jobsWithError := make([]int64, 0)
			jobsWithErrorMux := sync.Mutex{}

			worker := func(cliCtx context.CLIContext, id int, jobs <-chan int64, jobsDone chan bool) {
				for height := range jobs {
					fmt.Printf("\nworking %d is fetching txs of block %d", id, height)
					block, err := cliCtx.Client.Block(&height)
					if err != nil {
						fmt.Println("err", err)
						jobsWithErrorMux.Lock()
						jobsWithError = append(jobsWithError, height)
						jobsWithErrorMux.Unlock()
						jobsDone <- true
						continue
					}
					fmt.Printf(". found %d txs", len(block.Block.Txs))
					resultsMux.Lock()
					results = append(results, block.Block.Txs...)
					resultsMux.Unlock()
					jobsDone <- true
				}
			}

			for w := 1; w <= 20; w++ {
				go worker(cliCtx, w, jobs, jobsDone)
			}
			for height := int64(1); height <= status.SyncInfo.LatestBlockHeight; height++ {
				jobs <- height
			}
			close(jobs)

			for a := int64(0); a <= numJobs; a++ {
				<-jobsDone
			}

			json, err := cdc.MarshalJSONIndent(results, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(json))

			fmt.Println("jobsWithError", jobsWithError)
			return nil
		},
	}

	return flags.GetCommands(cmd)[0]
}

func importTxsCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "importTxs",
		Short: "import txs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var txs []types.Tx
			err := cdc.UnmarshalJSON([]byte(args[0]), &txs)
			if err != nil {
				return err
			}

			for _, tx := range txs {
				res, err := cliCtx.BroadcastTxSync(tx)
				if err != nil {
					return err
				}
				cliCtx.PrintOutput(res)
			}
			return nil
		},
	}

	return flags.GetCommands(cmd)[0]
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
