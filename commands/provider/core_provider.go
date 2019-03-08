package provider

import (
	"io"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/protobuf/coreapi"
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
	return p.d.Stop()
}

// Status returns daemon status.
func (p *CoreProvider) Status() (container.StatusType, error) {
	return p.d.Status()
}

// Logs returns daemon logs reader.
func (p *CoreProvider) Logs() (io.ReadCloser, error) {
	return p.d.Logs()
}
