package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

type serviceExecuteCmd struct {
	baseCmd

	taskKey  string
	jsonData string
	jsonFile string

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
	c.cmd.Flags().StringVarP(&c.jsonFile, "json", "j", c.jsonData, "JSON containing the data required to run the task")
	c.cmd.Flags().StringVarP(&c.jsonFile, "json-file", "f", c.jsonFile, "Path to a JSON file containing the data required to run the task")

	c.cmd.MarkPersistentFlagRequired("task")
	return c
}

func (c *serviceExecuteCmd) preRunE(cmd *cobra.Command, args []string) error {
	if cmd.Flag("json").Changed && cmd.Flag("json-file").Changed {
		return errors.New("Only one of '--json' or '--json-file' options can be specified")
	}

	if c.jsonFile != "" {
		b, err := ioutil.ReadFile(c.jsonFile)
		if err != nil {
			return err
		}
		c.jsonData = string(b)
	}

	// validate json data
	var js json.RawMessage
	if err := json.Unmarshal([]byte(c.jsonData), &js); err == nil {
		return fmt.Errorf("invalid json: %s", err)
	}

	return nil
}

func (c *serviceExecuteCmd) runE(cmd *cobra.Command, args []string) error {
	// Create an unique tag that will be used to listen to the result of this exact execution
	tags := []string{uuid.NewV4().String()}

	listenResultsC, resultsErrC, err := c.e.ServiceListenResults(args[0], c.taskKey, "", tags)
	if err != nil {
		return err
	}

	// XXX: sleep because listen stream may not be ready to stream the data
	// and execution will done before stream is ready. In that case the response
	// wlll never come TODO: investigate
	time.Sleep(1 * time.Second)

	if err := c.e.ServiceExecuteTask(args[0], c.taskKey, c.jsonData, tags); err != nil {
		return err
	}

	select {
	case result := <-listenResultsC:
		fmt.Printf("Task %s returned output %s with data:\n%s\n",
			pretty.Success(c.taskKey),
			pretty.Colorize(pretty.FgBlue, result.OutputKey),
			pretty.Bold(result.OutputData),
		)

	case err := <-resultsErrC:
		return err

	// XXX: double check if sleep before was too short.
	case <-time.After(5 * time.Second):
		return errors.Errorf("Task %s didn't get any response", c.taskKey)
	}
	return nil
}
