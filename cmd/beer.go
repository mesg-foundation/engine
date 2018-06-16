package cmd

import (
	"fmt"

	"github.com/kyokomi/emoji"
	"github.com/mesg-foundation/core/cmd/utils"

	"github.com/spf13/cobra"
)

// Beer send us token to buy a beer
var Beer = &cobra.Command{
	Use:               "beer",
	Short:             "Pay us a beer",
	DisableAutoGenTag: true,
	Run:               beerHandler,
}

func beerHandler(cmd *cobra.Command, args []string) {
	account := utils.AccountFromFlagOrAsk(cmd, "Select your account")
	amount, err := utils.GetOrAskAmount(cmd, "Select the amount of MESG you want to kindly send")
	if err != nil {
		panic(err)
	}
	if utils.Confirm(cmd, "Are you sure ?") {
		fmt.Println(emoji.Sprint("Thanks you, we will have a nice :beer:"), account, amount)
	}
}

// TODO this command is disabled for now waiting for the transfer feature to be implemented
func init() {
	// RootCmd.AddCommand(Beer)
	// utils.Accountable(Beer)
	// utils.Payable(Beer)
	// utils.Confirmable(Beer)
}
