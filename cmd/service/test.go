package cmdService

import (
	"log"

	"github.com/mesg-foundation/core/cmd/utils"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
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

func listenEvents(service *service.Service, callback func(event *types.EventReply)) {
	listener := pubsub.Subscribe(service.EventSubscriptionChannel())
	go func() {
		for event := range listener {
			callback(event.(*types.EventReply))
		}
	}()
	return
}

func startServer() {
	server := api.Server{}
	err := server.Serve()
	defer server.Stop()
	if err != nil {
		log.Panicln(err)
	}
}

func testHandler(cmd *cobra.Command, args []string) {
	service := loadService(defaultPath(args))

	go startServer()
	go listenEvents(service, func(event *types.EventReply) {
		filter := cmd.Flag("event").Value.String()
		if filter == "*" || filter == event.Event {
			log.Println("Receive event", aurora.Green(event.Event), ":", aurora.Bold(event.Data))
		}
	})
	go startService(service)
	defer service.Stop()

	<-cmdUtils.WaitForCancel()

	// if cmd.Flag("task").Value.String() != "" {
	// 	fmt.Println("Calling task ", cmd.Flag("task").Value.String(), " with data ", cmd.Flag("data").Value.String())
	// }

	// if cmd.Flag("keep-alive").Value.String() != "true" {
	// 	service.Stop()
	// }
}

func init() {
	Test.Flags().StringP("event", "e", "*", "Only log a specific event")
	// Test.Flags().StringP("task", "t", "", "Run a specific task")
	// Test.Flags().StringP("data", "d", "", "Path to the file containing the data required to run the task")
	// Test.Flags().BoolP("keep-alive", "", false, "Leave the service runs after the end of the test")
}
