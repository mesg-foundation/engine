package dockermanager

import (
	"sync"

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
				if err := m.c.DeleteVolume(source); err != nil {
					errs.Append(err)
				}
			}(volume.Source)
		}
	}
	wg.Wait()
	return errs.ErrorOrNil()
}
