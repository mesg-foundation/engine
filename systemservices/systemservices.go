package systemservices

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/systemservices/resolver"
	"github.com/mesg-foundation/core/x/xerrors"
)

// SystemServices is managing all system services.
// It is responsible to start all system services when the core start.
// It reads the services' ID from the config package.
// All system services should runs all the time.
// Any interaction with the system services are done by using the api package.
type SystemServices struct {
	api                *api.API
	systemServicesPath string
}

// New creates a new SystemServices instance.
// It accepts an instance of the api package.
// It accepts the system services path.
// It reads the services' ID from the config package.
// It starts all system services.
// It waits for all system services to run.
// If services' ID are not in the config, it should return an error.
// IF the services don't start properly, it should return an error.
func New(api *api.API, systemServicesPath string) (*SystemServices, error) {
	s := &SystemServices{
		api:                api,
		systemServicesPath: systemServicesPath,
	}
	_, err := s.deploySystemServices()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Resolver returns the Resolver instance using the running Resolver service.
func (s *SystemServices) Resolver() (*resolver.Resolver, error) {
	return nil, nil
}

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
