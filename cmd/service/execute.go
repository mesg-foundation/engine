package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/xpflag"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var executeData map[string]string

// Execute a task from a service
var Execute = &cobra.Command{
	Use:               "execute",
	Short:             "Execute a task of a service",
	Example:           `mesg-core service execute SERVICE_ID`,
	Args:              cobra.ExactArgs(1),
	PreRun:            executePreRun,
	Run:               executeHandler,
	DisableAutoGenTag: true,
}

func init() {
	Execute.Flags().StringP("task", "t", "", "Run the given task")
	Execute.Flags().VarP(xpflag.NewStringToStringValue(&executeData, nil), "data", "d", "data required to run the task")
	Execute.Flags().StringP("json", "j", "", "Path to a JSON file containing the data required to run the task")
}

func executePreRun(cmd *cobra.Command, args []string) {
	if cmd.Flag("data").Changed && cmd.Flag("json").Changed {
		utils.HandleError(errors.New("You can specify only one of '--data' or '--json' options"))
	}
}

func executeHandler(cmd *cobra.Command, args []string) {
	serviceID := args[0]
	serviceReply, err := cli().GetService(context.Background(), &core.GetServiceRequest{
		ServiceID: serviceID,
	})
	utils.HandleError(err)
	taskKey := getTaskKey(cmd, serviceReply.Service)
	taskData, err := getData(cmd, taskKey, serviceReply.Service, executeData)
	utils.HandleError(err)

	// Create an unique tag that will be used to listen to the result of this exact execution
	tags := []string{uuid.NewV4().String()}
	stream, err := cli().ListenResult(context.Background(), &core.ListenResultRequest{
		ServiceID:  serviceID,
		TaskFilter: taskKey,
		TagFilters: tags,
	})
	utils.HandleError(err)

	var execution *core.ResultData
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Executing task " + aurora.Green(taskKey).String() + "..."}, func() {
		// TODO: Fix this, it's a bit messy to have a sleep here
		go func() {
			time.Sleep(1 * time.Second)
			_, err = executeTask(serviceID, taskKey, taskData, tags)
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

func executeTask(serviceID string, task string, data string, tags []string) (execution *core.ExecuteTaskReply, err error) {
	return cli().ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
		ServiceID:     serviceID,
		TaskKey:       task,
		InputData:     data,
		ExecutionTags: tags,
	})
}

func taskKeysFromService(s *service.Service) []string {
	var taskKeys []string
	for key := range s.Tasks {
		taskKeys = append(taskKeys, key)
	}
	return taskKeys
}

func getTaskKey(cmd *cobra.Command, s *service.Service) string {
	taskKey := cmd.Flag("task").Value.String()
	if taskKey == "" {
		if survey.AskOne(&survey.Select{
			Message: "Select the task to execute",
			Options: taskKeysFromService(s),
		}, &taskKey, nil) != nil {
			os.Exit(0)
		}
	}
	return taskKey
}

func getData(cmd *cobra.Command, taskKey string, s *service.Service, dataStruct map[string]string) (string, error) {
	data := cmd.Flag("data").Value.String()
	jsonFile := cmd.Flag("json").Value.String()

	if data != "" {
		castData, err := s.Cast(taskKey, dataStruct)
		if err != nil {
			return "", err
		}

		b, err := json.Marshal(castData)
		return string(b), err
	}

	if jsonFile == "" {
		if survey.AskOne(&survey.Input{Message: "Enter the filepath to the inputs"}, &jsonFile, nil) != nil {
			os.Exit(0)
		}
	}
	return readJSONFile(jsonFile)
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
