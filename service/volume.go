package service

import (
	"sync"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/x/xerrors"
)

// DeleteVolumes deletes the data volumes of service and its dependencies.
func (s *Service) DeleteVolumes(c container.Container) error {
	var (
		wg   sync.WaitGroup
		errs xerrors.SyncErrors
	)
	for _, d := range append(s.Dependencies, s.Configuration) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		for _, volume := range d.extractVolumes(s) {
			wg.Add(1)
			go func(source string) {
				defer wg.Done()
				if err := c.DeleteVolume(source); err != nil {
					errs.Append(err)
				}
			}(volume.Source)
		}
	}
	wg.Wait()
	return errs.ErrorOrNil()
}
