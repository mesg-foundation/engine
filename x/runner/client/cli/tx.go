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
		GetCmdCreate(cdc),
		GetCmdDelete(cdc),
	)...)

	return runnerTxCmd
}

func GetCmdCreate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [runnerHash] [key=val]...",
		Short: "Creates a new runner",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO:
			return nil
		},
	}
}

func GetCmdDelete(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [runnerHash] [key=val]...",
		Short: "Creates a new runner",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO:
			return nil
		},
	}
}
