// Package deployer deploys and starts system services by using api
// and provides ids of services by their names.
package deployer

import (
	"path/filepath"
	"sync"

	"github.com/mesg-foundation/core/systemservices"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/sirupsen/logrus"
)

// Deployer deploys and starts system services by using api
// and provides ids of services by their names.
type Deployer struct {
	api *api.API

	// instance of the system services manager to update
	ss *systemservices.SystemServices

	// absolute path of system services dir.
	systemServicesPath string
}

// New creates a new Deployer.
// It accepts an instance of the api package.
// It accepts the system services path.
func New(api *api.API, systemServicesPath string, ss *systemservices.SystemServices) *Deployer {
	return &Deployer{
		api:                api,
		ss:                 ss,
		systemServicesPath: systemServicesPath,
	}
}

// Deploy deploys and starts system services from given services list.
// It waits for all system services to run.
// If services are not found, it should return an error.
// If services doesn't start properly, it should return an error.
func (d *Deployer) Deploy(services []string) error {
	if err := d.deployServices(services); err != nil {
		return err
	}
	return d.startServices(services)
}

// deployServices deploys system services.
func (d *Deployer) deployServices(services []string) error {
	var (
		// errs are the deployment errors.
		errs xerrors.Errors
		m    sync.Mutex

		wg sync.WaitGroup
	)

	for _, srv := range services {
		wg.Add(1)
		go func(service string) {
			defer wg.Done()
			logrus.Infof("Deploying system service %q", service)
			sr, err := d.deployService(service)
			m.Lock()
			defer m.Unlock()
			if err != nil {
				errs = append(errs, err)
				return
			}
			logrus.Infof("System service %q deployed", service)
			d.ss.RegisterSystemService(service, sr)
		}(srv)
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
func (d *Deployer) startServices(services []string) error {
	var (
		// errs are the service starting errors.
		errs xerrors.Errors
		m    sync.Mutex

		wg sync.WaitGroup
	)

	for _, srv := range services {
		wg.Add(1)
		go func(service string) {
			defer wg.Done()

			logrus.Infof("Starting system service %q", service)
			if err := d.api.StartService(d.ss.GetServiceID(service)); err != nil {
				m.Lock()
				defer m.Unlock()
				errs = append(errs, err)
			} else {
				logrus.Infof("Starting system service %q", service)
			}
		}(srv)
	}

	wg.Wait()

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	logrus.Info("all system services started")
	return nil
}
