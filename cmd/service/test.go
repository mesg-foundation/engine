package cmdService

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/config"
	"google.golang.org/grpc"

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
	server := api.Server{
		Network: config.Api.Server.Network(),
		Address: config.Api.Server.Address(),
	}
	err := server.Serve()
	defer server.Stop()
	if err != nil {
		log.Panicln(err)
	}
}

func executeTask(service *service.Service, task string, dataPath string) (reply *types.TaskReply, err error) {
	connection, err := grpc.Dial(config.Api.Client.Target(), grpc.WithInsecure())
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return
	}

	cli := types.NewTaskClient(connection)
	reply, err = cli.Execute(context.Background(), &types.ExecuteTaskRequest{
		Service: &types.ProtoService{
			Name: service.Name,
		},
		Task: task,
		Data: string(data),
	})
	log.Println("Execute task", task, "with data", string(data))
	return
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
	startService(service)
	log.Println(aurora.Green("Service started"))
	defer service.Stop()
	task := cmd.Flag("task").Value.String()
	if task != "" {
		_, err := executeTask(service, task, cmd.Flag("data").Value.String())
		if err != nil {
			log.Println(aurora.Red(err))
		}
	}

	<-cmdUtils.WaitForCancel()
}

func init() {
	Test.Flags().StringP("event", "e", "*", "Only log a specific event")
	Test.Flags().StringP("task", "t", "", "Run a specific task")
	Test.Flags().StringP("data", "d", "", "Path to the file containing the data required to run the task")
}
