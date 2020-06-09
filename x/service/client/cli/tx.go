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
	"github.com/mesg-foundation/engine/x/service/internal/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	serviceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", strings.Title(types.ModuleName)),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	serviceTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreate(cdc),
	)...)

	return serviceTxCmd
}

// GetCmdCreate is the CLI command for creating a service.
func GetCmdCreate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [definition]",
		Short: "Create a service from its definition",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			if cliCtx.FromAddress.Empty() {
				return fmt.Errorf("flag --from is required")
			}

			rawMsg := fmt.Sprintf(`{"type":"service/create","value":%s}`, args[0])
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
