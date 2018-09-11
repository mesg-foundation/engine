package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/fatih/color"
	"github.com/mesg-foundation/core/utils/pretty"
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
		id    string
		valid bool
		err   error
	)
	pretty.Progress("Deploying the service...", func() {
		id, valid, err = c.e.ServiceDeploy(c.path)
	})
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("service deploy invalid")
	}
	fmt.Printf("%s Service deployed\n", pretty.SuccessSign)
	defer pretty.Progress("Removing the service...", func() { c.e.ServiceDelete(id) })

	pretty.Progress("Starting the service...", func() { err = c.e.ServiceStart(id) })
	if err != nil {
		return err
	}
	fmt.Printf("%s Service started with success: %s\n", pretty.SuccessSign, pretty.Success(id))

	listenEventsC, eventsErrC, err := c.e.ServiceListenEvents(c.path, c.eventFilter)
	if err != nil {
		return err
	}

	listenResultsC, resultsErrC, err := c.e.ServiceListenResults(c.path, c.taskFilter, c.outputFilter, nil)
	if err != nil {
		return err
	}

	reader, err := c.e.ServiceLogs(id)
	if err != nil {
		return err
	}
	defer reader.Close()

	go stdcopy.StdCopy(os.Stdout, os.Stderr, reader)

	abort := xsignal.WaitForInterrupt()

loop:
	for {
		select {
		case e := <-listenEventsC:
			fmt.Printf("Receive event %s: %s\n",
				pretty.Success(e.EventKey),
				pretty.ColorizeJSON(pretty.FgCyan, pretty.FgMagenta, []byte(e.EventData)),
			)
		case err := <-eventsErrC:
			fmt.Fprintf(os.Stderr, "%s Listening events error: %s", pretty.FailSign, err)
		case r := <-listenResultsC:
			fmt.Printf("Receive result %s %s with data\n%s\n",
				pretty.Success(r.TaskKey),
				pretty.Colorize(color.New(color.FgCyan), r.OutputKey),
				pretty.ColorizeJSON(pretty.FgBlue, pretty.FgMagenta, []byte(r.OutputData)),
			)
		case err := <-resultsErrC:
			fmt.Fprintf(os.Stderr, "%s Listening results error: %s", pretty.FailSign, err)
		case <-abort:
			break loop
		}
	}
	return nil
}
