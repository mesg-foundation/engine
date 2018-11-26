package commands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/pretty"
	casting "github.com/mesg-foundation/core/utils/servicecasting"
	"github.com/mesg-foundation/core/x/xjson"
	"github.com/mesg-foundation/core/x/xpflag"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type serviceExecuteCmd struct {
	baseCmd

	executeData map[string]string
	taskKey     string
	jsonFile    string

	e ServiceExecutor
}

func newServiceExecuteCmd(e ServiceExecutor) *serviceExecuteCmd {
	c := &serviceExecuteCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "execute",
		Short:   "Execute a task of a service",
		Example: `mesg-core service execute SERVICE`,
		Args:    cobra.ExactArgs(1),
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})
	c.cmd.Flags().StringVarP(&c.taskKey, "task", "t", c.taskKey, "Run the given task")
	c.cmd.Flags().VarP(xpflag.NewStringToStringValue(&c.executeData, nil), "data", "d", "data required to run the task")
	c.cmd.Flags().StringVarP(&c.jsonFile, "json", "j", c.jsonFile, "Path to a JSON file containing the data required to run the task")
	return c
}

func (c *serviceExecuteCmd) preRunE(cmd *cobra.Command, args []string) error {
	if cmd.Flag("data").Changed && cmd.Flag("json").Changed {
		return errors.New("Only one of '--data' or '--json' options can be specified")
	}
	return nil
}

func (c *serviceExecuteCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		s              *coreapi.Service
		result         *coreapi.ResultData
		listenResultsC chan *coreapi.ResultData
		inputData      string
		resultsErrC    chan error
		err            error
	)

	pretty.Progress("Loading the service...", func() {
		s, err = c.e.ServiceByID(args[0])
	})
	if err != nil {
		return err
	}

	if err = c.getTaskKey(s); err != nil {
		return err
	}

	inputData, err = c.getData(c.taskKey, s, c.executeData)
	if err != nil {
		return err
	}
	pretty.Progress(fmt.Sprintf("Executing task %q...", c.taskKey), func() {
		// Create an unique tag that will be used to listen to the result of this exact execution
		tags := []string{uuid.NewV4().String()}

		listenResultsC, resultsErrC, err = c.e.ServiceListenResults(args[0], c.taskKey, "", tags)
		if err != nil {
			return
		}

		// XXX: sleep because listen stream may not be ready to stream the data
		// and execution will done before stream is ready. In that case the response
		// wlll never come TODO: investigate
		time.Sleep(1 * time.Second)

		err = c.e.ServiceExecuteTask(args[0], c.taskKey, inputData, tags)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Task %q executed\n", pretty.SuccessSign, c.taskKey)

	pretty.Progress("Waiting for result...", func() {
		select {
		case result = <-listenResultsC:
		case err = <-resultsErrC:
		}
	})
	if err != nil {
		return err
	}
	fmt.Printf("Task %s returned output %s with data:\n%s\n",
		pretty.Success(c.taskKey),
		pretty.Colorize(pretty.FgCyan, result.OutputKey),
		pretty.ColorizeJSON(pretty.FgCyan, nil, true, []byte(result.OutputData)),
	)
	return nil
}

func (c *serviceExecuteCmd) getTaskKey(s *coreapi.Service) error {
	if c.taskKey == "" {
		keys := taskKeysFromService(s)
		if len(keys) == 1 {
			c.taskKey = keys[0]
			return nil
		}

		if survey.AskOne(&survey.Select{
			Message: "Select the task to execute",
			Options: keys,
		}, &c.taskKey, nil) != nil {
			return errors.New("no task to execute")
		}
	}
	return nil
}

func (c *serviceExecuteCmd) getData(taskKey string, s *coreapi.Service, dataStruct map[string]string) (string, error) {
	if dataStruct != nil {
		castData, err := casting.TaskInputs(s, taskKey, dataStruct)
		if err != nil {
			return "", err
		}

		b, err := json.Marshal(castData)
		return string(b), err
	}

	// see if task has no inputs.
	for _, task := range s.Tasks {
		if task.Key == taskKey {
			if len(task.Inputs) == 0 {
				return "{}", nil
			}
			break
		}
	}

	if c.jsonFile == "" {
		if survey.AskOne(&survey.Input{Message: "Enter the filepath to the inputs"}, &c.jsonFile, nil) != nil {
			return "", errors.New("no filepath given")
		}
	}

	content, err := xjson.ReadFile(c.jsonFile)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func taskKeysFromService(s *coreapi.Service) []string {
	var taskKeys []string
	for _, task := range s.Tasks {
		taskKeys = append(taskKeys, task.Key)
	}
	return taskKeys
}
