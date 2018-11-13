// Package deployer deploys and starts system services by using api
// and provides ids of services by their names.
package deployer

import (
	"path/filepath"
	"sync"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/sirupsen/logrus"
)

// systemService represents a system service.
type systemService struct {
	*service.Service

	// name is the unique name of system service.
	// it's also the relative path of system service in the filesystem.
	name string
}

// Deployer deploys and starts system services by using api
// and provides ids of services by their names.
type Deployer struct {
	api *api.API

	// absolute path of system services dir.
	systemServicesPath string

	// all deployed system services
	services []*systemService
}

// New creates a new Deployer.
// It accepts an instance of the api package.
// It accepts the system services path.
func New(api *api.API, systemServicesPath string) *Deployer {
	return &Deployer{
		api:                api,
		systemServicesPath: systemServicesPath,
	}
}

// Deploy deploys and starts system services from given services list.
// It waits for all system services to run.
// If services are not found, it should return an error.
// If services doesn't start properly, it should return an error.
func (d *Deployer) Deploy(services []string) error {
	for _, name := range services {
		d.services = append(d.services, &systemService{name: name})
	}
	if err := d.deployServices(); err != nil {
		return err
	}
	return d.startServices()
}

// GetServiceID returns the service id of a system service that matches with name.
// name compared with the unique name/relative path of system service.
func (d *Deployer) GetServiceID(name string) string {
	for _, s := range d.services {
		if s.name == name {
			return s.ID
		}
	}
	panic("unreachable")
}

// deployServices deploys system services.
func (d *Deployer) deployServices() error {
	var (
		// errs are the deployment errors.
		errs xerrors.Errors
		m    sync.Mutex

		wg sync.WaitGroup
	)

	logrus.Infof("deploying (%d) system services...", len(d.services))

	for _, ss := range d.services {
		wg.Add(1)
		go func(ss *systemService) {
			defer wg.Done()
			sr, err := d.deployService(ss.name)
			m.Lock()
			defer m.Unlock()
			if err != nil {
				errs = append(errs, err)
				return
			}
			logrus.Infof("'%s' system service deployed", ss.name)
			ss.Service = sr
		}(ss)
	}

	wg.Wait()
	return errs.ErrorOrNil()
}

// deployService deploys a system service living in relativePath.
func (d *Deployer) deployService(relativePath string) (*service.Service, error) {
	path := filepath.Join(d.systemServicesPath, relativePath)
	exists, err := xos.DirExists(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &systemServiceNotFoundError{name: relativePath}
	}

	archive, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return nil, err
	}

	sr, validationErr, err := d.api.DeployService(archive)
	if err != nil {
		return nil, err
	}
	if validationErr != nil {
		return nil, validationErr
	}
	return sr, nil
}

// startService starts the system services.
func (d *Deployer) startServices() error {
	var (
		// errs are the service starting errors.
		errs xerrors.Errors
		m    sync.Mutex

		wg sync.WaitGroup
	)

	logrus.Info("starting system services...")

	for _, ss := range d.services {
		wg.Add(1)
		go func(ss *systemService) {
			defer wg.Done()
			if err := d.api.StartService(ss.ID); err != nil {
				m.Lock()
				defer m.Unlock()
				errs = append(errs, err)
			}
		}(ss)
	}

	wg.Wait()

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	logrus.Info("all system services started")
	return nil
}
