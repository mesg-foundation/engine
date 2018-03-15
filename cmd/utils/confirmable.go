package cmdUtils

import (
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Confirm checks that the flag "confirm" is set and otherwise ask a confirmation in the command line
func Confirm(cmd *cobra.Command, message string) bool {
	confirm := cmd.Flag("confirm").Value.String() == "true"
	if !confirm {
		survey.AskOne(&survey.Confirm{Message: message}, &confirm, nil)
	}
	return confirm
}

// Confirmable mark a command as confirmable so will have the --confirm flag
func Confirmable(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP("confirm", "c", false, "Confirm")
	return cmd
}
