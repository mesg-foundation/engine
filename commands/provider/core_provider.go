package provider

import (
	"context"
	"io"
	"sync"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xerrors"
)

type CoreProvider struct {
	client core.CoreClient
}

func NewCoreProvider(client core.CoreClient) *CoreProvider {
	return &CoreProvider{
		client: client,
	}
}

func (p *CoreProvider) Start() error {
	_, err := daemon.Start()
	return err
}

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
			errs = append(errs, err)
		default:
			break loop
		}
	}

	if err := daemon.Stop(); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (p *CoreProvider) Status() (container.StatusType, error) {
	return daemon.Status()
}

func (p *CoreProvider) Logs() (io.ReadCloser, error) {
	return daemon.Logs()
}
