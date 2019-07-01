package instancesdk

import (
	"sync"

	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/x/xerrors"
)

// Stop stops an instance.
func (i *Instance) stop(inst *instance.Instance) error {
	srv, err := i.service.Get(inst.ServiceHash)
	if err != nil {
		return err
	}

	var (
		wg         sync.WaitGroup
		errs       xerrors.SyncErrors
		sNamespace = instanceNamespace(inst.Hash)
	)
	for _, d := range append([]*service.Dependency{srv.Configuration}, srv.Dependencies...) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		wg.Add(1)
		go func(namespace string) {
			defer wg.Done()
			if err := i.container.StopService(namespace); err != nil {
				errs.Append(err)
			}
		}(dependencyNamespace(sNamespace, d.Key))
	}
	wg.Wait()
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	return i.container.DeleteNetwork(sNamespace)
}
