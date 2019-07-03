package client

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "github.com/mesg-foundation/engine/protobuf/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// env variables for configure mesg client.
	envMesgEndpoint     = "MESG_ENDPOINT"
	envMesgInstanceHash = "MESG_TOKEN"
)

// Client is a client to connect to all mesg exposed API.
type Client struct {
	// all clients registered by mesg server.
	pb.EventClient
	pb.ExecutionClient
	pb.InstanceClient
	pb.ServiceClient

	// instance hash that could be used in api calls.
	InstanceHash string
}

// New creates a new client from env variables supplied by mesg engine.
func New() (*Client, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("client: mesg server address env(%s) is empty", envMesgEndpoint)
	}

	instanceHash := os.Getenv(envMesgInstanceHash)
	if instanceHash == "" {
		return nil, fmt.Errorf("client: mesg instance hash env(%s) is empty", envMesgInstanceHash)
	}

	conn, err := grpc.DialContext(context.Background(), endpoint, dialoptions()...)
	if err != nil {
		return nil, fmt.Errorf("client: connection error: %s", err)
	}

	return &Client{
		EventClient:     pb.NewEventClient(conn),
		ExecutionClient: pb.NewExecutionClient(conn),
		InstanceClient:  pb.NewInstanceClient(conn),
		ServiceClient:   pb.NewServiceClient(conn),
		InstanceHash:    instanceHash,
	}, nil
}

// dialoptions returns all options used to create client connection.
func dialoptions() []grpc.DialOption {
	return []grpc.DialOption{
		// Keep alive prevents Docker network to drop TCP idle connections after 15 minutes.
		// See: https://forum.mesg.com/t/solution-summary-for-docker-dropping-connections-after-15-min/246
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 5 * time.Minute, // 5 minutes is the minimun time of gRPC enforcement policy.
		}),
		grpc.WithTimeout(10 * time.Second),
		grpc.WithInsecure(),
	}
}

// TaskRunner returns runner that could process execution
// without direct access to MESG api.
func (c *Client) TaskRunner() *TaskRunner {
	return &TaskRunner{
		defs:   make(map[string]TaskFn),
		client: c,
	}
}
