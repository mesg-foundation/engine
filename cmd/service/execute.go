package service

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
)

var executions []*core.ResultData

// Execute a task from a service
var Execute = &cobra.Command{
	Use:               "execute",
	Short:             "Execute a task from a service",
	Example:           `mesg-core service execute`,
	Run:               executeHandler,
	DisableAutoGenTag: true,
}

func init() {
	Execute.Flags().StringP("service", "s", "", "ID of a previously deployed service")
	Execute.Flags().StringP("key", "t", "", "Run the given task")
	Execute.Flags().StringP("json", "d", "", "Path to the file containing the data required to run the task")
}

func executeHandler(cmd *cobra.Command, args []string) {
	serviceID := cmd.Flag("service").Value.String()
	taskKey := cmd.Flag("key").Value.String()
	taskData, err := getJSON(cmd.Flag("json").Value.String())
	utils.HandleError(err)

	stream, err := cli.ListenResult(context.Background(), &core.ListenResultRequest{
		ServiceID:  serviceID,
		TaskFilter: taskKey,
	})
	// TODO: Fix this, it's a bit messy to have a sleep here
	go func() {
		time.Sleep(1 * time.Second)
		_, err = executeTask(serviceID, taskKey, taskData)
		utils.HandleError(err)
	}()
	for {
		execution, err := stream.Recv()
		utils.HandleError(err)
		log.Println("Result of task", aurora.Green(taskKey), "with output", aurora.Blue(execution.OutputKey), aurora.Bold(execution.OutputData))
		break
	}
}

func getJSON(path string) (string, error) {
	if path == "" {
		return "{}", nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func executeTask(serviceID string, task string, data string) (execution *core.ExecuteTaskReply, err error) {
	return cli.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID: serviceID,
		TaskKey:   task,
		InputData: data,
	})
}
