package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/x/process/internal/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group process queries under a subcommand
	processQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	processQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdGet(queryRoute, cdc),
			GetCmdList(queryRoute, cdc),
		)...,
	)

	return processQueryCmd
}

func GetCmdGet(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "get",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryGet, args[0]), nil)
			if err != nil {
				fmt.Printf("could not get process\n%s\n", err.Error())
				return nil
			}

			var out *process.Process
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryList), nil)
			if err != nil {
				fmt.Printf("could not list processes\n%s\n", err.Error())
				return nil
			}

			var out []*process.Process
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
