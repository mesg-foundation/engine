package commands

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/spf13/cobra"
)

type serviceDevCmd struct {
	baseCmd

	eventFilter  string
	taskFilter   string
	outputFilter string
	path         string

	e ServiceExecutor
}

func newServiceDevCmd(e ServiceExecutor) *serviceDevCmd {
	c := &serviceDevCmd{
		e:           e,
		eventFilter: "*",
	}

	c.cmd = newCommand(&cobra.Command{
		Use:     "dev",
		Short:   "Run your service in development mode",
		Example: "mesg-core service dev PATH",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: c.preRunE,
		RunE:    c.runE,
	})
	c.cmd.Flags().StringVarP(&c.eventFilter, "event-filter", "e", c.eventFilter, "Only log the data of the given event")
	c.cmd.Flags().StringVarP(&c.taskFilter, "task-filter", "t", "", "Only log the result of the given task")
	c.cmd.Flags().StringVarP(&c.outputFilter, "output-filter", "o", "", "Only log the data of the given output of a task result. If set, you also need to set the task in --task-filter")
	return c
}

func (c *serviceDevCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args, "./")
	return nil
}

func (c *serviceDevCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		statuses = make(chan provider.DeployStatus)
		wg       sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		printDeployStatuses(statuses)
	}()

	id, validationError, err := c.e.ServiceDeploy(c.path, statuses)
	wg.Wait()

	pretty.DestroySpinner()
	if err != nil {
		return err
	}
	if validationError != nil {
		return xerrors.Errors{
			validationError,
			errors.New("To get more information, run: mesg-core service validate"),
		}
	}
	fmt.Printf("%s Service deployed with ID: %v\n", pretty.SuccessSign, pretty.Success(id))
	defer func() {
		var err error
		pretty.Progress("Removing the service...", func() {
			err = c.e.ServiceDelete(id)
		})
		if err != nil {
			fmt.Printf("%s Error while removing the service\n", pretty.FailSign)
			fmt.Printf(pretty.Failln(err))
		} else {
			fmt.Printf("%s Service removed\n", pretty.SuccessSign)
		}
	}()

	pretty.Progress("Starting the service...", func() { err = c.e.ServiceStart(id) })
	if err != nil {
		return err
	}
	fmt.Printf("%s Service started\n", pretty.SuccessSign)

	listenEventsC, eventsErrC, err := c.e.ServiceListenEvents(id, c.eventFilter)
	if err != nil {
		return err
	}

	listenResultsC, resultsErrC, err := c.e.ServiceListenResults(id, c.taskFilter, c.outputFilter, nil)
	if err != nil {
		return err
	}

	closer, err := showLogs(c.e, id)
	if err != nil {
		return err
	}
	defer closer()

	abort := xsignal.WaitForInterrupt()

	for {
		select {
		case <-abort:
			return nil

		case e := <-listenEventsC:
			fmt.Printf("Receive event %s: %s\n",
				pretty.Success(e.EventKey),
				pretty.ColorizeJSON(pretty.FgCyan, nil, false, []byte(e.EventData)),
			)

		case err := <-eventsErrC:
			fmt.Fprintf(os.Stderr, "%s Listening events error: %s", pretty.FailSign, err)

		case r := <-listenResultsC:
			fmt.Printf("Receive result %s %s: %s\n",
				pretty.Success(r.TaskKey),
				pretty.Colorize(color.New(color.FgCyan), r.OutputKey),
				pretty.ColorizeJSON(pretty.FgCyan, nil, false, []byte(r.OutputData)),
			)

		case err := <-resultsErrC:
			fmt.Fprintf(os.Stderr, "%s Listening results error: %s", pretty.FailSign, err)
		}
	}
}
