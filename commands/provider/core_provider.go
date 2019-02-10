package provider

import (
	"context"
	"io"
	"sync"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/x/xerrors"
)

// CoreProvider is a struct that provides all methods required by core command.
type CoreProvider struct {
	d      daemon.Daemon
	client coreapi.CoreClient
}

// NewCoreProvider creates new CoreProvider.
func NewCoreProvider(client coreapi.CoreClient, d daemon.Daemon) *CoreProvider {
	return &CoreProvider{
		client: client,
		d:      d,
	}
}

// Start starts core daemon.
func (p *CoreProvider) Start() error {
	return p.d.Start()
}

// Stop stops core daemon and all running services.
func (p *CoreProvider) Stop() error {
	ids, err := ListRunning()
	if err != nil {
		return err
	}

	var (
		idsLen = len(ids)
		errC   = make(chan error, idsLen)
		wg     sync.WaitGroup
	)

	wg.Add(idsLen)
	for _, id := range ids {
		go func(id string) {
			defer wg.Done()
			_, err := p.client.StopService(context.Background(), &coreapi.StopServiceRequest{
				ServiceID: id,
			})
			if err != nil {
				errC <- err
			}
		}(id)
	}

	wg.Wait()
	close(errC)

	var errs xerrors.Errors
	for err := range errC {
		errs = append(errs, err)
	}

	if err := p.d.Stop(); err != nil {
		errs = append(errs, err)
	}

	return errs.ErrorOrNil()
}

// Status returns daemon status.
func (p *CoreProvider) Status() (container.StatusType, error) {
	return p.d.Status()
}

// Logs returns daemon logs reader.
func (p *CoreProvider) Logs() (io.ReadCloser, error) {
	return p.d.Logs()
}

// ListRunning returns all the running services.2
// TODO: should move to another file
func ListRunning() ([]string, error) {
	cfg, err := config.Global()
	if err != nil {
		return nil, err
	}

	// TODO(ilgooz): remove this line after ListRunning refactored.
	c, err := container.New()
	if err != nil {
		return nil, err
	}
	services, err := c.ListServices("mesg.hash", "mesg.core="+cfg.Core.Name)
	if err != nil {
		return nil, err
	}
	// Make service list unique. One mesg service can have multiple docker service.
	mapRes := make(map[string]uint)
	for _, service := range services {
		serviceName := service.Spec.Annotations.Labels["mesg.hash"]
		mapRes[serviceName]++
	}
	res := make([]string, 0, len(mapRes))
	for k := range mapRes {
		res = append(res, k)
	}
	return res, nil
}
