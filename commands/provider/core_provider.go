package provider

import (
	"context"
	"io"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xerrors"
)

// CoreProvider is a struct that provides all methods required by core command.
type CoreProvider struct {
	client core.CoreClient
}

// NewCoreProvider creates new CoreProvider.
func NewCoreProvider(client core.CoreClient) *CoreProvider {
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
	ids, err := service.ListRunning()
	if err != nil {
		return err
	}

	var (
		idsLen = len(ids)
		errC   = make(chan error, idsLen)
	)

	for _, id := range ids {
		go func(id string) {
			_, err := p.client.StopService(context.Background(), &core.StopServiceRequest{
				ServiceID: id,
			})
			errC <- err
		}(id)
	}

	var errs xerrors.Errors

	for i := 0; i < idsLen; i++ {
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
