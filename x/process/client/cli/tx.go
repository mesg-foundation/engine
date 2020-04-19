package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/x/process/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	processTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	processTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreate(cdc),
		GetCmdDelete(cdc),
	)...)

	return processTxCmd
}

func GetCmdCreate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [serviceHash] [key=val]...",
		Short: "Creates a new process",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO:
			return nil
		},
	}
}

func GetCmdDelete(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [serviceHash] [key=val]...",
		Short: "Creates a new process",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO:
			return nil
		},
	}
}
