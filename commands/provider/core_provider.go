package provider

import (
	"context"
	"io"
	"sync"

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

// Stop stops all running services and core daemon.
func (p *CoreProvider) Stop() error {
	listReply, err := p.client.ListServices(context.Background(), &coreapi.ListServicesRequest{})
	if err != nil {
		return err
	}

	var (
		serviceLen = len(listReply.Services)
		errC       = make(chan error, serviceLen)
		wg         sync.WaitGroup
	)

	wg.Add(serviceLen)
	for _, service := range listReply.Services {
		go func(hash string) {
			defer wg.Done()
			_, err := p.client.StopService(context.Background(), &coreapi.StopServiceRequest{
				ServiceID: hash,
			})
			if err != nil {
				errC <- err
			}
		}(service.Hash)
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
