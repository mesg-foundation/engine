package cmdService

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	servicePackage "github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
)

// Test a service
var Test = &cobra.Command{
	Use:   "test",
	Short: "Test a service",
	Long: `Test a service by listening to events or calling tasks.

See more detail on the [Test page from the documentation](https://docs.mesg.tech/service/test.html)`,
	Example: `mesg-core service test
mesg-core service test ./SERVICE_FOLDER
mesg-core service test --event EVENT_NAME
mesg-core service test --task TASK_NAME --data ./PATH_TO_DATA_FILE.yml
mesg-core service test --keep-alive`,
	Run:               testHandler,
	DisableAutoGenTag: true,
}

func listenEvents(serviceID string, filter string) {
	stream, err := cli.ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID: serviceID,
	})
	handleError(err)
	for {
		event, err := stream.Recv()
		if err != nil {
			log.Println(aurora.Red(err))
			return
		}
		if filter == "*" || filter == event.EventKey {
			log.Println("Receive event", aurora.Green(event.EventKey), ":", aurora.Bold(event.EventData))
		}
	}
}

func listenResults(serviceID string) {
	stream, err := cli.ListenResult(context.Background(), &core.ListenResultRequest{
		ServiceID: serviceID,
	})
	handleError(err)
	for {
		result, err := stream.Recv()
		if err != nil {
			log.Println(aurora.Red(err))
			return
		}
		log.Println("Receive result", aurora.Green(result.TaskKey), aurora.Green(result.OutputKey), ":", aurora.Bold(result.OutputData))
	}
}

func executeTask(serviceID string, task string, dataPath string) (execution *core.ExecuteTaskReply, err error) {
	if task == "" {
		return
	}
	var data = []byte("{}")
	if dataPath != "" {
		data, err = ioutil.ReadFile(dataPath)
		handleError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	execution, err = cli.ExecuteTask(ctx, &core.ExecuteTaskRequest{
		ServiceID: serviceID,
		TaskKey:   task,
		TaskData:  string(data),
	})
	handleError(err)
	log.Println("Execute task", aurora.Green(task), "with data", aurora.Bold(string(data)))
	return
}

func testHandler(cmd *cobra.Command, args []string) {
	service := loadService(defaultPath(args))
	_, err := buildDockerImage(defaultPath(args), service.Name)
	handleError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	deployment, err := cli.DeployService(ctx, &core.DeployServiceRequest{
		Service: service,
	})
	handleError(err)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = cli.StartService(ctx, &core.StartServiceRequest{
		ServiceID: deployment.ServiceID,
	})
	handleError(err)

	go listenEvents(deployment.ServiceID, cmd.Flag("event").Value.String())

	go listenResults(deployment.ServiceID)

	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Wait service to run..."})
	<-service.WaitStatus(servicePackage.RUNNING, 1*time.Minute)
	s.Stop()

	executeTask(deployment.ServiceID, cmd.Flag("task").Value.String(), cmd.Flag("data").Value.String())

	<-cmdUtils.WaitForCancel()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = cli.StopService(ctx, &core.StopServiceRequest{
		ServiceID: deployment.ServiceID,
	})
	fmt.Println(err)
}

func init() {
	Test.Flags().StringP("event", "e", "*", "Only log a specific event")
	Test.Flags().StringP("task", "t", "", "Run a specific task")
	Test.Flags().StringP("data", "d", "", "Path to the file containing the data required to run the task")
}
