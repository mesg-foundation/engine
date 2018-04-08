package cmdService

import (
	"fmt"
	"log"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/queue"
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
	importedService := importService(path)
	if importedService == nil {
		return
	}
	_, err := importedService.Start()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	fmt.Println(aurora.Green("Service started"))
	q := queue.Queue{URL: "amqp://guest:guest@127.0.0.1:5672/"}

	if cmd.Flag("task").Value.String() != "" {
		fmt.Println("Calling task ", cmd.Flag("task").Value.String(), " with data ", cmd.Flag("data").Value.String())
		err := q.Publish(importedService.Namespace(), []queue.Channel{
			queue.Channel{
				Kind: queue.Tasks,
				Name: cmd.Flag("task").Value.String(),
			},
		}, cmd.Flag("data").Value.String())
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("listening events...")
	err = q.Listen(importedService.Namespace(), []queue.Channel{
		queue.Channel{
			Kind: queue.Events,
			Name: cmd.Flag("event").Value.String(),
		},
	}, onEvent)

	if err != nil {
		panic(err)
	}

	if cmd.Flag("keep-alive").Value.String() != "true" {
		importedService.Stop()
	}
}

func onEvent(event interface{}) {
	log.Println(event)
}

func init() {
	Test.Flags().StringP("event", "e", "*", "Only log a specific event")
	Test.Flags().StringP("task", "t", "", "Run a specific task")
	Test.Flags().StringP("data", "d", "", "Path to the file containing the data required to run the task")
	Test.Flags().BoolP("keep-alive", "", false, "Leave the service runs after the end of the test")
}
