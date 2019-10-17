package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
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
	InstanceHash []byte
}

// New creates a new client from env variables supplied by mesg engine.
func New() (*Client, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("client: mesg server address env(%s) is empty", envMesgEndpoint)
	}

	instanceHash, err := hash.Decode(os.Getenv(envMesgInstanceHash))
	if err != nil {
		return nil, fmt.Errorf("client: error with mesg's instance hash env(%s): %s", envMesgInstanceHash, err.Error())
	}

	dialoptions := []grpc.DialOption{
		// Keep alive prevents Docker network to drop TCP idle connections after 15 minutes.
		// See: https://forum.mesg.com/t/solution-summary-for-docker-dropping-connections-after-15-min/246
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 5 * time.Minute, // 5 minutes is the minimun time of gRPC enforcement policy.
		}),
		grpc.WithTimeout(10 * time.Second),
		grpc.WithInsecure(),
	}

	conn, err := grpc.DialContext(context.Background(), endpoint, dialoptions...)
	if err != nil {
		return nil, fmt.Errorf("client: connection error: %s", err)
	}

	return &Client{
		ExecutionClient: pb.NewExecutionClient(conn),
		InstanceHash:    instanceHash,
	}, nil
}

func main() {
	client, err := New()
	if err != nil {
		log.Fatal(err)
	}

	stream, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		Filter: &pb.StreamExecutionRequest_Filter{
			Statuses:     []execution.Status{execution.Status_InProgress},
			InstanceHash: client.InstanceHash,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		exec, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}

		outputs, err := r.defs[exec.TaskKey](exec.Inputs)
		req := &pb.UpdateExecutionRequest{
			Hash: exec.Hash,
		}
		req.Result = &pb.UpdateExecutionRequest_Outputs{
			Outputs: outputs,
		}
		if _, err := client.ExecutionClient.Update(context.Background(), req); err != nil {
			return err
		}
	}
}
