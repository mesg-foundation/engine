package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	executionTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", strings.Title(types.ModuleName)),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	executionTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreate(cdc),
		GetCmdUpdate(cdc),
	)...)
	return executionTxCmd
}

// GetCmdCreate is the command to create a execution.
func GetCmdCreate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [definition]",
		Short: "Creates a new execution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			rawMsg := fmt.Sprintf(`{"type":"execution/create","value":%s}`, args[0])
			var msg types.MsgCreate
			if err := cdc.UnmarshalJSON([]byte(rawMsg), &msg); err != nil {
				return err
			}
			if !msg.Signer.Empty() && !msg.Signer.Equals(cliCtx.FromAddress) {
				return fmt.Errorf("the signer set in the definition is not the same as the from flag")
			}
			msg.Signer = cliCtx.FromAddress
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUpdate is the command to update a execution.
func GetCmdUpdate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update [definition]",
		Short: "Updates an execution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			rawMsg := fmt.Sprintf(`{"type":"execution/update","value":%s}`, args[0])
			var msg types.MsgUpdate
			if err := cdc.UnmarshalJSON([]byte(rawMsg), &msg); err != nil {
				return err
			}
			if !msg.Executor.Empty() && !msg.Executor.Equals(cliCtx.FromAddress) {
				return fmt.Errorf("the executor set in the definition is not the same as the from flag")
			}
			msg.Executor = cliCtx.FromAddress
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
