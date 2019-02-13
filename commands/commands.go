package commands

import (
	"io"
	"io/ioutil"

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

// MarketplaceExecutor is an interface that handles marketplace commands.
type MarketplaceExecutor interface {
	PublishDefinitionFile(path string) (string, error)
	CreateService(sid, from string) (*provider.TransactionOutput, error)
	CreateServiceVersion(sidHash, hash, manifest, manifestProtocol, from string) (*provider.TransactionOutput, error)
	CreateServiceOffer(sidHash, price, duration, from string) (*provider.TransactionOutput, error)
	DisableServiceOffer(sidHash, offerIndex, from string) (*provider.TransactionOutput, error)
	Purchase(sidHash, offerIndex, from string) (*provider.TransactionOutput, error)
	TransferServiceOwnership(sidHash, newOwner, from string) (*provider.TransactionOutput, error)
	SendSignedTransaction(signedTransaction string) (*provider.SendSignedTransactionTaskSuccessOutput, error)
	IsAuthorized(sidHash string) (bool, error)
}

// Executor is an interface that keeps all commands interfaces.
type Executor interface {
	RootExecutor
	ServiceExecutor
	MarketplaceExecutor
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
