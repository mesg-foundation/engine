package provider

import (
	"context"
	"io"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/x/xerrors"
)

// CoreProvider is a struct that provides all methods required by core command.
type CoreProvider struct {
	client coreapi.CoreClient
}

// NewCoreProvider creates new CoreProvider.
func NewCoreProvider(client coreapi.CoreClient) *CoreProvider {
	return &CoreProvider{
		client: client,
	}
}

// Start starts core daemon.
func (p *CoreProvider) Start() error {
	_, err := daemon.Start()
	return err
}

// Stop stops core daemon and all running services.
func (p *CoreProvider) Stop() error {
	res, err := p.client.ListServices(context.Background(), &coreapi.ListServicesRequest{
		FilterRunning: true,
	})
	if err != nil {
		return err
	}

	var (
		serviceLen = len(res.Services)
		errC       = make(chan error, serviceLen)
	)

	for _, service := range res.Services {
		go func(service *coreapi.Service) {
			_, err := p.client.StopService(context.Background(), &coreapi.StopServiceRequest{
				ServiceID: service.ID,
			})
			errC <- err
		}(service)
	}

	var errs xerrors.Errors

	for i := 0; i < serviceLen; i++ {
		if err := <-errC; err != nil {
			errs = append(errs, err)
		}
	}

	if err := daemon.Stop(); err != nil {
		errs = append(errs, err)
	}

	return errs.ErrorOrNil()
}

// Status returns daemon status.
func (p *CoreProvider) Status() (container.StatusType, error) {
	return daemon.Status()
}

// Logs returns daemon logs reader.
func (p *CoreProvider) Logs() (io.ReadCloser, error) {
	return daemon.Logs()
}
