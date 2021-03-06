package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/x/instance/internal/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group instance queries under a subcommand
	instanceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Query commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	instanceQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdGet(queryRoute, cdc),
			GetCmdList(queryRoute, cdc),
		)...,
	)

	return instanceQueryCmd
}

// GetCmdGet implements the get query command.
func GetCmdGet(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [hash]",
		Short: "Fetch an instance by its hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryGet, args[0]), nil)
			if err != nil {
				fmt.Printf("could not get instance\n%s\n", err.Error())
				return nil
			}

			var out *instance.Instance
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdList implements the list query command.
func GetCmdList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Query all the instances",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryList), nil)
			if err != nil {
				fmt.Printf("could not list instances\n%s\n", err.Error())
				return nil
			}

			var out []*instance.Instance
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
