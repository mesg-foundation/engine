package cmdService

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Test a service
var Test = &cobra.Command{
	Use:               "test SERVICE_FILE",
	Short:             "Start and test the service",
	Long:              "Test the interactions with the service, listening to events and calling tasks.",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service test service.yml",
	Run:               testHandler,
	DisableAutoGenTag: true,
}

func testHandler(cmd *cobra.Command, args []string) {
	if cmd.Flag("event").Value.String() == "" {
		fmt.Println("Logging all events")
	} else {
		fmt.Println("Logging only events ", cmd.Flag("event").Value.String())
	}

	if cmd.Flag("task").Value.String() != "" {
		fmt.Println("Calling task ", cmd.Flag("task").Value.String(), " with data ", cmd.Flag("data").Value.String())
	}
	// TODO add real testing
	fmt.Println("Test complete, keep alive option: ", cmd.Flag("keep-alive").Value.String())
}

func init() {
	Test.Flags().StringP("event", "e", "", "Event filter, will only log those events")
	Test.Flags().StringP("task", "t", "", "Run a specific task")
	Test.Flags().StringP("data", "d", "", "File with the data required to run a specific task")
	Test.Flags().BoolP("keep-alive", "", false, "Let the service run event after the end of the test")
}
