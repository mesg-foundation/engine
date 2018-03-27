package cmdUtils

import (
	"os"

	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Confirm checks that the flag "confirm" is set and otherwise ask a confirmation in the command line
func Confirm(cmd *cobra.Command, message string) (confirmed bool) {
	confirmed = cmd.Flag("confirm").Value.String() == "true"
	if !confirmed && survey.AskOne(&survey.Confirm{Message: message}, &confirmed, nil) != nil {
		os.Exit(0)
	}
	return
}

// Confirmable marks a command as confirmable so will have the --confirm flag
func Confirmable(cmd *cobra.Command) {
	cmd.Flags().BoolP("confirm", "c", false, "Confirm")
}
