package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/cobra"
)

type serviceStatus struct {
	service *service.Service
	status  service.StatusType
}

func (s serviceStatus) String() string {
	statusText := map[service.StatusType]aurora.Value{
		service.STOPPED: aurora.Red("[Stopped]"),
		service.RUNNING: aurora.Green("[Running]"),
		service.PARTIAL: aurora.Brown("[Partial]"),
	}
	return strings.Join([]string{
		"-",
		statusText[s.status].String(),
		aurora.Bold(s.service.Hash()).String(),
		s.service.Name,
	}, " ")
}

type byStatus []serviceStatus

func (a byStatus) Len() int           { return len(a) }
func (a byStatus) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byStatus) Less(i, j int) bool { return a[j].status < a[i].status }

// List all the services
var List = &cobra.Command{
	Use:   "list",
	Short: "List all published services",
	Long: `This command returns all published services with basic information.
Optionally, you can filter the services published by a specific developer:
To have more details, see the [detail command](mesg-core_service_detail.md).`,
	Example:           `mesg-core service list`,
	Run:               listHandler,
	DisableAutoGenTag: true,
}

func listHandler(cmd *cobra.Command, args []string) {
	reply, err := cli.ListServices(context.Background(), &core.ListServicesRequest{})
	utils.HandleError(err)
	status, err := servicesWithStatus(reply.Services)
	utils.HandleError(err)
	sort.Sort(byStatus(status))
	for _, serviceStatus := range status {
		fmt.Println(serviceStatus)
	}
}

func servicesWithStatus(services []*service.Service) (status []serviceStatus, err error) {
	for _, s := range services {
		st, err := s.Status()
		if err != nil {
			break
		}
		status = append(status, serviceStatus{
			service: s,
			status:  st,
		})
	}
	return
}
