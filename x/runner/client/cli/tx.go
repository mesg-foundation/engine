package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/x/runner/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	runnerTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	runnerTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateRunner(cdc),
		GetCmdDeleteRunner(cdc),
	)...)

	return runnerTxCmd
}

func GetCmdCreateRunner(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "createRunner [serviceHash] [key=val]...",
		Short: "Creates a new runner",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO:
			return nil
		},
	}
}

func GetCmdDeleteRunner(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "createRunner [serviceHash] [key=val]...",
		Short: "Creates a new runner",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO:
			return nil
		},
	}
}
