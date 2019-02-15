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
	CreateManifest(path string) (provider.ManifestData, error)
	UploadServiceFiles(path string, manifest provider.ManifestData) (protocol string, source string, err error)
	CreateService(sid, from string) (*provider.Transaction, error)
	CreateServiceVersion(sidHash, hash, manifest, manifestProtocol, from string) (*provider.Transaction, error)
	CreateServiceOffer(sidHash, price, duration, from string) (*provider.Transaction, error)
	DisableServiceOffer(sidHash, offerIndex, from string) (*provider.Transaction, error)
	Purchase(sidHash, offerIndex, from string) (*provider.Transaction, error)
	TransferServiceOwnership(sidHash, newOwner, from string) (*provider.Transaction, error)
	SendSignedTransaction(signedTransaction string) (*provider.TransactionReceipt, error)
	IsAuthorized(sidHash string) (bool, error)
	ServiceExist(sid string) (bool, error)
	GetService(sid string) (*provider.MarketplaceService, error)
}

// WalletExecutor is an interface that handles wallet commands.
type WalletExecutor interface {
	List() ([]string, error)
	Create(passphrase string) (string, error)
	Delete(address string, passphrase string) (string, error)
	Export(address string, passphrase string) (provider.EncryptedKeyJSONV3, error)
	Import(account provider.EncryptedKeyJSONV3, passphrase string) (string, error)
	ImportFromPrivateKey(privateKey string, passphrase string) (string, error)
	Sign(address string, passphrase string, transaction *provider.Transaction) (string, error)
}

// Executor is an interface that keeps all commands interfaces.
type Executor interface {
	RootExecutor
	ServiceExecutor
	MarketplaceExecutor
	WalletExecutor
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
