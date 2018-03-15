package cmdUtils

import (
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func convert(value string) int {
	// TODO convert the string into MESG
	return 42
}

// GetOrAskAmount return the amount in MESG based on the flag or the user input
func GetOrAskAmount(cmd *cobra.Command, message string) int {
	amount := cmd.Flag("amount").Value.String()
	if amount == "" {
		survey.AskOne(&survey.Input{Message: message}, &amount, nil)
	}
	return convert(amount)
}

// Payable mark a command as payable
func Payable(cmd *cobra.Command) {
	cmd.Flags().StringP("amount", "", "", "The amount of MESG")
}
