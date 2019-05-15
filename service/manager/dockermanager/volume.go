package dockermanager

import (
	"sync"

	"github.com/docker/docker/client"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xerrors"
)

// Delete deletes the data volumes of service and its dependencies.
func (m *DockerManager) Delete(s *service.Service) error {
	var (
		wg   sync.WaitGroup
		errs xerrors.SyncErrors
	)
	for _, d := range append(s.Dependencies, s.Configuration) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		for _, volume := range extractVolumes(s, d) {
			wg.Add(1)
			go func(source string) {
				defer wg.Done()
				// if service is never started before, data volume won't be created and Docker Engine
				// will return with the not found error. therefore, we can safely ignore it.
				if err := m.c.DeleteVolume(source); err != nil && !client.IsErrNotFound(err) {
					errs.Append(err)
				}
			}(volume.Source)
		}
	}
	wg.Wait()
	return errs.ErrorOrNil()
}
