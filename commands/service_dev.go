package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/x/xpflag"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/spf13/cobra"
)

type serviceDevCmd struct {
	baseCmd

	eventFilter  string
	taskFilter   string
	outputFilter string
	path         string
	env          map[string]string

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
	c.cmd.Flags().Var(xpflag.NewStringToStringValue(&c.env, nil), "env", "set env defined in mesg.yml (configuration.env)")
	return c
}

func (c *serviceDevCmd) preRunE(cmd *cobra.Command, args []string) error {
	c.path = getFirstOrDefault(args)
	return nil
}

func (c *serviceDevCmd) runE(cmd *cobra.Command, args []string) error {
	sid, hash, err := deployService(c.e, c.path, c.env)
	if err != nil {
		return err
	}
	fmt.Printf("%s Service deployed with sid %s and hash %s\n", pretty.SuccessSign, pretty.Success(sid), pretty.Success(hash))

	defer func() {
		var err error
		pretty.Progress("Removing the service...", func() {
			err = c.e.ServiceDelete(false, hash)
		})
		if err != nil {
			fmt.Printf("%s Removing the service completed with an error: %s\n", pretty.FailSign, err)
		} else {
			fmt.Printf("%s Service removed\n", pretty.SuccessSign)
		}
	}()

	pretty.Progress("Starting the service...", func() { err = c.e.ServiceStart(hash) })
	if err != nil {
		return err
	}
	fmt.Printf("%s Service started\n", pretty.SuccessSign)

	listenEventsC, eventsErrC, err := c.e.ServiceListenEvents(hash, c.eventFilter)
	if err != nil {
		return err
	}

	listenResultsC, resultsErrC, err := c.e.ServiceListenResults(hash, c.taskFilter, c.outputFilter, nil)
	if err != nil {
		return err
	}

	closer, logsErrC, err := showLogs(c.e, hash)
	if err != nil {
		return err
	}
	defer closer()

	abort := xsignal.WaitForInterrupt()
	fmt.Println("Listening for results and events from the service...")

	for {
		select {
		case <-abort:
			return nil

		case err := <-logsErrC:
			return err

		case e := <-listenEventsC:
			fmt.Printf("Receive event %s: %s\n",
				pretty.Success(e.EventKey),
				pretty.ColorizeJSON(pretty.FgCyan, nil, false, []byte(e.EventData)),
			)

		case err := <-eventsErrC:
			fmt.Fprintf(os.Stderr, "%s Listening events error: %s\n", pretty.FailSign, err)
			return nil

		case r := <-listenResultsC:
			if r.Error != "" {
				fmt.Printf("Receive execution error on %s task %s: %s\n",
					pretty.Fail(r.TaskKey),
					pretty.Fail(r.OutputKey),
					pretty.Fail(r.Error),
				)
			} else {
				fmt.Printf("Receive execution result on %s task %s: %s\n",
					pretty.Success(r.TaskKey),
					pretty.Colorize(color.New(color.FgCyan), r.OutputKey),
					pretty.ColorizeJSON(pretty.FgCyan, nil, false, []byte(r.OutputData)),
				)
			}

		case err := <-resultsErrC:
			fmt.Fprintf(os.Stderr, "%s Listening results error: %s\n", pretty.FailSign, err)
			return nil
		}
	}
}
