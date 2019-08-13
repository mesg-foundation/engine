package instancesdk

import (
	"sync"

	"github.com/docker/docker/client"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/x/xerrors"
)

// deleteData deletes the data volumes of instance and its dependencies.
// TODO: right now deleteData() is not compatible to work with multiple instances of same
// service since extractVolumes() accepts service, not instance. same applies in the start
// api as well. make it work with multiple instances.
func (i *Instance) deleteData(inst *instance.Instance) error {
	s, err := i.service.Get(inst.ServiceHash)
	if err != nil {
		return err
	}
	var (
		wg      sync.WaitGroup
		errs    xerrors.SyncErrors
		volumes = make([]container.Mount, 0)
	)

	for _, d := range s.Dependencies {
		volumes = append(volumes, extractVolumes(s, d.Configuration, d.Key)...)
	}
	volumes = append(volumes, extractVolumes(s, s.Configuration, service.MainServiceKey)...)

	for _, volume := range volumes {
		// TODO: is it actually necessary to remvoe in parallel? I worry that some volume will be deleted at the same time and create side effect
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
	wg.Wait()
	return errs.ErrorOrNil()
}
