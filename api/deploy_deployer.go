package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/x/xdocker/xarchive"
	"github.com/mesg-foundation/core/x/xgit"
	uuid "github.com/satori/go.uuid"
)

// serviceDeployer provides functionalities to deploy a MESG service.
type serviceDeployer struct {
	// statuses receives status messages produced during deployment.
	statuses chan DeployStatus

	api *API
}

// StatusType indicates the type of status message.
type StatusType int

const (
	// RUNNING indicates that status message belongs to an active state.
	RUNNING StatusType = iota + 1

	// DONE indicates that status message belongs to completed state.
	DONE
)

// DeployStatus represents the deployment status.
type DeployStatus struct {
	Message string
	Type    StatusType
}

// newServiceDeployer creates a new serviceDeployer with given api and options.
func newServiceDeployer(api *API, options ...DeployServiceOption) *serviceDeployer {
	d := &serviceDeployer{
		api: api,
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// FromGitURL deploys a service hosted at a Git url.
func (d *serviceDeployer) FromGitURL(url string) (*service.Service, *importer.ValidationError, error) {
	d.sendStatus("Downloading service...", RUNNING)
	path, err := d.createTempDir()
	if err != nil {
		return nil, nil, err
	}
	defer os.RemoveAll(path)
	if err := xgit.Clone(url, path); err != nil {
		return nil, nil, err
	}
	d.sendStatus(fmt.Sprintf("%s Service downloaded with success.", aurora.Green("✔")), DONE)
	r, err := xarchive.GzippedTar(path)
	if err != nil {
		return nil, nil, err
	}
	return d.deploy(r)
}

// FromGzippedTar deploys a service from a gzipped tarball.
func (d *serviceDeployer) FromGzippedTar(r io.Reader) (*service.Service, *importer.ValidationError, error) {
	return d.deploy(r)
}

// deploy deploys a service in path.
func (d *serviceDeployer) deploy(r io.Reader) (*service.Service, *importer.ValidationError, error) {
	statuses := make(chan service.DeployStatus)
	go d.forwardDeployStatuses(statuses)

	s, err := service.New(r,
		service.ContainerOption(d.api.container),
		service.DeployStatusOption(statuses),
	)
	validationErr, err := d.assertValidationError(err)
	if err != nil {
		return nil, nil, err
	}
	if validationErr != nil {
		return nil, validationErr, nil
	}
	return s, nil, services.Save(s)
}

func (d *serviceDeployer) createTempDir() (path string, err error) {
	return ioutil.TempDir("", "mesg-"+uuid.NewV4().String())
}

// sendStatus sends a status message.
func (d *serviceDeployer) sendStatus(message string, typ StatusType) {
	if d.statuses != nil {
		d.statuses <- DeployStatus{
			Message: message,
			Type:    typ,
		}
	}
}

// forwardStatuses forwards status messages.
func (d *serviceDeployer) forwardDeployStatuses(statuses chan service.DeployStatus) {
	for status := range statuses {
		var t StatusType
		switch status.Type {
		case service.DRUNNING:
			t = RUNNING
		case service.DDONE:
			t = DONE
		}
		d.sendStatus(status.Message, t)
	}
}

func (d *serviceDeployer) assertValidationError(err error) (*importer.ValidationError, error) {
	if err == nil {
		return nil, nil
	}
	if validationError, ok := err.(*importer.ValidationError); ok {
		return validationError, nil
	}
	return nil, err
}
