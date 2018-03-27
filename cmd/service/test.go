package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// Test a service
var Test = &cobra.Command{
	Use:   "test",
	Short: "Test a service",
	Long: `Test a service by listening to events or calling tasks.

See more detail on the [Test page from the documentation](https://docs.mesg.tech/service/test.html)`,
	Example: `mesg-cli service test
mesg-cli service test ./SERVICE_FOLDER
mesg-cli service test --event EVENT_NAME
mesg-cli service test --task TASK_NAME --data ./PATH_TO_DATA_FILE.yml
mesg-cli service test --keep-alive`,
	Run:               testHandler,
	DisableAutoGenTag: true,
}

func testHandler(cmd *cobra.Command, args []string) {
	path := "./"
	if len(args) > 0 {
		path = args[0]
	}
	service := importService(path)
	if service == nil {
		return
	}
	_, err := service.Start()
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
	Test.Flags().StringP("event", "e", "", "Only log a specific event")
	Test.Flags().StringP("task", "t", "", "Run a specific task")
	Test.Flags().StringP("data", "d", "", "Path to the file containing the data required to run the task")
	Test.Flags().BoolP("keep-alive", "", false, "Leave the service runs after the end of the test")
}
