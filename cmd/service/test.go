package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	"github.com/mesg-foundation/application/service"
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
	service, err := service.ImportFromFile(args[0])
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	_, err = service.Start()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	fmt.Println(aurora.Green("Service started"))

	if cmd.Flag("event").Value.String() == "" {
		fmt.Println("Logging all events")
	} else {
		fmt.Println("Logging only events ", cmd.Flag("event").Value.String())
	}

	if cmd.Flag("task").Value.String() != "" {
		fmt.Println("Calling task ", cmd.Flag("task").Value.String(), " with data ", cmd.Flag("data").Value.String())
	}

	if cmd.Flag("keep-alive").Value.String() != "true" {
		service.Stop()
	}
}

func init() {
	Test.Flags().StringP("event", "e", "", "Event filter, will only log those events")
	Test.Flags().StringP("task", "t", "", "Run a specific task")
	Test.Flags().StringP("data", "d", "", "File with the data required to run a specific task")
	Test.Flags().BoolP("keep-alive", "", false, "Let the service run event after the end of the test")
}
