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
