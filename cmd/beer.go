package cmd

import (
	"fmt"

	"github.com/kyokomi/emoji"
	"github.com/mesg-foundation/application/cmd/utils"

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
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select your account")
	amount, err := cmdUtils.GetOrAskAmount(cmd, "Select the amount of MESG you want to kindly send")
	if err != nil {
		panic(err)
	}
	if cmdUtils.Confirm(cmd, "Are you sure ?") {
		fmt.Println(emoji.Sprint("Thanks you, we will have a nice :beer:"), account, amount)
	}
}

func init() {
	RootCmd.AddCommand(Beer)
	cmdUtils.Accountable(Beer)
	cmdUtils.Payable(Beer)
	cmdUtils.Confirmable(Beer)
}
