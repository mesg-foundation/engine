package commands

import (
	"io"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/utils/servicetemplate"
	"github.com/spf13/cobra"
)

type RootExecutor interface {
	Start() error
	Stop() error
	Status() (container.StatusType, error)
	Logs() (io.ReadCloser, error)
}

type ServiceExecutor interface {
	ServiceByID(id string) (*core.Service, error)
	ServiceDeleteAll() error
	ServiceDelete(ids ...string) error
	ServiceDeploy(path string) (id string, valid bool, err error)
	ServiceListenEvents(id, eventFilter string) (chan *core.EventData, chan error, error)
	ServiceListenResults(id, taskFilter, outputFilter string, tagFilters []string) (chan *core.ResultData, chan error, error)
	ServiceLogs(id string) (io.ReadCloser, error)
	ServiceDependencyLogs(id string, dependency string) ([]io.ReadCloser, error)
	ServiceExecuteTask(id, taskKey, inputData string, tags []string) error
	ServiceStart(id string) error
	ServiceStop(id string) error
	ServiceValidate(path string) (string, error)
	ServiceGenerateDocs(path string) error
	ServiceList() ([]*core.Service, error)
	ServiceInitTemplateList() ([]*servicetemplate.Template, error)
	ServiceInitDownloadTemplate(t *servicetemplate.Template, dst string) error
	ServiceInitExecuteTemplate(dst string, option servicetemplate.ConfigOption) error
}

type Executor interface {
	RootExecutor
	ServiceExecutor
}

func Build(e Executor) *cobra.Command {
	return newRootCmd(e).cmd
}

type baseCmd struct {
	cmd *cobra.Command
}

// newCommand set default options for given command.
func newCommand(c *cobra.Command) *cobra.Command {
	c.DisableAutoGenTag = true
	return c
}

func getFirstOrDefault(args []string, def string) string {
	if len(args) > 0 {
		return args[0]
	}
	return def
}
