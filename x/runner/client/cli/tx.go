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
	"github.com/mesg-foundation/engine/x/runner/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	runnerTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", strings.Title(types.ModuleName)),
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

// GetCmdCreate is the command to create a runner.
func GetCmdCreate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [serviceHash] [envHash]",
		Short: "Create a new runner",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			serviceHash, err := hash.Decode(args[0])
			if err != nil {
				return fmt.Errorf("arg serviceHash error: %w", err)
			}

			var envHash hash.Hash
			if len(args) >= 2 && args[1] != "" {
				if envHash, err = hash.Decode(args[1]); err != nil {
					return fmt.Errorf("arg envHash error: %w", err)
				}
			}

			if cliCtx.FromAddress.Empty() {
				return fmt.Errorf("flag --from is required")
			}

			msg := types.MsgCreate{
				ServiceHash: serviceHash,
				EnvHash:     envHash,
				Owner:       cliCtx.FromAddress,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDelete is the command to delete a runner.
func GetCmdDelete(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [runnerHash]",
		Short: "Delete a runner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			runHash, err := hash.Decode(args[0])
			if err != nil {
				return fmt.Errorf("arg runHash error: %w", err)
			}

			msg := types.MsgDelete{
				Hash:  runHash,
				Owner: cliCtx.FromAddress,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
