package commands

import (
	"io"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/servicetemplate"
	"github.com/spf13/cobra"
)

// RootExecutor is an interface that handles core commands.
type RootExecutor interface {
	Start() error
	Stop() error
	Status() (container.StatusType, error)
	Logs() (io.ReadCloser, error)
}

// ServiceExecutor is an interface that handles services commands.
type ServiceExecutor interface {
	ServiceByID(id string) (*coreapi.Service, error)
	ServiceDeleteAll(deleteData bool) error
	ServiceDelete(deleteData bool, ids ...string) error
	ServiceDeploy(path string, env map[string]string, statuses chan provider.DeployStatus) (id string, validationError, err error)
	ServicePublishDefinitionFile(path string) (string, error)
	ServiceListenEvents(id, eventFilter string) (chan *coreapi.EventData, chan error, error)
	ServiceListenResults(id, taskFilter, outputFilter string, tagFilters []string) (chan *coreapi.ResultData, chan error, error)
	ServiceLogs(id string, dependencies ...string) (logs []*provider.Log, closer func(), err error)
	ServiceExecuteTask(id, taskKey, inputData string, tags []string) error
	ServiceStart(id string) error
	ServiceStop(id string) error
	ServiceValidate(path string) (string, error)
	ServiceGenerateDocs(path string) error
	ServiceList() ([]*coreapi.Service, error)
	ServiceInitTemplateList() ([]*servicetemplate.Template, error)
	ServiceInitDownloadTemplate(t *servicetemplate.Template, dst string) error
}

// WalletExecutor is an interface that handles wallet commands.
type WalletExecutor interface {
	List() ([]common.Address, error)
	Create(passphrase string) (common.Address, error)
	Delete(passphrase string) error
	Export(passphrase string) ([]byte, error)
	Import(passphrase string) error
	Sign(passphrase string) (*types.Transaction, error)
}

// Executor is an interface that keeps all commands interfaces.
type Executor interface {
	RootExecutor
	ServiceExecutor
}

// Build constructs root command and returns it.
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

// discardOutput discards usage and error messages.
func (c *baseCmd) discardOutput() {
	c.cmd.SetOutput(ioutil.Discard)
}

// getFirstOrDefault returns directory if args len is gt 0 or current directory.
func getFirstOrDefault(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "./"
}
