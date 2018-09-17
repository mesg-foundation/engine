package provider

import (
	"context"
	"io"
	"sync"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/protobuf/core"
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
	var wg sync.WaitGroup

	ids, err := service.ListRunning()
	if err != nil {
		return err
	}

	var errC = make(chan error, len(ids))
	wg.Add(len(ids))

	for _, id := range ids {
		go func(id string) {
			defer wg.Done()
			if _, err := p.client.StopService(context.Background(), &core.StopServiceRequest{
				ServiceID: id,
			}); err == nil {
				errC <- err
			}
		}(id)
	}
	wg.Wait()

	var errs xerrors.Errors
loop:
	for {
		select {
		case err := <-errC:
			if err != nil {
				errs = append(errs, err)
			}
		default:
			break loop
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
