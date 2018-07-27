package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var executions []*core.ResultData

// Execute a task from a service
var Execute = &cobra.Command{
	Use:               "execute",
	Short:             "Execute a task of a service",
	Example:           `mesg-core service execute SERVICE_ID`,
	Args:              cobra.MinimumNArgs(1),
	Run:               executeHandler,
	DisableAutoGenTag: true,
}

func init() {
	Execute.Flags().StringP("task", "t", "", "Run the given task")
	Execute.Flags().StringP("json", "j", "", "Path to a JSON file containing the data required to run the task")
}

func executeHandler(cmd *cobra.Command, args []string) {
	serviceID := args[0]
	taskKey := getTaskKey(cmd, serviceID)
	json := getJSON(cmd)
	taskData, err := readJSONFile(json)
	utils.HandleError(err)

	stream, err := cli().ListenResult(context.Background(), &core.ListenResultRequest{
		ServiceID:  serviceID,
		TaskFilter: taskKey,
	})

	var execution *core.ResultData
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Executing task " + aurora.Green(taskKey).String() + "..."}, func() {
		// TODO: Fix this, it's a bit messy to have a sleep here
		go func() {
			time.Sleep(1 * time.Second)
			_, err = executeTask(serviceID, taskKey, taskData)
			utils.HandleError(err)
		}()
		for {
			execution, err = stream.Recv()
			break
		}
	})
	utils.HandleError(err)
	fmt.Println("Task " + aurora.Green(taskKey).String() + " returned output " + aurora.Blue(execution.OutputKey).String() + " with data:")
	fmt.Println(aurora.Bold(execution.OutputData).String())
}

func executeTask(serviceID string, task string, data string) (execution *core.ExecuteTaskReply, err error) {
	return cli().ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID: serviceID,
		TaskKey:   task,
		InputData: data,
	})
}

func taskKeysFromService(s *service.Service) []string {
	var taskKeys []string
	for key := range s.Tasks {
		taskKeys = append(taskKeys, key)
	}
	return taskKeys
}

func getTaskKey(cmd *cobra.Command, serviceID string) string {
	taskKey := cmd.Flag("task").Value.String()
	if taskKey == "" {
		serviceReply, err := cli().GetService(context.Background(), &core.GetServiceRequest{
			ServiceID: serviceID,
		})
		utils.HandleError(err)
		if survey.AskOne(&survey.Select{
			Message: "Select the task to execute",
			Options: taskKeysFromService(serviceReply.Service),
		}, &taskKey, nil) != nil {
			os.Exit(0)
		}
	}
	return taskKey
}

func getJSON(cmd *cobra.Command) string {
	json := cmd.Flag("json").Value.String()
	if json == "" {
		if survey.AskOne(&survey.Input{Message: "Enter the filepath to the inputs"}, &json, nil) != nil {
			os.Exit(0)
		}
	}
	return json
}

func readJSONFile(path string) (string, error) {
	if path == "" {
		return "{}", nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
