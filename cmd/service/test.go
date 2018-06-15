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
mesg-core service test --event-filter EVENT_NAME
mesg-core service test --task TASK_NAME --data ./PATH_TO_DATA_FILE.json
mesg-core service test --task-filter TASK_NAME --output-filter OUTPUT_NAME
mesg-core service test --serviceID SERVICE_ID --keep-alive`,
	Run:               testHandler,
	DisableAutoGenTag: true,
}

func init() {
	Test.Flags().StringP("task", "t", "", "Run the given task")
	Test.Flags().StringP("data", "d", "", "Path to the file containing the data required to run the task")
	Test.Flags().StringP("serviceID", "s", "", "ID of a previously deployed service")
	Test.Flags().BoolP("keep-alive", "", false, "Do not stop the service at the end of this command")
	Test.Flags().StringP("event-filter", "e", "*", "Only log the data of the given event")
	Test.Flags().StringP("task-filter", "r", "", "Only log the result of the given task")
	Test.Flags().StringP("output-filter", "o", "", "Only log the data of the given output of a task result. If set, you also need to set the task in --task-filter")
}

func listenEvents(serviceID string, filter string) {
	stream, err := cli.ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID:   serviceID,
		EventFilter: filter,
	})
	cmdUtils.HandleError(err)
	fmt.Println(aurora.Cyan("Listening for events from the service..."))
	for {
		event, err := stream.Recv()
		if err != nil {
			log.Println(aurora.Red(err))
			return
		}
		log.Println("Receive event", aurora.Green(event.EventKey), ":", aurora.Bold(event.EventData))
	}
}

func listenResults(serviceID string, result string, output string) {
	stream, err := cli.ListenResult(context.Background(), &core.ListenResultRequest{
		ServiceID:    serviceID,
		TaskFilter:   result,
		OutputFilter: output,
	})
	cmdUtils.HandleError(err)
	fmt.Println(aurora.Cyan("Listening for results from the service..."))
	for {
		result, err := stream.Recv()
		if err != nil {
			log.Println(aurora.Red(err))
			return
		}
		log.Println("Receive result", aurora.Green(result.TaskKey), aurora.Cyan(result.OutputKey), "with data", aurora.Bold(result.OutputData))
	}
}

func executeTask(serviceID string, task string, dataPath string) (execution *core.ExecuteTaskReply, err error) {
	if task == "" {
		return
	}
	var data = []byte("{}")
	if dataPath != "" {
		data, err = ioutil.ReadFile(dataPath)
		cmdUtils.HandleError(err)
	}

	execution, err = cli.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID: serviceID,
		TaskKey:   task,
		TaskData:  string(data),
	})
	cmdUtils.HandleError(err)
	log.Println("Execute task", aurora.Green(task), "with data", aurora.Bold(string(data)))
	return
}

func testHandler(cmd *cobra.Command, args []string) {
	var err error
	serviceID := cmd.Flag("serviceID").Value.String()
	if serviceID == "" {
		service := prepareService(defaultPath(args))
		deployment, err := cli.DeployService(context.Background(), &core.DeployServiceRequest{
			Service: service,
		})
		cmdUtils.HandleError(err)
		serviceID = deployment.ServiceID
		fmt.Println(aurora.Green("Service deployed with success"))
		fmt.Println("Service ID:", serviceID)

		cmdUtils.ShowSpinnerForFunc(cmdUtils.SpinnerOptions{Text: "Starting service..."}, func() {
			_, err = cli.StartService(context.Background(), &core.StartServiceRequest{
				ServiceID: serviceID,
			})
		})
		cmdUtils.HandleError(err)
		fmt.Println(aurora.Green("Service started"))
	}

	go listenEvents(serviceID, cmd.Flag("event-filter").Value.String())
	go listenResults(serviceID, cmd.Flag("task-filter").Value.String(), cmd.Flag("output-filter").Value.String())

	time.Sleep(time.Second)
	executeTask(serviceID, cmd.Flag("task").Value.String(), cmd.Flag("data").Value.String())
	<-cmdUtils.WaitForCancel()

	if cmd.Flag("keep-alive").Value.String() != "true" {
		cmdUtils.ShowSpinnerForFunc(cmdUtils.SpinnerOptions{Text: "Stopping service..."}, func() {
			_, err = cli.StopService(context.Background(), &core.StopServiceRequest{
				ServiceID: serviceID,
			})
		})
		cmdUtils.HandleError(err)
		fmt.Println(aurora.Green("Service stopped"))
	}
}
