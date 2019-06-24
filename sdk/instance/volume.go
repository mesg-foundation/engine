package instancesdk

import (
	"sync"

	"github.com/docker/docker/client"
	"github.com/mesg-foundation/core/instance"
	"github.com/mesg-foundation/core/x/xerrors"
)

// deleteData deletes the data volumes of instance and its dependencies.
// TODO: right now deleteData() is not compatible to work with multiple instances of same
// service since extractVolumes() accepts service, not instance. same applies in the start
// api as well. make it work with multiple instances.
func (i *Instance) deleteData(inst *instance.Instance) error {
	s, err := i.serviceDB.Get(inst.ServiceHash)
	if err != nil {
		return err
	}
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
				if err := i.container.DeleteVolume(source); err != nil && !client.IsErrNotFound(err) {
					errs.Append(err)
				}
			}(volume.Source)
		}
	}
	wg.Wait()
	return errs.ErrorOrNil()
}
