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
	fmt.Println("Deploying the service...")
	id, valid, err := c.e.ServiceDeploy(c.path)
	if err != nil {
		return err
	}
	fmt.Printf("%s Service deployed\n", pretty.SuccessSign)

	if !valid {
		return errors.New("service deploy invalid")
	}
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
	defer close(listenEventsC)
	defer close(eventsErrC)

	listenResultsC, resultsErrC, err := c.e.ServiceListenResults(c.path, c.taskFilter, c.outputFilter, nil)
	if err != nil {
		return err
	}
	defer close(listenResultsC)
	defer close(resultsErrC)

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
			fmt.Println("Receive event", pretty.Success(e.EventKey), ":", pretty.Bold(e.EventData))
		case err := <-eventsErrC:
			fmt.Fprintln(os.Stderr, "Listening events error:", err)
		case r := <-listenResultsC:
			fmt.Println("Receive result", pretty.Success(r.TaskKey), pretty.Colorize(color.New(color.FgCyan), r.OutputKey), "with data", pretty.Bold(r.OutputData))
		case err := <-resultsErrC:
			fmt.Fprintln(os.Stderr, "Listening results error:", err)
		case <-abort:
			break loop
		}
	}
	return nil
}
