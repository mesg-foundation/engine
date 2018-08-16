package mesg

import (
	"io"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
)

// DeployServiceOption is a configuration func for Deploy methods.
type DeployServiceOption func(*serviceDeployer)

// DeployServiceStatusOption receives chan statuses to send deploy statuses.
func DeployServiceStatusOption(statuses chan string) DeployServiceOption {
	return func(deployer *serviceDeployer) {
		deployer.Statuses = statuses
	}
}

// DeployService deploys a service from a gzipped tarball.
func (m *MESG) DeployService(r io.Reader, options ...DeployServiceOption) (*service.Service, *importer.ValidationError, error) {
	deployer := newServiceDeployer(m)
	for _, option := range options {
		option(deployer)
	}
	return deployer.FromGzippedTar(r)
}

// DeployServiceFromURL deploys a service lives at a Git host.
func (m *MESG) DeployServiceFromURL(url string, options ...DeployServiceOption) (*service.Service, *importer.ValidationError, error) {
	deployer := newServiceDeployer(m)
	for _, option := range options {
		option(deployer)
	}
	return deployer.FromGitURL(url)
}
