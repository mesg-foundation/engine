package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type serviceListCmd struct {
	baseCmd

	e ServiceExecutor
}

// newServiceListCmd receives e to do API calls and w to output structured table logs.
func newServiceListCmd(e ServiceExecutor) *serviceListCmd {
	c := &serviceListCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "list",
		Short: "List all published services",
		Long: `This command returns all published services with basic information.
Optionally, you can filter the services published by a specific developer:
To have more details, see the [detail command](mesg-core_service_detail.md).`,
		Example: `mesg-core service list`,
		Args:    cobra.NoArgs,
		RunE:    c.runE,
	})
	return c
}

func (c *serviceListCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		services []*coreapi.Service
		err      error
	)
	pretty.Progress("Listing services...", func() {
		services, err = c.e.ServiceList()
	})
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	fmt.Fprintf(w, "STATUS\tSERVICE\tNAME\t\n")
	for _, s := range services {
		var status string
		switch s.Status {
		case coreapi.Service_UNKNOWN:
			status = "unknown"
		case coreapi.Service_STOPPED:
			status = "stopped"
		case coreapi.Service_STARTING:
			status = "starting"
		case coreapi.Service_PARTIAL:
			status = "partial"
		case coreapi.Service_RUNNING:
			status = "running"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", status, s.ID, s.Name)
	}
	return w.Flush()
}
