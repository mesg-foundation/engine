package workflow

import (
	"context"
	"io"
	"time"

	core "github.com/mesg-foundation/core/systemservices/sources/workflow/workflow/proto"
	"google.golang.org/grpc"
)

// coreClientProvider provides functionalities to create a new mesg core client.
type coreClientProvider interface {
	// New creates a new mesg core client for given address.
	New(address string) (coreClient, error)
}

// coreClient provides functionalities from core APIs of mesg.
type coreClient interface {
	// CoreClient provides core APIs of mesg.
	core.CoreClient

	// Close closes client connection.
	Close() error
}

// defaultCoreClientProvider implements core CoreClientProvider.
type defaultCoreClientProvider struct {
	timeout time.Duration
}

// defaultCoreClient implements CoreClient.
type defaultCoreClient struct {
	core.CoreClient
	io.Closer
}

// newCoreClient returns a new client for address.
func (p *defaultCoreClientProvider) New(address string) (coreClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := core.NewCoreClient(conn)
	return &defaultCoreClient{
		CoreClient: client,
		Closer:     conn,
	}, nil
}
