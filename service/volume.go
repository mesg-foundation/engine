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
	for _, d := range s.Dependencies {
		wg.Add(1)
		go func(d *Dependency) {
			defer wg.Done()
			if err := d.DeleteVolumes(c, s); err != nil {
				errs.Append(err)
			}
		}(d)
	}
	wg.Wait()
	return errs.ErrorOrNil()
}

// DeleteVolumes deletes the data volumes of service's dependency.
func (d *Dependency) DeleteVolumes(c container.Container, s *Service) error {
	volumes := d.extractVolumes(s)
	var (
		wg   sync.WaitGroup
		errs xerrors.SyncErrors
	)
	for _, mount := range volumes {
		wg.Add(1)
		go func(mount container.Mount) {
			defer wg.Done()
			if err := c.DeleteVolume(mount.Source); err != nil {
				errs.Append(err)
			}
		}(mount)
	}
	wg.Wait()
	return errs.ErrorOrNil()
}
