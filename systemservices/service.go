package systemservices

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xerrors"
)

// deploySystemServices deploys system services placed inside the system services path.
func (s *SystemServices) deploySystemServices() ([]*service.Service, error) {
	files, err := ioutil.ReadDir(s.systemServicesPath)
	if err != nil {
		return nil, err
	}

	var (
		services []*service.Service
		errs     xerrors.Errors
		m        sync.Mutex

		wg sync.WaitGroup
	)

	for _, file := range files {
		dirName := file.Name()

		// ignore dot files/dirs. (e.g. .DS_Store)
		if strings.HasPrefix(dirName, ".") {
			continue
		}
		if !file.IsDir() {
			return nil, &notDirectoryError{fileName: dirName}
		}

		wg.Add(1)
		go func(dirName string) {
			defer wg.Done()
			path := filepath.Join(s.systemServicesPath, dirName)
			s, err := s.deployService(path)
			m.Lock()
			defer m.Unlock()
			if err != nil {
				errs = append(errs, err)
				return
			}
			services = append(services, s)
		}(dirName)
	}

	wg.Wait()
	return services, errs.ErrorOrNil()
}

// deployService deploys a service living in path.
func (s *SystemServices) deployService(path string) (*service.Service, error) {
	archive, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return nil, err
	}

	service, validationErr, err := s.api.DeployService(archive)
	if err != nil {
		return nil, err
	}
	if validationErr != nil {
		return nil, validationErr
	}
	return service, nil
}

// startService starts the services.
func (s *SystemServices) startServices(services []*service.Service) error {
	var (
		errs xerrors.Errors
		m    sync.Mutex

		wg sync.WaitGroup
	)

	wg.Add(len(services))
	for _, ss := range services {
		go func(id string) {
			defer wg.Done()
			if err := s.api.StartService(id); err != nil {
				m.Lock()
				defer m.Unlock()
				errs = append(errs, err)
			}
		}(ss.ID)
	}

	wg.Wait()
	return errs.ErrorOrNil()
}

// getServiceID returns the service id of the service that it's name matches with name.
func (s *SystemServices) getServiceID(services []*service.Service, name string) string {
	for _, s := range services {
		if s.Name == name {
			return s.ID
		}
	}
	return ""
}
