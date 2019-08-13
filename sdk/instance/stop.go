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
		namespaces = make([]string, 0)
	)

	for _, d := range srv.Dependencies {
		namespaces = append(namespaces, dependencyNamespace(sNamespace, d.Key))
	}
	namespaces = append(namespaces, dependencyNamespace(sNamespace, service.MainServiceKey))

	for _, namespace := range namespaces {
		wg.Add(1)
		go func(namespace string) {
			defer wg.Done()
			if err := i.container.StopService(namespace); err != nil {
				errs.Append(err)
			}
		}(namespace)
	}
	wg.Wait()
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	return i.container.DeleteNetwork(sNamespace)
}
