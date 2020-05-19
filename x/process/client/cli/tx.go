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
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/x/process/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	processTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", strings.Title(types.ModuleName)),
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

// GetCmdCreate is the command to create a process.
func GetCmdCreate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [definition]",
		Short: "Create a new process from its definition",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			rawMsg := fmt.Sprintf(`{"type":"process/create","value":%s}`, args[0])
			var msg types.MsgCreate
			if err := cdc.UnmarshalJSON([]byte(rawMsg), &msg); err != nil {
				return err
			}
			if !msg.Owner.Empty() && !msg.Owner.Equals(cliCtx.FromAddress) {
				return fmt.Errorf("the owner set in the definition is not the same as the from flag")
			}
			msg.Owner = cliCtx.FromAddress
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDelete is the command to delete a process.
func GetCmdDelete(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [processHash]",
		Short: "Delete a process",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			processHash, err := hash.Decode(args[0])
			if err != nil {
				return fmt.Errorf("arg processHash error: %w", err)
			}

			msg := types.MsgDelete{
				Hash:  processHash,
				Owner: cliCtx.FromAddress,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
